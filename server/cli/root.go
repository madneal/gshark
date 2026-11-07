package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/madneal/gshark/client"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var server, token, output string
	var timeout time.Duration
	root := &cobra.Command{Use: "gshark-cli", Short: "Manage a deployed GShark instance", SilenceUsage: true, SilenceErrors: true}
	root.PersistentFlags().StringVar(&server, "server", envDefault("GSHARK_SERVER", "http://127.0.0.1:8888"), "Backend URL, or frontend URL including /api (GSHARK_SERVER)")
	root.PersistentFlags().StringVar(&token, "token", "", "Login token (prefer GSHARK_TOKEN to avoid shell history)")
	root.PersistentFlags().StringVar(&output, "output", "json", "Output format: json or text")
	root.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "HTTP request timeout")
	connect := func(cmd *cobra.Command) (*client.Client, error) {
		credential := token
		if credential == "" {
			credential = os.Getenv("GSHARK_TOKEN")
		}
		endpoint := server
		saved, loadErr := loadCredentials()
		if loadErr == nil && !cmd.Flags().Changed("server") && os.Getenv("GSHARK_SERVER") == "" {
			endpoint = saved.Server
		}
		useSaved := credential == ""
		if useSaved {
			if loadErr != nil {
				return nil, loadErr
			}
			if strings.TrimRight(endpoint, "/") != strings.TrimRight(saved.Server, "/") {
				return nil, fmt.Errorf("no saved login for this server; run auth login --server URL")
			}
			credential = saved.Token
		}
		c, err := client.New(endpoint, credential, timeout)
		if err != nil {
			return nil, err
		}
		if useSaved {
			c.OnToken = func(next string) error { saved.Token = next; return saveCredentials(saved) }
		}
		return c, nil
	}
	execute := func(cmd *cobra.Command, path string, q url.Values, results bool) error {
		if output != "json" && output != "text" {
			return fmt.Errorf("--output must be json or text")
		}
		c, err := connect(cmd)
		if err != nil {
			return err
		}
		data, err := c.Get(cmd.Context(), path, q)
		if err != nil {
			return err
		}
		raw, _ := cmd.Flags().GetBool("include-sensitive")
		if results && !raw {
			var value interface{}
			if err = json.Unmarshal(data, &value); err != nil {
				return err
			}
			redact(value)
			data, err = json.Marshal(value)
			if err != nil {
				return err
			}
		}
		if output == "json" {
			var pretty bytes.Buffer
			if err = json.Indent(&pretty, data, "", "  "); err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), pretty.String())
			return err
		}
		var value interface{}
		if err = json.Unmarshal(data, &value); err != nil {
			return err
		}
		return renderText(cmd, value)
	}
	for _, spec := range []struct {
		group, name, path string
		filters           map[string]string
		results           bool
	}{
		{"results", "list", "searchResult/getSearchResultList", map[string]string{"keyword": "keyword", "query": "query"}, true},
		{"rules", "list", "rule/getRuleList", map[string]string{"type": "ruleType", "name": "name", "content": "content"}, false},
		{"repos", "list", "repo/getRepoList", nil, false},
		{"filters", "list", "filter/getFilterList", nil, false},
		{"scans", "logs", "scanLog/getScanLogList", map[string]string{"provider": "provider", "status": "status", "cycle-id": "cycleId"}, false},
	} {
		group := findGroup(root, spec.group)
		var page, limit int
		filters := map[string]*string{}
		var status int
		cmd := &cobra.Command{Use: spec.name, Short: "List " + spec.group, Args: cobra.NoArgs}
		cmd.Flags().IntVar(&page, "page", 1, "Page number (starts at 1)")
		cmd.Flags().IntVar(&limit, "limit", 20, "Page size (1–100)")
		for flag := range spec.filters {
			filters[flag] = cmd.Flags().String(flag, "", "Filter by "+flag)
		}
		if spec.results {
			cmd.Flags().IntVar(&status, "status", -1, "Result status: -1 all, 0 unprocessed, 1 processed, 2 ignored")
			cmd.Flags().Bool("include-sensitive", false, "Include raw matches, keywords and source locations")
		}
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			if page < 1 || limit < 1 || limit > 100 {
				return fmt.Errorf("--page must be positive and --limit must be 1–100")
			}
			q := url.Values{"page": {strconv.Itoa(page)}, "pageSize": {strconv.Itoa(limit)}}
			for flag, key := range spec.filters {
				if *filters[flag] != "" {
					q.Set(key, *filters[flag])
				}
			}
			if spec.results {
				if status < -1 || status > 2 {
					return fmt.Errorf("--status must be -1, 0, 1 or 2")
				}
				q.Set("status", strconv.Itoa(status))
			}
			return execute(cmd, spec.path, q, spec.results)
		}
		group.AddCommand(cmd)
	}
	get := &cobra.Command{Use: "get ID", Short: "Get a result by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil || id == 0 {
			return fmt.Errorf("ID must be a positive integer")
		}
		return execute(cmd, "searchResult/findSearchResult", url.Values{"ID": {args[0]}}, true)
	}}
	get.Flags().Bool("include-sensitive", false, "Include raw matches, keywords and source locations")
	findGroup(root, "results").AddCommand(get)
	findGroup(root, "scans").AddCommand(&cobra.Command{Use: "overview", Short: "Show the latest scan cycle", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return execute(cmd, "scanLog/getScanLogOverview", nil, false)
	}})
	addOperations(root, connect, &output)
	root.AddCommand(authCommand(&server, &timeout))
	return root
}

func findGroup(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	cmd := &cobra.Command{Use: name, Short: "Manage " + name}
	root.AddCommand(cmd)
	return cmd
}

func envDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

// Result locations and search keywords can themselves contain leaked credentials.
func redact(value interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, child := range v {
			switch strings.ToLower(key) {
			case "matches", "text_matches", "keyword", "repo", "repourl", "path", "url":
				v[key] = "[redacted]"
			default:
				redact(child)
			}
		}
	case []interface{}:
		for _, child := range v {
			redact(child)
		}
	}
}

func renderText(cmd *cobra.Command, value interface{}) error {
	obj, ok := value.(map[string]interface{})
	if !ok {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), value)
		return err
	}
	if list, ok := obj["list"].([]interface{}); ok {
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
		for _, item := range list {
			row, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			keys := []string{"ID", "repo", "name", "provider", "status", "keyword", "content", "message"}
			var cells []string
			for _, key := range keys {
				if v, exists := row[key]; exists {
					cells = append(cells, key+"="+strings.Map(func(r rune) rune {
						if r < 32 || r == 127 {
							return ' '
						}
						return r
					}, fmt.Sprint(v)))
				}
			}
			if _, err := fmt.Fprintln(w, strings.Join(cells, "\t")); err != nil {
				return err
			}
		}
		if err := w.Flush(); err != nil {
			return err
		}
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "Total: %v; page: %v; pageSize: %v\n", obj["total"], obj["page"], obj["pageSize"])
		return err
	}
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		encoded, err := json.Marshal(obj[key])
		if err != nil {
			return err
		}
		if _, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, encoded); err != nil {
			return err
		}
	}
	return nil
}
