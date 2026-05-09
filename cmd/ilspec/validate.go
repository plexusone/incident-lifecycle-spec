package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/plexusone/incident-lifecycle-spec/pkg/schema"
	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
	"github.com/spf13/cobra"
)

func validateCmd() *cobra.Command {
	var (
		quiet      bool
		schemaOnly bool
	)

	cmd := &cobra.Command{
		Use:   "validate <incident.json>",
		Short: "Validate an incident JSON file against the schema",
		Long: `Validate an incident JSON file to ensure it conforms to the
incident lifecycle schema. Checks for:

  - Valid JSON syntax
  - JSON Schema validation (structure, types, required fields)
  - Go type validation (semantic checks, enum values)
  - Valid structure for timeline, hypotheses, action_items, evidence

Exit codes:
  0 - Valid
  1 - Invalid or error`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]

			data, err := os.ReadFile(inputFile)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			// JSON Schema validation
			validator, err := schema.NewValidator()
			if err != nil {
				return fmt.Errorf("creating schema validator: %w", err)
			}

			schemaErrors, err := validator.ValidateBytesDetailed(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
				os.Exit(1)
			}

			if len(schemaErrors) > 0 {
				for _, e := range schemaErrors {
					if e.Path != "" {
						fmt.Fprintf(os.Stderr, "Schema error at %s: %s\n", e.Path, e.Message)
					} else {
						fmt.Fprintf(os.Stderr, "Schema error: %s\n", e.Message)
					}
				}
				os.Exit(1)
			}

			// Skip Go validation if only schema validation requested
			if schemaOnly {
				if !quiet {
					fmt.Printf("✓ %s passes JSON Schema validation\n", inputFile)
				}
				return nil
			}

			// Go type validation for additional semantic checks
			var incident types.Incident
			if err := json.Unmarshal(data, &incident); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
				os.Exit(1)
			}

			errors := validateIncident(&incident)
			if len(errors) > 0 {
				for _, e := range errors {
					fmt.Fprintf(os.Stderr, "Error: %s\n", e)
				}
				os.Exit(1)
			}

			if !quiet {
				fmt.Printf("✓ %s is valid\n", inputFile)
				fmt.Printf("  Phase: %s\n", incident.Phase)
				fmt.Printf("  Severity: %s\n", incident.Severity)
				if incident.Status != "" {
					fmt.Printf("  Status: %s\n", incident.Status)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output on success")
	cmd.Flags().BoolVar(&schemaOnly, "schema-only", false, "Only run JSON Schema validation (skip Go type validation)")

	return cmd
}

func validateIncident(incident *types.Incident) []string {
	var errors []string

	// Required fields
	if incident.IncidentID == "" {
		errors = append(errors, "missing required field: incident_id")
	}
	if incident.Title == "" {
		errors = append(errors, "missing required field: title")
	}
	if incident.Phase == "" {
		errors = append(errors, "missing required field: phase")
	}
	if incident.Severity == "" {
		errors = append(errors, "missing required field: severity")
	}
	if incident.CreatedAt == nil {
		errors = append(errors, "missing required field: created_at")
	}

	// Valid phase
	if incident.Phase != "" {
		validPhases := map[types.Phase]bool{
			types.PhasePremortem:   true,
			types.PhaseIntraMortem: true,
			types.PhasePostmortem:  true,
		}
		if !validPhases[incident.Phase] {
			errors = append(errors, fmt.Sprintf("invalid phase: %q (must be premortem, intra_mortem, or postmortem)", incident.Phase))
		}
	}

	// Valid severity
	if incident.Severity != "" {
		validSeverities := map[types.Severity]bool{
			types.SeveritySEV0: true,
			types.SeveritySEV1: true,
			types.SeveritySEV2: true,
			types.SeveritySEV3: true,
		}
		if !validSeverities[incident.Severity] {
			errors = append(errors, fmt.Sprintf("invalid severity: %q (must be SEV0, SEV1, SEV2, or SEV3)", incident.Severity))
		}
	}

	// Valid status (if provided)
	if incident.Status != "" {
		validStatuses := map[types.Status]bool{
			types.StatusHypothetical:  true,
			types.StatusInvestigating: true,
			types.StatusIdentified:    true,
			types.StatusMitigating:    true,
			types.StatusResolved:      true,
			types.StatusClosed:        true,
		}
		if !validStatuses[incident.Status] {
			errors = append(errors, fmt.Sprintf("invalid status: %q", incident.Status))
		}
	}

	// Validate timeline events
	for i, event := range incident.Timeline {
		if event.EventID == "" {
			errors = append(errors, fmt.Sprintf("timeline[%d]: missing required field: event_id", i))
		}
		if event.Description == "" {
			errors = append(errors, fmt.Sprintf("timeline[%d]: missing required field: description", i))
		}
	}

	// Validate hypotheses
	for i, hyp := range incident.Hypotheses {
		if hyp.HypothesisID == "" {
			errors = append(errors, fmt.Sprintf("hypotheses[%d]: missing required field: hypothesis_id", i))
		}
		if hyp.Description == "" {
			errors = append(errors, fmt.Sprintf("hypotheses[%d]: missing required field: description", i))
		}
		if hyp.Status == "" {
			errors = append(errors, fmt.Sprintf("hypotheses[%d]: missing required field: status", i))
		}
	}

	// Validate action items
	for i, action := range incident.ActionItems {
		if action.ActionID == "" {
			errors = append(errors, fmt.Sprintf("action_items[%d]: missing required field: action_id", i))
		}
		if action.Description == "" {
			errors = append(errors, fmt.Sprintf("action_items[%d]: missing required field: description", i))
		}
		if action.Priority == "" {
			errors = append(errors, fmt.Sprintf("action_items[%d]: missing required field: priority", i))
		}
		if action.Status == "" {
			errors = append(errors, fmt.Sprintf("action_items[%d]: missing required field: status", i))
		}
	}

	// Validate evidence
	for i, ev := range incident.Evidence {
		if ev.EvidenceID == "" {
			errors = append(errors, fmt.Sprintf("evidence[%d]: missing required field: evidence_id", i))
		}
		if ev.EvidenceType == "" {
			errors = append(errors, fmt.Sprintf("evidence[%d]: missing required field: evidence_type", i))
		}
		if ev.Description == "" {
			errors = append(errors, fmt.Sprintf("evidence[%d]: missing required field: description", i))
		}
	}

	return errors
}
