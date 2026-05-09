package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var (
		phase      string
		severity   string
		title      string
		outputFile string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new incident JSON file",
		Long: `Create a new incident JSON file with required fields populated.

The file can be created for any lifecycle phase:
  - premortem:   Proactive failure simulation
  - intra_mortem: Active incident tracking
  - postmortem:  Post-resolution analysis

Example:
  ilspec init --phase intra_mortem --severity SEV1 --title "Database outage"
  ilspec init -p premortem -s SEV2 -t "Failover risk analysis" -o risk.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate phase
			var p types.Phase
			switch phase {
			case "premortem":
				p = types.PhasePremortem
			case "intra_mortem":
				p = types.PhaseIntraMortem
			case "postmortem":
				p = types.PhasePostmortem
			default:
				return fmt.Errorf("invalid phase %q: must be premortem, intra_mortem, or postmortem", phase)
			}

			// Validate severity
			var sev types.Severity
			switch severity {
			case "SEV0":
				sev = types.SeveritySEV0
			case "SEV1":
				sev = types.SeveritySEV1
			case "SEV2":
				sev = types.SeveritySEV2
			case "SEV3":
				sev = types.SeveritySEV3
			default:
				return fmt.Errorf("invalid severity %q: must be SEV0, SEV1, SEV2, or SEV3", severity)
			}

			now := time.Now().UTC()
			incident := createIncidentTemplate(p, sev, title, now)

			data, err := json.MarshalIndent(incident, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling incident: %w", err)
			}

			// Determine output file
			outFile := outputFile
			if outFile == "" {
				outFile = generateFilename(p, now)
			}

			if outFile == "-" {
				fmt.Println(string(data))
			} else {
				if err := os.WriteFile(outFile, data, 0600); err != nil {
					return fmt.Errorf("writing file: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Created %s\n", outFile)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&phase, "phase", "p", "intra_mortem", "Lifecycle phase (premortem, intra_mortem, postmortem)")
	cmd.Flags().StringVarP(&severity, "severity", "s", "SEV2", "Severity level (SEV0, SEV1, SEV2, SEV3)")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Incident title (required)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: auto-generated, use '-' for stdout)")

	_ = cmd.MarkFlagRequired("title")

	return cmd
}

func createIncidentTemplate(phase types.Phase, severity types.Severity, title string, now time.Time) *types.Incident {
	incidentID := fmt.Sprintf("INC-%s", now.Format("2006-0102-150405"))
	if phase == types.PhasePremortem {
		incidentID = fmt.Sprintf("PRE-%s", now.Format("2006-0102-150405"))
	}

	incident := &types.Incident{
		IncidentID: incidentID,
		Title:      title,
		Phase:      phase,
		Severity:   severity,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	// Set phase-appropriate defaults
	switch phase {
	case types.PhasePremortem:
		incident.Status = types.StatusHypothetical
		incident.Summary = "TODO: Describe the potential failure scenario being analyzed."
		incident.Hypotheses = []types.Hypothesis{
			{
				HypothesisID: "hyp-001",
				Description:  "TODO: Describe a potential failure mode",
				Status:       types.HypothesisProposed,
				Confidence:   0.5,
			},
		}
		incident.ActionItems = []types.ActionItem{
			{
				ActionID:    "act-001",
				Description: "TODO: Add preventive action",
				Priority:    types.PriorityP1,
				Status:      types.ActionStatusOpen,
			},
		}

	case types.PhaseIntraMortem:
		incident.Status = types.StatusInvestigating
		incident.StartedAt = &now
		incident.Summary = "TODO: Describe what is currently happening."
		incident.Timeline = []types.TimelineEvent{
			{
				EventID:     "evt-001",
				Timestamp:   now,
				Description: "Incident detected",
				Source:      types.EventSourceMonitoring,
				Confidence:  types.ConfidenceConfirmed,
			},
		}
		incident.Hypotheses = []types.Hypothesis{
			{
				HypothesisID: "hyp-001",
				Description:  "TODO: Initial hypothesis about root cause",
				Status:       types.HypothesisInvestigating,
				Confidence:   0.5,
			},
		}

	case types.PhasePostmortem:
		incident.Status = types.StatusClosed
		incident.StartedAt = &now
		resolved := now.Add(time.Hour)
		incident.ResolvedAt = &resolved
		incident.Summary = "TODO: Summarize what happened."
		incident.RootCause = "TODO: Describe the root cause."
		incident.Timeline = []types.TimelineEvent{
			{
				EventID:     "evt-001",
				Timestamp:   now,
				Description: "Incident started",
				Source:      types.EventSourceMonitoring,
				Confidence:  types.ConfidenceConfirmed,
			},
		}
		incident.WhatWentWell = []string{"TODO: What went well"}
		incident.WhatWentWrong = []string{"TODO: What went wrong"}
		incident.LessonsLearned = []string{"TODO: Key takeaways"}
		incident.ActionItems = []types.ActionItem{
			{
				ActionID:    "act-001",
				Description: "TODO: Follow-up action to prevent recurrence",
				Priority:    types.PriorityP1,
				Status:      types.ActionStatusOpen,
			},
		}
	}

	return incident
}

func generateFilename(phase types.Phase, now time.Time) string {
	dateStr := now.Format("2006-01-02")
	switch phase {
	case types.PhasePremortem:
		return fmt.Sprintf("premortem-%s.json", dateStr)
	case types.PhaseIntraMortem:
		return fmt.Sprintf("incident-%s.json", dateStr)
	case types.PhasePostmortem:
		return fmt.Sprintf("postmortem-%s.json", dateStr)
	default:
		return fmt.Sprintf("incident-%s.json", dateStr)
	}
}
