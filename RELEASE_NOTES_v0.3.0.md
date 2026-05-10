# Release Notes - v0.3.0

## Overview

This release adds D2 diagram visualization for incident data, including hypothesis lifecycle tracking and timeline visualization with evidence links.

## Highlights

- `ilspec visualize` command for generating D2/SVG diagrams
- Hypothesis lifecycle diagram showing state transitions
- Timeline diagram with evidence links and event source shapes
- `--diagrams` flag to embed D2 code blocks in rendered Markdown
- schemakit navigable profile integration

## New Features

### Visualize Command

Generate D2 diagrams from incident data:

```bash
# Hypothesis lifecycle diagram (D2 code)
ilspec visualize incident.json --type hypothesis

# Timeline diagram as SVG
ilspec visualize incident.json --type timeline --format svg -o timeline.svg

# All diagrams combined
ilspec visualize incident.json --type all
```

**Diagram Types:**

| Type | Description |
|------|-------------|
| `hypothesis` | Shows hypothesis state transitions (proposed → investigating → validated/invalidated) |
| `timeline` | Shows timeline events with evidence links, uses shapes for event sources |
| `all` | Combined diagram with all visualizations |

**Output Formats:**

| Format | Description |
|--------|-------------|
| `d2` | D2 source code (default) |
| `svg` | Rendered SVG image |
| `png` | Rendered PNG image (requires playwright) |

### Diagrams in Rendered Markdown

Embed D2 code blocks directly in Markdown output:

```bash
ilspec render incident.json --diagrams
```

The D2 code blocks render natively on GitHub and GitLab.

### Hypothesis Lifecycle Visualization

Shows hypothesis state transitions across phases with color coding:

```d2
hypotheses: Hypothesis Lifecycle {
  proposed: Proposed { style.fill: "#e3f2fd" }
  investigating: Investigating { style.fill: "#fff3cd" }
  validated: Validated ✓ { style.fill: "#d4edda" }
  invalidated: Invalidated ✗ { style.fill: "#f8d7da" }
}
```

### Timeline Visualization

Shows timeline events with:

- Event source shapes (hexagon for monitoring, rectangle for human, diamond for automated)
- Confidence level colors (green for confirmed, yellow for suspected)
- Sequential flow arrows
- Evidence links

## API

### pkg/visualize

New package for D2 diagram generation:

```go
import "github.com/plexusone/incident-lifecycle-spec/pkg/visualize"

// Create visualizer (with SVG rendering)
vis, _ := visualize.New()

// Generate D2 code
d2Code, _ := vis.Generate(incident, visualize.DiagramHypothesis)

// Render to SVG
svg, _ := vis.Render(ctx, incident, visualize.DiagramTimeline, visualize.FormatSVG)
```

## schemakit Integration

The design principles from incident-lifecycle-spec have been codified into a new `navigable` profile in schemakit:

```bash
schemakit lint --profile navigable --property-case snake_case schema/incident.schema.json
```

The navigable profile enforces:

- Maximum 2 levels of object nesting
- ID fields in array items for cross-referencing
- Flat structure for human reviewability

## Dependencies

- `github.com/grokify/d2vision` - D2 code generation and SVG rendering

## Links

- [README](README.md)
- [CHANGELOG](CHANGELOG.md)
- [v0.2.0 Release Notes](RELEASE_NOTES_v0.2.0.md)
- [schemakit navigable profile](https://github.com/grokify/schemakit/blob/main/docs/reference/profiles.md)
