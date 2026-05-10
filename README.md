# Incident Lifecycle Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/incident-lifecycle-spec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/incident-lifecycle-spec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/incident-lifecycle-spec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/incident-lifecycle-spec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/incident-lifecycle-spec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/incident-lifecycle-spec/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/incident-lifecycle-spec
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/incident-lifecycle-spec
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/incident-lifecycle-spec
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/incident-lifecycle-spec
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/incident-lifecycle-spec
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fincident-lifecycle-spec
 [loc-svg]: https://tokei.rs/b1/github/plexusone/incident-lifecycle-spec
 [repo-url]: https://github.com/plexusone/incident-lifecycle-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/incident-lifecycle-spec/blob/main/LICENSE

A unified schema and tooling for incident lifecycle management, spanning premortem, intra-mortem, and postmortem phases.

## Why "Lifecycle" Instead of "Postmortem"

Traditional incident tooling treats postmortems as standalone documents written after resolution. This misses two opportunities:

1. **Premortems** — Proactive failure simulation before incidents occur
2. **Intra-mortems** — Real-time incident understanding during active incidents

This spec models incidents as **a single artifact that evolves over time**, not three separate documents. Information flows forward: a premortem's risk hypotheses become an intra-mortem's investigation targets, which become a postmortem's validated or invalidated causes.

## The Three Phases

### Premortem (Before Incident)

Anticipate failure modes before they happen.

- What could break?
- What signals would we see first?
- What mitigations exist?
- What assumptions are risky?

### Intra-Mortem (During Incident)

Maintain shared reality while things are unstable.

- What do we know *right now*?
- What is confirmed vs suspected?
- What hypotheses are being tested?
- What actions are in progress?

This is the most underserved phase. Most teams rely on Slack chaos or fragmented war room notes. A structured format enables clearer communication and handoffs.

### Postmortem (After Resolution)

Reconstruct truth and extract learning.

- What actually happened?
- Why did it happen?
- Why wasn't it prevented?
- What should change?

## Architecture: One Truth, Many Views

```
┌────────────────────────────────────────────────────┐
│                  Incident JSON                     │
│            (single evolving artifact)              │
│                                                    │
│  phase: premortem | intra_mortem | postmortem      │
└────────────────────────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
    ┌──────────┐   ┌──────────────┐   ┌──────────┐
    │ Premortem│   │ Intra-Mortem │   │Postmortem│
    │ Markdown │   │   Markdown   │   │ Markdown │
    └──────────┘   └──────────────┘   └──────────┘
```

- **JSON** is the source of truth — structured, machine-readable, agent-fillable
- **Markdown** is for human consumption — rendered views optimized for scanning

The JSON structure and Markdown structure don't need to match. Renderers can flatten, reorder, or summarize as needed.

## Installation

```bash
go install github.com/plexusone/incident-lifecycle-spec/cmd/ilspec@latest
```

Or build from source:

```bash
git clone https://github.com/plexusone/incident-lifecycle-spec
cd incident-lifecycle-spec
go build -o ilspec ./cmd/ilspec
```

## Usage

### Initialize a new incident

```bash
# Create intra-mortem incident (default)
ilspec init --title "Database outage" --severity SEV1

# Create premortem for failure simulation
ilspec init -p premortem -s SEV2 -t "Failover risk analysis"

# Create postmortem
ilspec init -p postmortem -s SEV0 -t "Auth service outage" -o postmortem.json

# Output to stdout
ilspec init -p intra_mortem -s SEV1 -t "API failures" -o -
```

### Render an incident to Markdown

```bash
# Auto-generates output filename based on phase
# incident.json (intra_mortem) → incident-update.md
ilspec render incident.json

# Explicit output filename
ilspec render incident.json -o custom-name.md

# Output to stdout
ilspec render incident.json -o -

# Use specific template
ilspec render incident.json --template postmortem.md.tmpl

# Use custom templates directory
ilspec render incident.json --template-dir ./my-templates
```

### Validate an incident file

```bash
# Full validation (schema + Go type validation)
ilspec validate incident.json

# Schema validation only
ilspec validate --schema-only incident.json

# Quiet mode (no output on success)
ilspec validate incident.json -q
```

### Validate schema definition (using schemakit)

