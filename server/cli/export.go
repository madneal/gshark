package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func exportCommand(connect connector) *cobra.Command {
	var file, format, keyword, query string
	var status, maxResults int
	var raw bool
	cmd := &cobra.Command{Use: "export", Short: "Export paginated results to a new JSONL or CSV file", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) (err error) {
		if file == "" {
			return fmt.Errorf("--file is required")
		}
		if format != "jsonl" && format != "csv" {
			return fmt.Errorf("--format must be jsonl or csv")
		}
		if status < -1 || status > 2 {
			return fmt.Errorf("--status must be -1, 0, 1 or 2")
		}
		if maxResults < 1 {
			return fmt.Errorf("--max-results must be positive")
		}
		c, err := connect(cmd)
		if err != nil {
			return err
		}
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		success := false
		defer func() {
			f.Close()
			if !success {
				os.Remove(file)
			}
		}()
		csvWriter := csv.NewWriter(f)
		keys := []string{"ID", "repo", "RepoUrl", "matches", "keyword", "path", "url", "status"}
		if format == "csv" {
			if err = csvWriter.Write(keys); err != nil {
				return err
			}
		}
		encoder := json.NewEncoder(f)
		count, total := 0, int64(0)
		for page := 1; ; page++ {
			q := url.Values{"page": {strconv.Itoa(page)}, "pageSize": {"100"}, "status": {strconv.Itoa(status)}}
			if keyword != "" {
				q.Set("keyword", keyword)
			}
			if query != "" {
				q.Set("query", query)
			}
			data, err := c.Get(cmd.Context(), "searchResult/getSearchResultList", q)
			if err != nil {
				return err
			}
			var result struct {
				List  []map[string]interface{} `json:"list"`
				Total int64                    `json:"total"`
			}
			if err = json.Unmarshal(data, &result); err != nil {
				return err
			}
			total = result.Total
			for _, row := range result.List {
				if count == maxResults {
					break
				}
				if !raw {
					redact(row)
				}
				if format == "jsonl" {
					if err = encoder.Encode(row); err != nil {
						return err
					}
				} else {
					record := make([]string, len(keys))
					for i, key := range keys {
						if val, ok := row[key]; ok {
							record[i] = fmt.Sprint(val)
						}
						// CSVs are commonly opened in spreadsheets; neutralize formula cells.
						trimmed := strings.TrimLeft(record[i], " \t\r\n")
						if len(trimmed) > 0 && strings.ContainsAny(trimmed[:1], "=+-@") {
							record[i] = "'" + record[i]
						}
					}
					if err = csvWriter.Write(record); err != nil {
						return err
					}
				}
				count++
			}
			if len(result.List) == 0 || count == maxResults || int64(count) >= total {
				break
			}
		}
		if format == "csv" {
			csvWriter.Flush()
			if err = csvWriter.Error(); err != nil {
				return err
			}
		}
		if err = f.Close(); err != nil {
			return err
		}
		success = true
		return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]interface{}{"file": file, "exported": count, "total": total, "truncated": int64(count) < total})
	}}
	cmd.Flags().StringVar(&file, "file", "", "Destination file (must not already exist)")
	cmd.Flags().StringVar(&format, "format", "jsonl", "Export format: jsonl or csv")
	cmd.Flags().StringVar(&keyword, "keyword", "", "Exact keyword filter")
	cmd.Flags().StringVar(&query, "query", "", "Repository/content search")
	cmd.Flags().IntVar(&status, "status", -1, "-1 all, 0 unprocessed, 1 processed, 2 ignored")
	cmd.Flags().IntVar(&maxResults, "max-results", 10000, "Maximum records; summary reports truncation")
	cmd.Flags().BoolVar(&raw, "include-sensitive", false, "Include original sensitive result fields")
	return cmd
}
