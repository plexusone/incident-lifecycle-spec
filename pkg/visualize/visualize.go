// Package visualize generates D2 diagrams from incident data.
package visualize

import (
	"context"
	"fmt"

	"github.com/grokify/d2vision/generate"
	"github.com/grokify/d2vision/render"
	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

// DiagramType specifies which diagram to generate.
type DiagramType string

const (
	DiagramHypothesis DiagramType = "hypothesis"
	DiagramTimeline   DiagramType = "timeline"
	DiagramAll        DiagramType = "all"
)

// Format specifies the output format.
type Format string

const (
	FormatD2  Format = "d2"
	FormatSVG Format = "svg"
	FormatPNG Format = "png"
)

// Visualizer generates diagrams from incident data.
type Visualizer struct {
	generator *generate.Generator
	renderer  *render.Renderer
}

// New creates a new Visualizer.
func New() (*Visualizer, error) {
	renderer, err := render.New()
	if err != nil {
		return nil, fmt.Errorf("creating renderer: %w", err)
	}

	return &Visualizer{
		generator: generate.NewGenerator(),
		renderer:  renderer,
	}, nil
}

// NewWithoutRenderer creates a Visualizer that only generates D2 code (no SVG/PNG).
func NewWithoutRenderer() *Visualizer {
	return &Visualizer{
		generator: generate.NewGenerator(),
	}
}

// Generate creates a diagram of the specified type.
func (v *Visualizer) Generate(incident *types.Incident, diagramType DiagramType) (string, error) {
	switch diagramType {
	case DiagramHypothesis:
		return v.HypothesisLifecycle(incident), nil
	case DiagramTimeline:
		return v.Timeline(incident), nil
	case DiagramAll:
		return v.All(incident), nil
	default:
		return "", fmt.Errorf("unknown diagram type: %s", diagramType)
	}
}

// Render generates a diagram and renders it to the specified format.
func (v *Visualizer) Render(ctx context.Context, incident *types.Incident, diagramType DiagramType, format Format) ([]byte, error) {
	d2Code, err := v.Generate(incident, diagramType)
	if err != nil {
		return nil, err
	}

	switch format {
	case FormatD2:
		return []byte(d2Code), nil
	case FormatSVG:
		if v.renderer == nil {
			return nil, fmt.Errorf("renderer not initialized; use New() instead of NewWithoutRenderer()")
		}
		return v.renderer.RenderSVG(ctx, d2Code, nil)
	case FormatPNG:
		if v.renderer == nil {
			return nil, fmt.Errorf("renderer not initialized; use New() instead of NewWithoutRenderer()")
		}
		return v.renderer.RenderPNG(ctx, d2Code, nil)
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}

// All generates a combined diagram with all visualization types.
func (v *Visualizer) All(incident *types.Incident) string {
	spec := &generate.DiagramSpec{
		Direction: "down",
	}

	// Add hypothesis section if there are hypotheses
	if len(incident.Hypotheses) > 0 {
		hypContainer := v.buildHypothesisContainer(incident)
		spec.Containers = append(spec.Containers, hypContainer)
	}

	// Add timeline section if there are events
	if len(incident.Timeline) > 0 {
		timelineContainer := v.buildTimelineContainer(incident)
		spec.Containers = append(spec.Containers, timelineContainer)
	}

	return v.generator.Generate(spec)
}
