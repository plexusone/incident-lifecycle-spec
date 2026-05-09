package main

import (
	"context"
	"fmt"
	"os"

	"github.com/plexusone/incident-lifecycle-spec/pkg/render"
	"github.com/plexusone/incident-lifecycle-spec/pkg/visualize"
	"github.com/spf13/cobra"
)

func visualizeCmd() *cobra.Command {
	var (
		diagramType string
		format      string
		outputFile  string
	)

	cmd := &cobra.Command{
		Use:   "visualize <incident.json>",
		Short: "Generate D2 diagrams from incident data",
		Long: `Generate D2 diagrams visualizing incident data.

Diagram types:
  hypothesis - Shows hypothesis state transitions (proposed → investigating → validated/invalidated)
  timeline   - Shows timeline events with evidence links
  all        - Combined diagram with all visualizations

Output formats:
  d2  - D2 source code (default)
  svg - Rendered SVG image

Examples:
  # Generate hypothesis lifecycle diagram as D2
  ilspec visualize incident.json --type hypothesis

  # Generate timeline as SVG
  ilspec visualize incident.json --type timeline --format svg -o timeline.svg

  # Generate all diagrams to stdout
  ilspec visualize incident.json --type all -o -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]

			// Load incident
			incident, err := render.LoadIncident(inputFile)
			if err != nil {
				return fmt.Errorf("loading incident: %w", err)
			}

			// Parse diagram type
			var dt visualize.DiagramType
			switch diagramType {
			case "hypothesis":
				dt = visualize.DiagramHypothesis
			case "timeline":
				dt = visualize.DiagramTimeline
			case "all":
				dt = visualize.DiagramAll
			default:
				return fmt.Errorf("invalid diagram type %q: must be hypothesis, timeline, or all", diagramType)
			}

			// Parse format
			var fmt_ visualize.Format
			switch format {
			case "d2":
				fmt_ = visualize.FormatD2
			case "svg":
				fmt_ = visualize.FormatSVG
			case "png":
				fmt_ = visualize.FormatPNG
			default:
				return fmt.Errorf("invalid format %q: must be d2, svg, or png", format)
			}

			// Create visualizer
			var vis *visualize.Visualizer
			if fmt_ == visualize.FormatD2 {
				vis = visualize.NewWithoutRenderer()
			} else {
				vis, err = visualize.New()
				if err != nil {
					return fmt.Errorf("creating visualizer: %w", err)
				}
			}

			// Generate diagram
			output, err := vis.Render(context.Background(), incident, dt, fmt_)
			if err != nil {
				return fmt.Errorf("generating diagram: %w", err)
			}

			// Write output
			if outputFile == "-" || outputFile == "" {
				_, err = os.Stdout.Write(output)
				return err
			}

			if err := os.WriteFile(outputFile, output, 0600); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Generated %s\n", outputFile)
			return nil
		},
	}

	cmd.Flags().StringVarP(&diagramType, "type", "t", "all", "Diagram type (hypothesis, timeline, all)")
	cmd.Flags().StringVarP(&format, "format", "f", "d2", "Output format (d2, svg, png)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "-", "Output file (use '-' for stdout)")

	return cmd
}
