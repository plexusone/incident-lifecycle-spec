# Design Principles

Design principles for the incident-lifecycle-spec JSON Schema.

## 1. Two Levels of Nesting Maximum

Like OpenAPI's `components/schemas/`, nesting should not exceed two levels. This keeps the schema navigable—you can drill into a section and see most of its structure without scrolling or expanding multiple levels.

**Good:**

```json
{
  "root_cause_analysis": {
    "primary_cause": "",
    "contributing_factors": []
  }
}
```

**Avoid:**

```json
{
  "analysis": {
    "root_cause": {
      "details": {
        "primary": ""
      }
    }
  }
}
```

## 2. Each Object Locally Comprehensible

Any object in the schema should be understandable when viewed in isolation. A reader should not need to hold context from parent objects or other sections to understand what they're looking at.

This means:

- Self-contained objects with clear field names
- No implicit dependencies on sibling data
- Field names that are self-documenting

## 3. Cross-References Allowed but Shallow

References between objects are acceptable when they are:

- **Single hop** — never reference chains (A → B → C)
- **Predictable targets** — always to a known top-level collection
- **Tool-resolvable** — tooling can inline them for display

**Good examples:**

- `action_items[].related_incident_id` → external incident
- `timeline[].evidence_ref` → entry in top-level `evidence` array
- `hypotheses[].validated_by` → a `timeline` event ID

**Avoid:**

- References that require following multiple links
- References to arbitrary nested locations
- Circular reference patterns

## 4. JSON for Object-Level Review, Markdown for Linear Scanning

The schema has two consumption modes:

| Format | Purpose | User Action |
|--------|---------|-------------|
| JSON | Detailed review and editing | Navigate to specific object, view/edit in isolation |
| Markdown | Quick scanning and sharing | Read linearly, get overview |

The JSON structure and Markdown structure do not need to match exactly. The renderer can flatten, reorder, or summarize as needed for readability.

## 5. Manual-Compatible but Agent-Optimized

The schema must work for manual editing (human opens JSON in editor), but is optimized for agent authoring with human-in-the-loop review.

This means:

- Field names a human can understand without documentation
- No clever encodings or abbreviations
- Incremental population (agents can fill fields over time)
- Errors in one section don't invalidate the whole document

## 6. Flat Over Nested When Equivalent

When two structures convey the same information, prefer the flatter one.

**Prefer:**

```json
{
  "root_cause": "",
  "contributing_factors": []
}
```

**Over:**

```json
{
  "root_cause_analysis": {
    "primary_cause": "",
    "contributing_factors": []
  }
}
```

Exception: grouping is acceptable when the nested object is locally comprehensible and the grouping aids navigation (like OpenAPI's `components`).

## Future: Linting Rules

These principles should eventually be enforced via automated linting using [schemakit](https://github.com/grokify/schemakit).

### Proposed Profile: `navigable`

Schemakit supports profiles for different use cases. These principles could become a `navigable` profile optimized for human-editable, agent-authored schemas.

| Principle | Lint Code | Severity | Description |
|-----------|-----------|----------|-------------|
| Two levels max | `max-nesting-depth` | error | Nesting exceeds two levels |
| Shallow refs | `reference-chain-depth` | error | Reference requires multiple hops to resolve |
| Self-documenting | `abbreviated-field-name` | warning | Field name uses unclear abbreviation |
| Predictable refs | `unpredictable-ref-target` | error | `$ref` points outside known top-level collections |

### Candidate Implementation

```go
const ProfileNavigable Profile = "navigable"

type NavigableConfig struct {
    MaxNestingDepth     int      // default: 2
    MaxReferenceHops    int      // default: 1
    AllowedRefTargets   []string // e.g., ["#/$defs/", "#/components/"]
    AbbreviationDenylist []string // e.g., ["cfg", "mgr", "impl"]
}
```

### Analysis Modes

Schemakit can support two analysis modes:

| Mode | Speed | Deterministic | Catches |
|------|-------|---------------|---------|
| Static | Fast | Yes | Structural issues (nesting, refs) |
| LLM | Slower | No | Semantic issues (clarity, comprehensibility) |

Results should indicate which analysis was performed:

```json
{
  "issues": [...],
  "analysis": {
    "static": true,
    "llm": true,
    "llm_model": "claude-sonnet-4-20250514"
  }
}
```

When LLM analysis is skipped (offline, cost, speed), the result should clearly indicate no semantic analysis was done rather than implying the schema passed.

### Static Checks

| Check | Code | Description |
|-------|------|-------------|
| Nesting depth | `max-nesting-depth` | Count levels, error if > 2 |
| Reference hops | `reference-chain-depth` | Follow `$ref`, error if > 1 hop |
| Predictable refs | `unpredictable-ref-target` | Refs must point to allowed prefixes |
| Abbreviated names | `abbreviated-field-name` | Denylist check (e.g., `cfg`, `mgr`) |

### LLM Checks

| Check | Code | Prompt Strategy |
|-------|------|-----------------|
| Field clarity | `unclear-field-name` | "Is this field name self-documenting without reading other fields?" |
| Local comprehensibility | `requires-external-context` | "Can this object be understood in isolation, or does it require context from parent/sibling objects?" |
| Naming consistency | `inconsistent-naming` | "Are field names consistent in style and terminology across the schema?" |

LLM checks should return:

- `pass` / `warn` / `fail`
- Confidence score (0-1)
- Explanation for human review

### Example Output

```json
{
  "code": "requires-external-context",
  "severity": "warning",
  "path": "$.components.schemas.TimelineEvent",
  "message": "Object may require external context to understand",
  "analysis_type": "llm",
  "confidence": 0.72,
  "explanation": "The 'ref' field references timeline entries by index, requiring the reader to understand the timeline array structure."
}
```

See: [schemakit profiles](https://github.com/grokify/schemakit/blob/main/docs/reference/profiles.md)
