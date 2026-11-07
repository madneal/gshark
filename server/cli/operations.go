package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/madneal/gshark/client"
	"github.com/spf13/cobra"
)

type connector func(*cobra.Command) (*client.Client, error)
type entitySpec struct {
	group, api, object string
	fields             map[string]string
	required           []string
}

func addOperations(root *cobra.Command, connect connector, output *string) {
	emit := func(cmd *cobra.Command, data json.RawMessage) error {
		if *output == "json" {
			var obj interface{}
			if err := json.Unmarshal(data, &obj); err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(obj)
		}
		var obj interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		return renderText(cmd, obj)
	}
	// Reject invalid output before performing writes.
	prepare := func(cmd *cobra.Command) (*client.Client, error) {
		if *output != "json" && *output != "text" {
			return nil, fmt.Errorf("--output must be json or text")
		}
		return connect(cmd)
	}
	for _, spec := range []entitySpec{
		{"rules", "rule", "Rule", map[string]string{"type": "ruleType", "name": "name", "content": "content", "description": "desc", "match-pattern": "matchPattern"}, []string{"type", "name", "content"}},
		{"repos", "repo", "Repo", map[string]string{"type": "type", "url": "url", "path": "path", "description": "desc", "project-id": "projectId", "status": "status"}, []string{"type", "url"}},
		{"filters", "filter", "Filter", map[string]string{"type": "filter_type", "class": "filter_class", "content": "content"}, []string{"type", "class", "content"}},
	} {
		group := findGroup(root, spec.group)
		group.AddCommand(&cobra.Command{Use: "get ID", Short: "Get " + spec.api + " details", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := positiveID(args[0]); err != nil {
				return err
			}
			c, err := prepare(cmd)
			if err != nil {
				return err
			}
			data, err := c.Get(cmd.Context(), spec.api+"/find"+spec.object, url.Values{"ID": {args[0]}})
			if err != nil {
				return err
			}
			return emit(cmd, data)
		}})
		for _, action := range []string{"create", "update"} {
			values := map[string]*string{}
			use := action
			if action == "update" {
				use += " ID"
			}
			cmd := &cobra.Command{Use: use, Short: strings.Title(action) + " a " + spec.api}
			if action == "update" {
				cmd.Args = cobra.ExactArgs(1)
			} else {
				cmd.Args = cobra.NoArgs
			}
			for flag := range spec.fields {
				values[flag] = cmd.Flags().String(flag, "", "Set "+flag)
			}
			if spec.group == "rules" && action == "create" {
				cmd.Flags().Bool("enabled", true, "Enable the new rule")
			}
			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				payload := map[string]interface{}{}
				changed := 0
				for flag, key := range spec.fields {
					if cmd.Flags().Changed(flag) {
						changed++
						if key == "projectId" || key == "status" {
							v, err := strconv.Atoi(*values[flag])
							if err != nil || v < 0 {
								return fmt.Errorf("--%s must be a non-negative integer", flag)
							}
							payload[key] = v
						} else {
							payload[key] = *values[flag]
						}
					}
				}
				if action == "create" {
					for _, flag := range spec.required {
						if strings.TrimSpace(*values[flag]) == "" {
							return fmt.Errorf("--%s is required", flag)
						}
					}
					if spec.group == "rules" {
						enabled, _ := cmd.Flags().GetBool("enabled")
						payload["status"] = enabled
					}
				} else if changed == 0 {
					return fmt.Errorf("specify at least one field to update")
				}
				if spec.group == "filters" {
					if v, ok := payload["filter_type"]; ok && v != "whitelist" && v != "blacklist" {
						return fmt.Errorf("--type must be whitelist or blacklist")
					}
					if v, ok := payload["filter_class"]; ok && v != "extension" && v != "keyword" {
						return fmt.Errorf("--class must be extension or keyword")
					}
				}
				var id uint64
				if action == "update" {
					var err error
					id, err = positiveID(args[0])
					if err != nil {
						return err
					}
				}
				c, err := prepare(cmd)
				if err != nil {
					return err
				}
				if action == "update" {
					data, err := c.Get(cmd.Context(), spec.api+"/find"+spec.object, url.Values{"ID": {args[0]}})
					if err != nil {
						return err
					}
					// Backend updates save full records, so merge flags into the existing object.
					var envelope map[string]map[string]interface{}
					if err = json.Unmarshal(data, &envelope); err != nil {
						return err
					}
					existing := envelope[spec.api]
					if existing == nil {
						return fmt.Errorf("missing %s in detail response", spec.api)
					}
					for key, val := range payload {
						existing[key] = val
					}
					existing["ID"] = id
					payload = existing
				}
				method := http.MethodPost
				if action == "update" {
					method = http.MethodPut
				}
				data, err := c.Mutate(cmd.Context(), method, spec.api+"/"+action+spec.object, payload)
				if err != nil {
					return err
				}
				return emit(cmd, data)
			}
			group.AddCommand(cmd)
		}
		addDelete(group, spec.api+"/delete"+spec.object+"ByIds", prepare, emit)
	}
	results := findGroup(root, "results")
	addDelete(results, "searchResult/deleteSearchResultByIds", prepare, emit)
	var status string
	mark := &cobra.Command{Use: "set-status ID [ID...]", Short: "Set result status by explicit IDs", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		states := map[string]int{"unprocessed": 0, "processed": 1, "ignored": 2}
		value, ok := states[status]
		if !ok {
			return fmt.Errorf("--status must be unprocessed, processed or ignored")
		}
		ids, err := parseIDs(args)
		if err != nil {
			return err
		}
		c, err := prepare(cmd)
		if err != nil {
			return err
		}
		data, err := c.Mutate(cmd.Context(), http.MethodPost, "searchResult/updateSearchResultStatusByIds", map[string]interface{}{"ids": ids, "status": value})
		if err != nil {
			return err
		}
		return emit(cmd, data)
	}}
	mark.Flags().StringVar(&status, "status", "", "unprocessed, processed or ignored")
	results.AddCommand(mark)
	for _, action := range []string{"enable", "disable"} {
		findGroup(root, "rules").AddCommand(&cobra.Command{Use: action + " ID", Short: strings.Title(action) + " a rule", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			id, err := positiveID(args[0])
			if err != nil {
				return err
			}
			c, err := prepare(cmd)
			if err != nil {
				return err
			}
			status := 0
			if action == "enable" {
				status = 1
			}
			data, err := c.Mutate(cmd.Context(), http.MethodPost, "rule/switchRuleStatus", map[string]interface{}{"ID": id, "status": status})
			if err != nil {
				return err
			}
			return emit(cmd, data)
		}})
	}
	results.AddCommand(exportCommand(prepare))
}

func positiveID(s string) (uint64, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("ID must be a positive 32-bit integer")
	}
	return id, nil
}
func parseIDs(args []string) ([]uint64, error) {
	if len(args) > 100 {
		return nil, fmt.Errorf("at most 100 IDs per operation")
	}
	ids := make([]uint64, 0, len(args))
	seen := map[uint64]bool{}
	for _, arg := range args {
		id, err := positiveID(arg)
		if err != nil {
			return nil, err
		}
		if !seen[id] {
			ids = append(ids, id)
			seen[id] = true
		}
	}
	return ids, nil
}
func addDelete(group *cobra.Command, path string, connect connector, emit func(*cobra.Command, json.RawMessage) error) {
	group.AddCommand(&cobra.Command{Use: "delete ID [ID...]", Short: "Delete explicitly selected records", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := parseIDs(args)
		if err != nil {
			return err
		}
		c, err := connect(cmd)
		if err != nil {
			return err
		}
		data, err := c.Mutate(cmd.Context(), http.MethodDelete, path, map[string]interface{}{"ids": ids})
		if err != nil {
			return err
		}
		return emit(cmd, data)
	}})
}