```bash
schemakit lint --property-case snake_case schema/incident.schema.json

# Use navigable profile for human-reviewable schema checks
schemakit lint --profile navigable --property-case snake_case schema/incident.schema.json
```

### Visualize incident data

Generate D2 diagrams showing hypothesis lifecycle and timeline:

```bash
# Generate hypothesis lifecycle diagram (D2 code to stdout)
ilspec visualize incident.json --type hypothesis

# Generate timeline diagram as SVG
ilspec visualize incident.json --type timeline --format svg -o timeline.svg

# Generate all diagrams
ilspec visualize incident.json --type all -o diagrams.d2

# Embed diagrams in rendered Markdown
ilspec render incident.json --diagrams
```

## Schema Overview

The schema uses a flat structure optimized for navigability:

```json
{
  "incident_id": "INC-2024-0042",
  "title": "Payment processing failures",
  "phase": "intra_mortem",
  "severity": "SEV1",
  "status": "mitigating",

  "summary": "...",
  "customer_impact_summary": "...",
  "services_affected": ["payment-gateway", "checkout"],

  "timeline": [...],
  "hypotheses": [...],
  "action_items": [...],
  "evidence": [...]
}
```

Key design decisions:

| Decision | Rationale |
|----------|-----------|
| Single schema with `phase` field | One artifact evolves over time |
| Top-level arrays | Flat structure enables single-hop cross-references |
| IDs on all entities | Enables linking (e.g., `validated_by_event_id`) |
| snake_case fields | More readable in JSON |

See [DESIGN_PRINCIPLES.md](DESIGN_PRINCIPLES.md) for full design rationale.

## Why This Design Works for AI Agents

### Determinism

No free-form sections required for core fields. Agents populate structured data, not prose.

### Tool Chaining

Works well with:

- Log analysis tools
- Tracing systems (Datadog, Honeycomb)
- Incident detectors (PagerDuty, Opsgenie)
- Ticketing systems (Jira, Linear)

### Multi-Agent Workflows

Responsibilities can be split:

- Agent A: Timeline reconstruction from logs
- Agent B: Root cause hypothesis generation
- Agent C: Action item generation from validated causes

### Human Review

Markdown export stays readable and consistent. Humans review rendered output, not raw JSON.

## Hypothesis Lifecycle

A key feature is tracking hypotheses across phases:

```
Premortem: hypothesis proposed (risk scenario)
     ↓
Intra-mortem: hypothesis investigating → validated/invalidated
     ↓
Postmortem: hypothesis becomes confirmed root cause or rejected
```

Each hypothesis includes:

- `status`: proposed → investigating → validated/invalidated
- `confidence`: 0-1 score
- `validated_by_event_id`: links to timeline event that confirmed/rejected it
- `evidence_ids`: supporting evidence

## Future Enhancements

### Confidence Scoring

```json
"analysis_confidence": {
  "root_cause": 0.85,
  "timeline_completeness": 0.92
}
```

### Evidence Linking

```json
"evidence": [{
  "evidence_id": "evi-001",
  "evidence_type": "trace",
  "source": "datadog",
  "url": "https://app.datadoghq.com/apm/trace/xyz"
}]
```

### Blameless Enforcement

Schema-level guardrails:

- No "human fault attribution" fields
- Root cause must be system/process-oriented
- Required contributing factors analysis

## Project Structure

```
incident-lifecycle-spec/
├── schema/
│   └── incident.schema.json    # JSON Schema (source of truth)
├── pkg/
│   ├── types/
│   │   └── incident.go         # Go types
│   ├── render/
│   │   ├── render.go           # Renderer
│   │   └── templates/
│   │       ├── intra-mortem.md.tmpl
│   │       ├── postmortem.md.tmpl
│   │       └── premortem.md.tmpl
│   └── schema/
│       └── schema.go           # Embedded schema validator
├── cmd/ilspec/
│   ├── main.go
│   ├── init.go
│   ├── render.go
│   └── validate.go
├── examples/
│   ├── intra-mortem-example.json
│   ├── premortem-example.json
│   └── postmortem-example.json
├── DESIGN_PRINCIPLES.md
└── README.md
```

## Related Projects

- [schemakit](https://github.com/grokify/schemakit) — JSON Schema linter used for validation

## License

MIT
