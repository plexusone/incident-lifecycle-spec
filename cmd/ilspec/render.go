package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/plexusone/incident-lifecycle-spec/pkg/render"
	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
	"github.com/plexusone/incident-lifecycle-spec/pkg/visualize"
	"github.com/spf13/cobra"
)

func renderCmd() *cobra.Command {
	var (
		outputFile   string
		templateName string
		templateDir  string
		diagrams     bool
	)

	cmd := &cobra.Command{
		Use:   "render <incident.json>",
		Short: "Render an incident artifact to Markdown",
		Long: `Render an incident JSON file to Markdown using the appropriate template.

By default, the template is selected based on the incident phase:
  - premortem    → premortem.md.tmpl
  - intra_mortem → intra-mortem.md.tmpl
  - postmortem   → postmortem.md.tmpl

You can override the template with --template.

Use --diagrams to embed D2 diagram code blocks in the output.
These render natively on GitHub and GitLab.`,
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

			// Append diagrams if requested
			if diagrams {
				diagramSection, err := generateDiagramSection(incident)
				if err != nil {
					return fmt.Errorf("generating diagrams: %w", err)
				}
				output = output + "\n" + diagramSection
			}

			outFile := outputFilename(outputFile, inputFile, incident.Phase)
			if outFile == "-" {
				fmt.Print(output)
			} else {
				if err := os.WriteFile(outFile, []byte(output), 0600); err != nil {
					return fmt.Errorf("writing output: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Wrote %s\n", outFile)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: auto-generated, use '-' for stdout)")
	cmd.Flags().StringVarP(&templateName, "template", "t", "", "Template name (default: auto-detect from phase)")
	cmd.Flags().StringVar(&templateDir, "template-dir", "", "Directory containing custom templates")
	cmd.Flags().BoolVar(&diagrams, "diagrams", false, "Embed D2 diagram code blocks in output")

	return cmd
}

// generateDiagramSection creates D2 code blocks for hypothesis and timeline diagrams.
func generateDiagramSection(incident *types.Incident) (string, error) {
	vis := visualize.NewWithoutRenderer()

	var sb strings.Builder
	sb.WriteString("## Diagrams\n\n")

	// Hypothesis lifecycle diagram
	if len(incident.Hypotheses) > 0 {
		sb.WriteString("### Hypothesis Lifecycle\n\n")
		sb.WriteString("```d2\n")
		d2Code, err := vis.Generate(incident, visualize.DiagramHypothesis)
		if err != nil {
			return "", err
		}
		sb.WriteString(d2Code)
		sb.WriteString("```\n\n")
	}

	// Timeline diagram
	if len(incident.Timeline) > 0 {
		sb.WriteString("### Timeline\n\n")
		sb.WriteString("```d2\n")
		d2Code, err := vis.Generate(incident, visualize.DiagramTimeline)
		if err != nil {
			return "", err
		}
		sb.WriteString(d2Code)
		sb.WriteString("```\n")
	}

	return sb.String(), nil
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

// outputFilename returns the output filename, auto-generating one if not specified.
// Use "-" for explicit stdout.
func outputFilename(specified, inputFile string, phase types.Phase) string {
	if specified != "" {
		return specified
	}
	// Auto-generate based on input filename and phase
	base := filepath.Base(inputFile)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	switch phase {
	case types.PhaseIntraMortem:
		return name + "-update.md"
	case types.PhasePostmortem:
		return name + "-postmortem.md"
	case types.PhasePremortem:
		return name + "-premortem.md"
	default:
		return name + ".md"
	}
}
