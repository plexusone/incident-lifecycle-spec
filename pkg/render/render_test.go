package render

import (
	"strings"
	"testing"
	"time"

	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

func TestNew(t *testing.T) {
	renderer, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if renderer == nil {
		t.Fatal("New() returned nil renderer")
	}
}

func TestRenderIntraMortem(t *testing.T) {
	renderer, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	now := time.Now().UTC()
	incident := &types.Incident{
		IncidentID:            "INC-2024-0001",
		Title:                 "Test Incident",
		Phase:                 types.PhaseIntraMortem,
		Severity:              types.SeveritySEV1,
		Status:                types.StatusInvestigating,
		CreatedAt:             &now,
		UpdatedAt:             &now,
		Summary:               "This is a test incident summary.",
		CustomerImpactSummary: "Customers affected",
		CustomerImpactScope:   types.ImpactScopePartial,
		ServicesAffected:      []string{"service-a", "service-b"},
		Timeline: []types.TimelineEvent{
			{
				EventID:     "evt-001",
				Timestamp:   now,
				Description: "First confirmed event",
				Confidence:  types.ConfidenceConfirmed,
			},
			{
				EventID:     "evt-002",
				Timestamp:   now,
				Description: "Suspected event",
				Confidence:  types.ConfidenceSuspected,
			},
		},
		Hypotheses: []types.Hypothesis{
			{
				HypothesisID: "hyp-001",
				Description:  "Active hypothesis",
				Status:       types.HypothesisInvestigating,
				Confidence:   0.75,
			},
			{
				HypothesisID: "hyp-002",
				Description:  "Proposed risk",
				Status:       types.HypothesisProposed,
				Confidence:   0.5,
			},
		},
		ActionItems: []types.ActionItem{
			{
				ActionID:    "act-001",
				Description: "In progress action",
				Owner:       "oncall",
				Status:      types.ActionStatusInProgress,
			},
		},
	}

	output, err := renderer.RenderIntraMortem(incident)
	if err != nil {
		t.Fatalf("RenderIntraMortem() failed: %v", err)
	}

	// Check for expected content
	checks := []string{
		"# Incident Update: Test Incident",
		"**Incident ID:** INC-2024-0001",
		"**Severity:** SEV1",
		"This is a test incident summary.",
		"service-a",
		"service-b",
		"First confirmed event",
		"Active hypothesis",
		"In progress action",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Output missing expected content: %q", check)
		}
	}
}

func TestRenderPostmortem(t *testing.T) {
	renderer, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	now := time.Now().UTC()
	incident := &types.Incident{
		IncidentID:          "INC-2024-0001",
		Title:               "Test Incident",
		Phase:               types.PhasePostmortem,
		Severity:            types.SeveritySEV1,
		Status:              types.StatusResolved,
		CreatedAt:           &now,
		StartedAt:           &now,
		ResolvedAt:          &now,
		Summary:             "Postmortem summary.",
		RootCause:           "Database connection pool exhaustion",
		ContributingFactors: []string{"High traffic", "Missing alerts"},
		WhatWentWell:        []string{"Quick detection"},
		WhatWentWrong:       []string{"Slow mitigation"},
		LessonsLearned:      []string{"Add connection pool monitoring"},
	}

	output, err := renderer.RenderPostmortem(incident)
	if err != nil {
		t.Fatalf("RenderPostmortem() failed: %v", err)
	}

	checks := []string{
		"# Postmortem: Test Incident",
		"**Incident ID:** INC-2024-0001",
		"Database connection pool exhaustion",
		"High traffic",
		"Quick detection",
		"Slow mitigation",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Output missing expected content: %q", check)
		}
	}
}

func TestIncidentViewHelpers(t *testing.T) {
	now := time.Now().UTC()
	incident := &types.Incident{
		Timeline: []types.TimelineEvent{
			{EventID: "1", Description: "Confirmed", Confidence: types.ConfidenceConfirmed},
			{EventID: "2", Description: "Suspected", Confidence: types.ConfidenceSuspected},
			{EventID: "3", Description: "Also confirmed", Confidence: types.ConfidenceConfirmed},
		},
		Hypotheses: []types.Hypothesis{
			{HypothesisID: "1", Description: "Investigating", Status: types.HypothesisInvestigating},
			{HypothesisID: "2", Description: "Proposed", Status: types.HypothesisProposed},
			{HypothesisID: "3", Description: "Validated", Status: types.HypothesisValidated},
		},
		ActionItems: []types.ActionItem{
			{ActionID: "1", Description: "Open", Status: types.ActionStatusOpen},
			{ActionID: "2", Description: "In Progress", Status: types.ActionStatusInProgress},
			{ActionID: "3", Description: "Done", Status: types.ActionStatusDone},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	view := newIncidentView(incident)

	// Test ConfirmedFacts
	facts := view.ConfirmedFacts()
	if len(facts) != 2 {
		t.Errorf("ConfirmedFacts() returned %d items, want 2", len(facts))
	}

	// Test ActiveHypotheses
	active := view.ActiveHypotheses()
	if len(active) != 1 {
		t.Errorf("ActiveHypotheses() returned %d items, want 1", len(active))
	}
	if active[0].Description != "Investigating" {
		t.Errorf("ActiveHypotheses()[0].Description = %q, want %q", active[0].Description, "Investigating")
	}

	// Test ProposedHypotheses
	proposed := view.ProposedHypotheses()
	if len(proposed) != 1 {
		t.Errorf("ProposedHypotheses() returned %d items, want 1", len(proposed))
	}

	// Test InProgressActions
	inProgress := view.InProgressActions()
	if len(inProgress) != 1 {
		t.Errorf("InProgressActions() returned %d items, want 1", len(inProgress))
	}

	// Test formatted timestamps
	if view.FormattedCreatedAt() == "" {
		t.Error("FormattedCreatedAt() returned empty string")
	}
	if view.FormattedUpdatedAt() == "" {
		t.Error("FormattedUpdatedAt() returned empty string")
	}
}

func TestLoadIncidentFromReader(t *testing.T) {
	jsonData := `{
		"incident_id": "INC-TEST",
		"title": "Test",
		"phase": "intra_mortem",
		"severity": "SEV2",
		"created_at": "2024-01-01T00:00:00Z"
	}`

	incident, err := LoadIncidentFromReader(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("LoadIncidentFromReader() failed: %v", err)
	}

	if incident.IncidentID != "INC-TEST" {
		t.Errorf("IncidentID = %q, want %q", incident.IncidentID, "INC-TEST")
	}
	if incident.Phase != types.PhaseIntraMortem {
		t.Errorf("Phase = %q, want %q", incident.Phase, types.PhaseIntraMortem)
	}
	if incident.Severity != types.SeveritySEV2 {
		t.Errorf("Severity = %q, want %q", incident.Severity, types.SeveritySEV2)
	}
}
