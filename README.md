# incident-lifecycle-spec

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

### Render an incident to Markdown

```bash
# Auto-detect template based on phase
ilspec render incident.json

# Write to file
ilspec render incident.json -o update.md

# Use specific template
ilspec render incident.json --template postmortem.md.tmpl

# Use custom templates directory
ilspec render incident.json --template-dir ./my-templates
```

### Validate schema (using schemakit)

```bash
schemakit lint --property-case snake_case schema/incident.schema.json
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
│   └── render/
│       ├── render.go           # Renderer
│       └── templates/
│           ├── intra-mortem.md.tmpl
│           └── postmortem.md.tmpl
├── cmd/ilspec/
│   ├── main.go
│   └── render.go
├── examples/
│   └── intra-mortem-example.json
├── DESIGN_PRINCIPLES.md
└── README.md
```

## Related Projects

- [schemakit](https://github.com/grokify/schemakit) — JSON Schema linter used for validation

## License

MIT
