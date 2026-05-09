package visualize

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

func TestHypothesisLifecycle(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		Hypotheses: []types.Hypothesis{
			{
				HypothesisID: "hyp-001",
				Description:  "Database connection pool exhaustion",
				Status:       types.HypothesisProposed,
				Confidence:   0.5,
			},
			{
				HypothesisID: "hyp-002",
				Description:  "Network partition",
				Status:       types.HypothesisInvalidated,
				Confidence:   0.0,
			},
		},
	}

	vis := NewWithoutRenderer()
	d2Code := vis.HypothesisLifecycle(incident)

	if !strings.Contains(d2Code, "hyp-001") {
		t.Error("D2 code should contain hyp-001")
	}
	if !strings.Contains(d2Code, "hyp-002") {
		t.Error("D2 code should contain hyp-002")
	}
	if !strings.Contains(d2Code, "Proposed") {
		t.Error("D2 code should contain Proposed status")
	}
	if !strings.Contains(d2Code, "Invalidated") {
		t.Error("D2 code should contain Invalidated status")
	}
}

func TestTimeline(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		Timeline: []types.TimelineEvent{
			{
				EventID:     "evt-001",
				Timestamp:   now,
				Description: "First alert",
				Source:      types.EventSourceMonitoring,
				Confidence:  types.ConfidenceConfirmed,
			},
			{
				EventID:     "evt-002",
				Timestamp:   now.Add(5 * time.Minute),
				Description: "Investigation started",
				Source:      types.EventSourceHuman,
				Confidence:  types.ConfidenceConfirmed,
			},
		},
	}

	vis := NewWithoutRenderer()
	d2Code := vis.Timeline(incident)

	if !strings.Contains(d2Code, "evt-001") {
		t.Error("D2 code should contain evt-001")
	}
	if !strings.Contains(d2Code, "evt-002") {
		t.Error("D2 code should contain evt-002")
	}
	if !strings.Contains(d2Code, `"evt-001" -> "evt-002"`) {
		t.Error("D2 code should contain edge from evt-001 to evt-002")
	}
}

func TestGenerate(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
	}

	vis := NewWithoutRenderer()

	tests := []struct {
		name        string
		diagramType DiagramType
		wantErr     bool
	}{
		{"hypothesis", DiagramHypothesis, false},
		{"timeline", DiagramTimeline, false},
		{"all", DiagramAll, false},
		{"invalid", DiagramType("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := vis.Generate(incident, tt.diagramType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenderD2(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Test incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
		Hypotheses: []types.Hypothesis{
			{
				HypothesisID: "hyp-001",
				Description:  "Test hypothesis",
				Status:       types.HypothesisInvestigating,
				Confidence:   0.7,
			},
		},
	}

	vis := NewWithoutRenderer()
	output, err := vis.Render(context.Background(), incident, DiagramHypothesis, FormatD2)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if len(output) == 0 {
		t.Error("Render() returned empty output")
	}
}

func TestNewWithRenderer(t *testing.T) {
	vis, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if vis == nil {
		t.Fatal("New() returned nil")
	}
}

func TestEmptyIncident(t *testing.T) {
	now := time.Now()
	incident := &types.Incident{
		IncidentID: "INC-001",
		Title:      "Empty incident",
		Phase:      types.PhaseIntraMortem,
		Severity:   types.SeveritySEV1,
		CreatedAt:  &now,
	}

	vis := NewWithoutRenderer()

	// Should not panic with empty hypotheses
	d2Code := vis.HypothesisLifecycle(incident)
	if !strings.Contains(d2Code, "No hypotheses") {
		t.Error("Should show 'No hypotheses' message")
	}

	// Should not panic with empty timeline
	d2Code = vis.Timeline(incident)
	if !strings.Contains(d2Code, "No timeline events") {
		t.Error("Should show 'No timeline events' message")
	}
}
