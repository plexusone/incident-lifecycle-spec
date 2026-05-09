package visualize

import (
	"fmt"

	"github.com/grokify/d2vision/generate"
	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

// Status colors for hypothesis states
var hypothesisColors = map[types.HypothesisStatus]string{
	types.HypothesisProposed:      "#e3f2fd", // light blue
	types.HypothesisInvestigating: "#fff3cd", // light yellow
	types.HypothesisValidated:     "#d4edda", // light green
	types.HypothesisInvalidated:   "#f8d7da", // light red
}

// HypothesisLifecycle generates a D2 diagram showing hypothesis state transitions.
func (v *Visualizer) HypothesisLifecycle(incident *types.Incident) string {
	spec := &generate.DiagramSpec{
		Direction: "right",
	}

	container := v.buildHypothesisContainer(incident)
	spec.Containers = append(spec.Containers, container)

	return v.generator.Generate(spec)
}

func (v *Visualizer) buildHypothesisContainer(incident *types.Incident) generate.ContainerSpec {
	container := generate.ContainerSpec{
		ID:        "hypotheses",
		Label:     "Hypothesis Lifecycle",
		Direction: "down",
	}

	if len(incident.Hypotheses) == 0 {
		container.Nodes = append(container.Nodes, generate.NodeSpec{
			ID:    "none",
			Label: "No hypotheses",
			Style: &generate.StyleSpec{Fill: "#f5f5f5"},
		})
		return container
	}

	// Group hypotheses by status for layout
	statusGroups := make(map[types.HypothesisStatus][]types.Hypothesis)
	for _, h := range incident.Hypotheses {
		statusGroups[h.Status] = append(statusGroups[h.Status], h)
	}

	// Create status containers in order
	statusOrder := []types.HypothesisStatus{
		types.HypothesisProposed,
		types.HypothesisInvestigating,
		types.HypothesisValidated,
		types.HypothesisInvalidated,
	}

	for _, status := range statusOrder {
		hyps, ok := statusGroups[status]
		if !ok || len(hyps) == 0 {
			continue
		}

		statusContainer := generate.ContainerSpec{
			ID:        string(status),
			Label:     formatStatus(status),
			Direction: "down",
			Style:     &generate.StyleSpec{Fill: hypothesisColors[status]},
		}

		for _, h := range hyps {
			label := truncate(h.Description, 40)
			if h.Confidence > 0 {
				label = fmt.Sprintf("%s\n(%.0f%%)", label, h.Confidence*100)
			}

			node := generate.NodeSpec{
				ID:    h.HypothesisID,
				Label: label,
			}
			statusContainer.Nodes = append(statusContainer.Nodes, node)
		}

		container.Containers = append(container.Containers, statusContainer)
	}

	// Add edges for validated_by_event_id links
	for _, h := range incident.Hypotheses {
		if h.ValidatedByEventID != "" {
			edge := generate.EdgeSpec{
				From:  h.HypothesisID,
				To:    fmt.Sprintf("timeline.%s", h.ValidatedByEventID),
				Label: "validated by",
				Style: &generate.StyleSpec{Stroke: "#28a745"},
			}
			container.Edges = append(container.Edges, edge)
		}
	}

	return container
}

func formatStatus(status types.HypothesisStatus) string {
	switch status {
	case types.HypothesisProposed:
		return "Proposed"
	case types.HypothesisInvestigating:
		return "Investigating"
	case types.HypothesisValidated:
		return "Validated ✓"
	case types.HypothesisInvalidated:
		return "Invalidated ✗"
	default:
		return string(status)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
