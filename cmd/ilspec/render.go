package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/plexusone/incident-lifecycle-spec/pkg/render"
	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
	"github.com/spf13/cobra"
)

func renderCmd() *cobra.Command {
	var (
		outputFile   string
		templateName string
		templateDir  string
	)

	cmd := &cobra.Command{
		Use:   "render <incident.json>",
		Short: "Render an incident artifact to Markdown",
		Long: `Render an incident JSON file to Markdown using the appropriate template.

By default, the template is selected based on the incident phase:
  - premortem    → premortem.md.tmpl (not yet implemented)
  - intra_mortem → intra-mortem.md.tmpl
  - postmortem   → postmortem.md.tmpl

You can override the template with --template.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]

			incident, err := render.LoadIncident(inputFile)
			if err != nil {
				return fmt.Errorf("loading incident: %w", err)
			}

			var renderer *render.Renderer
			if templateDir != "" {
				renderer, err = render.NewFromDir(templateDir)
			} else {
				renderer, err = render.New()
			}
			if err != nil {
				return fmt.Errorf("creating renderer: %w", err)
			}

			tmpl := templateName
			if tmpl == "" {
				tmpl = templateForPhase(incident.Phase)
			}

			output, err := renderer.Render(tmpl, incident)
			if err != nil {
				return fmt.Errorf("rendering: %w", err)
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
					return fmt.Errorf("writing output: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Wrote %s\n", outputFile)
			} else {
				fmt.Print(output)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringVarP(&templateName, "template", "t", "", "Template name (default: auto-detect from phase)")
	cmd.Flags().StringVar(&templateDir, "template-dir", "", "Directory containing custom templates")

	return cmd
}

func templateForPhase(phase types.Phase) string {
	switch phase {
	case types.PhasePremortem:
		return "premortem.md.tmpl"
	case types.PhaseIntraMortem:
		return "intra-mortem.md.tmpl"
	case types.PhasePostmortem:
		return "postmortem.md.tmpl"
	default:
		return "postmortem.md.tmpl"
	}
}

// suggestOutputFile suggests an output filename based on input and phase.
func suggestOutputFile(inputFile string, phase types.Phase) string {
	base := filepath.Base(inputFile)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	switch phase {
	case types.PhaseIntraMortem:
		return name + "-update.md"
	case types.PhasePostmortem:
		return name + "-postmortem.md"
	default:
		return name + ".md"
	}
}
