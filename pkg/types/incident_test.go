package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIncidentJSONRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	original := Incident{
		IncidentID:       "INC-2024-0001",
		Title:            "Test incident",
		Phase:            PhaseIntraMortem,
		Severity:         SeveritySEV1,
		Status:           StatusInvestigating,
		CreatedAt:        &now,
		Summary:          "Test summary",
		ServicesAffected: []string{"service-a", "service-b"},
		Timeline: []TimelineEvent{
			{
				EventID:     "evt-001",
				Timestamp:   now,
				Description: "First event",
				Source:      EventSourceMonitoring,
				Confidence:  ConfidenceConfirmed,
			},
		},
		Hypotheses: []Hypothesis{
			{
				HypothesisID: "hyp-001",
				Description:  "Root cause hypothesis",
				Status:       HypothesisInvestigating,
				Confidence:   0.75,
			},
		},
		ActionItems: []ActionItem{
			{
				ActionID:    "act-001",
				Description: "Fix the issue",
				Owner:       "oncall",
				Priority:    PriorityP0,
				Status:      ActionStatusInProgress,
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Incident
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.IncidentID != original.IncidentID {
		t.Errorf("IncidentID mismatch: got %q, want %q", decoded.IncidentID, original.IncidentID)
	}
	if decoded.Phase != original.Phase {
		t.Errorf("Phase mismatch: got %q, want %q", decoded.Phase, original.Phase)
	}
	if decoded.Severity != original.Severity {
		t.Errorf("Severity mismatch: got %q, want %q", decoded.Severity, original.Severity)
	}
	if len(decoded.Timeline) != len(original.Timeline) {
		t.Errorf("Timeline length mismatch: got %d, want %d", len(decoded.Timeline), len(original.Timeline))
	}
	if len(decoded.Hypotheses) != len(original.Hypotheses) {
		t.Errorf("Hypotheses length mismatch: got %d, want %d", len(decoded.Hypotheses), len(original.Hypotheses))
	}
	if len(decoded.ActionItems) != len(original.ActionItems) {
		t.Errorf("ActionItems length mismatch: got %d, want %d", len(decoded.ActionItems), len(original.ActionItems))
	}
}

func TestPhaseValues(t *testing.T) {
	tests := []struct {
		phase Phase
		want  string
	}{
		{PhasePremortem, "premortem"},
		{PhaseIntraMortem, "intra_mortem"},
		{PhasePostmortem, "postmortem"},
	}

	for _, tt := range tests {
		if string(tt.phase) != tt.want {
			t.Errorf("Phase %v: got %q, want %q", tt.phase, string(tt.phase), tt.want)
		}
	}
}

func TestSeverityValues(t *testing.T) {
	tests := []struct {
		severity Severity
		want     string
	}{
		{SeveritySEV0, "SEV0"},
		{SeveritySEV1, "SEV1"},
		{SeveritySEV2, "SEV2"},
		{SeveritySEV3, "SEV3"},
	}

	for _, tt := range tests {
		if string(tt.severity) != tt.want {
			t.Errorf("Severity %v: got %q, want %q", tt.severity, string(tt.severity), tt.want)
		}
	}
}

func TestHypothesisStatusValues(t *testing.T) {
	tests := []struct {
		status HypothesisStatus
		want   string
	}{
		{HypothesisProposed, "proposed"},
		{HypothesisInvestigating, "investigating"},
		{HypothesisValidated, "validated"},
		{HypothesisInvalidated, "invalidated"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("HypothesisStatus %v: got %q, want %q", tt.status, string(tt.status), tt.want)
		}
	}
}

func TestActionStatusValues(t *testing.T) {
	tests := []struct {
		status ActionStatus
		want   string
	}{
		{ActionStatusOpen, "open"},
		{ActionStatusInProgress, "in_progress"},
		{ActionStatusDone, "done"},
		{ActionStatusWontDo, "wont_do"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("ActionStatus %v: got %q, want %q", tt.status, string(tt.status), tt.want)
		}
	}
}
