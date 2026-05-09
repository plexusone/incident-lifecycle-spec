# Release Notes - v0.1.0

## Overview

Initial release of **incident-lifecycle-spec**, a unified schema and tooling for incident lifecycle management spanning premortem, intra-mortem, and postmortem phases.

## Highlights

- Unified JSON Schema for incident artifacts across all lifecycle phases
- Go types with type-safe enums for phase, severity, status, and confidence
- CLI tool for rendering incidents to Markdown with phase-aware template selection
- Design principles optimized for AI-agent authoring with human review

## The Problem

Traditional incident tooling treats postmortems as standalone documents written after resolution. This creates gaps:

1. **No forward flow** — Premortem risk hypotheses don't connect to postmortem root causes
2. **Lost context** — Real-time incident understanding (intra-mortem) lives in Slack chaos
3. **AI-unfriendly** — Narrative templates don't work well with agent authoring

## The Solution

A single artifact that evolves through three phases:

```
┌─────────────────────────────────────────────────────────┐
│                  Incident JSON                          │
│            (single evolving artifact)                   │
│                                                         │
│  phase: premortem | intra_mortem | postmortem          │
└─────────────────────────────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
    ┌──────────┐   ┌──────────────┐   ┌──────────┐
    │ Premortem│   │ Intra-Mortem │   │Postmortem│
    │ Markdown │   │   Markdown   │   │ Markdown │
    └──────────┘   └──────────────┘   └──────────┘
```

## New Features

### JSON Schema

Flat, navigable schema with top-level arrays for cross-referencing:

```json
{
  "incident_id": "INC-2024-0042",
  "title": "Payment processing failures",
  "phase": "intra_mortem",
  "severity": "SEV1",
  "timeline": [...],
  "hypotheses": [...],
  "action_items": [...],
  "evidence": [...]
}
```

Key design decisions:

| Decision | Rationale |
|----------|-----------|
| Single schema with `phase` | One artifact evolves over time |
| Top-level arrays | Single-hop cross-references |
| IDs on all entities | Enables linking (`validated_by_event_id`) |
| snake_case fields | More readable in JSON |

### Go Types

Type-safe structs matching the schema:

```go
import "github.com/plexusone/incident-lifecycle-spec/pkg/types"

incident := types.Incident{
    IncidentID: "INC-2024-0042",
    Phase:      types.PhaseIntraMortem,
    Severity:   types.SeveritySEV1,
    Status:     types.StatusMitigating,
}
```

### Markdown Renderer

Embedded templates with view helpers:

```go
import "github.com/plexusone/incident-lifecycle-spec/pkg/render"

renderer, _ := render.New()
incident, _ := render.LoadIncident("incident.json")
output, _ := renderer.RenderIntraMortem(incident)
```

View helpers filter by status:

- `ConfirmedFacts()` — Timeline events with confirmed confidence
- `ActiveHypotheses()` — Hypotheses being investigated
- `InProgressActions()` — Action items in progress

### CLI

```bash
# Auto-detect template based on phase
ilspec render incident.json

# Write to file
ilspec render incident.json -o update.md

# Use custom templates
ilspec render incident.json --template-dir ./my-templates
```

## Design Principles

This release establishes core design principles documented in [DESIGN_PRINCIPLES.md](DESIGN_PRINCIPLES.md):

1. **Two levels of nesting maximum** — Like OpenAPI's flat structure
2. **Each object locally comprehensible** — No implicit dependencies
3. **Cross-references shallow** — Single hop, predictable targets
4. **JSON for review, Markdown for scanning** — Two consumption modes
5. **Manual-compatible but agent-optimized** — Works for both

## Validation

Schema validates with [schemakit](https://github.com/grokify/schemakit):

```bash
schemakit lint --property-case snake_case schema/incident.schema.json
```

## What's Next

- Premortem template and examples
- `ilspec validate` command for schema validation
- Hypothesis lifecycle tracking across phases
- Integration with schemakit's planned `navigable` profile

## Links

- [README](README.md)
- [DESIGN_PRINCIPLES](DESIGN_PRINCIPLES.md)
- [CHANGELOG](CHANGELOG.md)
- [Example: Intra-Mortem](examples/intra-mortem-example.json)
