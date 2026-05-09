package visualize

import (
	"fmt"
	"sort"

	"github.com/grokify/d2vision/generate"
	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

// Confidence colors for timeline events
var confidenceColors = map[types.ConfidenceLevel]string{
	types.ConfidenceConfirmed:   "#d4edda", // light green
	types.ConfidenceLikely:      "#c3e6cb", // lighter green
	types.ConfidenceSuspected:   "#fff3cd", // light yellow
	types.ConfidenceUnconfirmed: "#e9ecef", // light gray
}

// Event source shapes
var sourceShapes = map[types.EventSource]string{
	types.EventSourceHuman:      "rectangle",
	types.EventSourceMonitoring: "hexagon",
	types.EventSourceAutomated:  "diamond",
	types.EventSourceCustomer:   "person",
	types.EventSourceAgent:      "circle",
	types.EventSourceOther:      "cloud",
}

// Timeline generates a D2 diagram showing timeline events with evidence links.
func (v *Visualizer) Timeline(incident *types.Incident) string {
	spec := &generate.DiagramSpec{
		Direction: "down",
	}

	container := v.buildTimelineContainer(incident)
	spec.Containers = append(spec.Containers, container)

	// Add evidence container if there's evidence
	if len(incident.Evidence) > 0 {
		evidenceContainer := v.buildEvidenceContainer(incident)
		spec.Containers = append(spec.Containers, evidenceContainer)
	}

	return v.generator.Generate(spec)
}

func (v *Visualizer) buildTimelineContainer(incident *types.Incident) generate.ContainerSpec {
	container := generate.ContainerSpec{
		ID:        "timeline",
		Label:     "Timeline",
		Direction: "down",
	}

	if len(incident.Timeline) == 0 {
		container.Nodes = append(container.Nodes, generate.NodeSpec{
			ID:    "none",
			Label: "No timeline events",
			Style: &generate.StyleSpec{Fill: "#f5f5f5"},
		})
		return container
	}

	// Sort events by timestamp
	events := make([]types.TimelineEvent, len(incident.Timeline))
	copy(events, incident.Timeline)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	// Create nodes for each event
	for _, event := range events {
		timeStr := event.Timestamp.Format("15:04")
		label := fmt.Sprintf("%s - %s", timeStr, truncate(event.Description, 35))

		shape := sourceShapes[event.Source]
		if shape == "" {
			shape = "rectangle"
		}

		fill := confidenceColors[event.Confidence]
		if fill == "" {
			fill = "#e9ecef"
		}

		node := generate.NodeSpec{
			ID:    event.EventID,
			Label: label,
			Shape: shape,
			Style: &generate.StyleSpec{Fill: fill},
		}
		container.Nodes = append(container.Nodes, node)
	}

	// Add edges connecting sequential events
	for i := 0; i < len(events)-1; i++ {
		edge := generate.EdgeSpec{
			From: events[i].EventID,
			To:   events[i+1].EventID,
		}
		container.Edges = append(container.Edges, edge)
	}

	return container
}

func (v *Visualizer) buildEvidenceContainer(incident *types.Incident) generate.ContainerSpec {
	container := generate.ContainerSpec{
		ID:        "evidence",
		Label:     "Evidence",
		Direction: "right",
	}

	// Evidence type shapes
	typeShapes := map[types.EvidenceType]string{
		types.EvidenceTypeLog:        "document",
		types.EvidenceTypeMetric:     "cylinder",
		types.EvidenceTypeTrace:      "parallelogram",
		types.EvidenceTypeScreenshot: "rectangle",
		types.EvidenceTypeDocument:   "page",
		types.EvidenceTypeOther:      "rectangle",
	}

	for _, ev := range incident.Evidence {
		shape := typeShapes[ev.EvidenceType]
		if shape == "" {
			shape = "rectangle"
		}

		node := generate.NodeSpec{
			ID:    ev.EvidenceID,
			Label: truncate(ev.Description, 30),
			Shape: shape,
			Style: &generate.StyleSpec{Fill: "#e8f4f8"},
		}
		container.Nodes = append(container.Nodes, node)
	}

	// Build a map of evidence to events
	eventToEvidence := make(map[string][]string)
	for _, event := range incident.Timeline {
		for _, eviID := range event.EvidenceIDs {
			eventToEvidence[event.EventID] = append(eventToEvidence[event.EventID], eviID)
		}
	}

	// Add edges from evidence to timeline events
	for eventID, eviIDs := range eventToEvidence {
		for _, eviID := range eviIDs {
			edge := generate.EdgeSpec{
				From:  eviID,
				To:    fmt.Sprintf("timeline.%s", eventID),
				Label: "supports",
				Style: &generate.StyleSpec{Stroke: "#6c757d"},
			}
			container.Edges = append(container.Edges, edge)
		}
	}

	return container
}
