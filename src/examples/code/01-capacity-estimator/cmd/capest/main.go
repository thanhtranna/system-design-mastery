// Command capest is a capacity estimator CLI.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanhtranna/system-design-mastery/examples/capest/estimator"
	"gopkg.in/yaml.v3"
)

var (
	preset string
	file   string
	format string
)

func main() {
	root := &cobra.Command{
		Use:   "capest",
		Short: "Capacity estimator for system design",
	}

	estimateCmd := &cobra.Command{
		Use:   "estimate",
		Short: "Run capacity estimate from preset or file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var s estimator.Scenario
			switch {
			case preset != "":
				p, ok := estimator.Presets[preset]
				if !ok {
					return fmt.Errorf("unknown preset %q; available: %v", preset, presetNames())
				}
				s = p
			case file != "":
				data, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}
				if err := yaml.Unmarshal(data, &s); err != nil {
					return fmt.Errorf("parse yaml: %w", err)
				}
			default:
				return fmt.Errorf("must provide --preset or --file")
			}

			if err := s.Validate(); err != nil {
				return fmt.Errorf("invalid scenario: %w", err)
			}

			r := estimator.Estimate(s)

			switch format {
			case "markdown":
				fmt.Print(estimator.FormatMarkdown(s, r))
			case "json":
				return json.NewEncoder(os.Stdout).Encode(map[string]any{
					"scenario": s, "result": r,
				})
			default:
				fmt.Print(estimator.FormatText(s, r))
			}
			return nil
		},
	}

	estimateCmd.Flags().StringVar(&preset, "preset", "", "preset name (twitter, instagram, whatsapp, uber, propertyhub)")
	estimateCmd.Flags().StringVar(&file, "file", "", "YAML scenario file")
	estimateCmd.Flags().StringVar(&format, "format", "text", "output format: text, markdown, json")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available presets",
		Run: func(_ *cobra.Command, _ []string) {
			for name, s := range estimator.Presets {
				fmt.Printf("%-15s — %s\n", name, s.Name)
			}
		},
	}

	root.AddCommand(estimateCmd, listCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func presetNames() []string {
	names := make([]string, 0, len(estimator.Presets))
	for name := range estimator.Presets {
		names = append(names, name)
	}
	return names
}
