package main

import (
	"testing"
	"time"

	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

func TestValidateIncident_Valid(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
	}

	errors := validateIncident(incident)
	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %v", errors)
	}
}

func TestValidateIncident_MissingRequiredFields(t *testing.T) {
	incident := &types.Incident{}

	errors := validateIncident(incident)

	expectedErrors := []string{
		"missing required field: incident_id",
		"missing required field: title",
		"missing required field: phase",
		"missing required field: severity",
		"missing required field: created_at",
	}

	if len(errors) != len(expectedErrors) {
		t.Errorf("Expected %d errors, got %d: %v", len(expectedErrors), len(errors), errors)
	}

	for _, expected := range expectedErrors {
		found := false
		for _, err := range errors {
			if err == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected error %q not found in %v", expected, errors)
		}
	}
}

func TestValidateIncident_InvalidPhase(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.Phase("invalid_phase"),
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
	}

	errors := validateIncident(incident)

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_InvalidSeverity(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.Severity("SEV9"),
		CreatedAt:  &now,
	}

	errors := validateIncident(incident)

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_InvalidStatus(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		Status:     types.Status("unknown_status"),
		CreatedAt:  &now,
	}

	errors := validateIncident(incident)

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_TimelineValidation(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		Timeline: []types.TimelineEvent{
			{EventID: "", Description: ""},          // Missing both
			{EventID: "evt-001", Description: ""},   // Missing description
			{EventID: "", Description: "Something"}, // Missing event_id
			{EventID: "evt-002", Description: "OK"}, // Valid
		},
	}

	errors := validateIncident(incident)

	// Should have 4 errors: 2 for first event, 1 each for second and third
	if len(errors) != 4 {
		t.Errorf("Expected 4 errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_HypothesesValidation(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		Hypotheses: []types.Hypothesis{
			{HypothesisID: "", Description: "", Status: ""},
			{HypothesisID: "hyp-001", Description: "Valid", Status: types.HypothesisProposed},
		},
	}

	errors := validateIncident(incident)

	// Should have 3 errors for first hypothesis
	if len(errors) != 3 {
		t.Errorf("Expected 3 errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_ActionItemsValidation(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		ActionItems: []types.ActionItem{
			{ActionID: "", Description: "", Priority: "", Status: ""},
			{ActionID: "act-001", Description: "Fix", Priority: types.PriorityP0, Status: types.ActionStatusOpen},
		},
	}

	errors := validateIncident(incident)

	// Should have 4 errors for first action item
	if len(errors) != 4 {
		t.Errorf("Expected 4 errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_EvidenceValidation(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		Evidence: []types.Evidence{
			{EvidenceID: "", EvidenceType: "", Description: ""},
			{EvidenceID: "evi-001", EvidenceType: types.EvidenceTypeLog, Description: "Log entry"},
		},
	}

	errors := validateIncident(incident)

	// Should have 3 errors for first evidence
	if len(errors) != 3 {
		t.Errorf("Expected 3 errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateIncident_AllPhasesValid(t *testing.T) {
	now := time.Now()
	phases := []types.Phase{
		types.PhasePremortem,
		types.PhaseIntraMortem,
		types.PhasePostmortem,
	}

	for _, phase := range phases {
		incident := &types.Incident{
			IncidentID: "INC-001",
			Title:      "Test incident",
			Phase:      phase,
			Severity:   types.SeveritySEV1,
			CreatedAt:  &now,
		}

		errors := validateIncident(incident)
		if len(errors) != 0 {
			t.Errorf("Phase %s: expected no errors, got %v", phase, errors)
		}
	}
}

func TestValidateIncident_AllSeveritiesValid(t *testing.T) {
	now := time.Now()
	severities := []types.Severity{
		types.SeveritySEV0,
		types.SeveritySEV1,
		types.SeveritySEV2,
		types.SeveritySEV3,
	}

	for _, severity := range severities {
		incident := &types.Incident{
			IncidentID: "INC-001",
			Title:      "Test incident",
			Phase:      types.PhaseIntraMortem,
			Severity:   severity,
			CreatedAt:  &now,
		}

		errors := validateIncident(incident)
		if len(errors) != 0 {
			t.Errorf("Severity %s: expected no errors, got %v", severity, errors)
		}
	}
}
