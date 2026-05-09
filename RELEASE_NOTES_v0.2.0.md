# Release Notes - v0.2.0

## Overview

This release adds scaffolding, JSON Schema validation, premortem templates, comprehensive tests, and complete examples for all lifecycle phases.

## Highlights

- `ilspec init` command for scaffolding new incidents with phase-aware templates
- JSON Schema runtime validation with embedded schema
- `ilspec validate` command for validating incident JSON files
- Premortem template for proactive failure simulation
- Complete example incidents for all phases (premortem, intra-mortem, postmortem)

## New Features

### Init Command

Scaffold new incident JSON files with phase-appropriate defaults:

```bash
# Create intra-mortem incident
ilspec init --phase intra_mortem --severity SEV1 --title "Database outage"

# Create premortem for failure simulation
ilspec init -p premortem -s SEV2 -t "Failover risk analysis" -o risk.json

# Output to stdout
ilspec init -p postmortem -s SEV0 -t "Auth service outage" -o -
```

Auto-generates:

- Incident ID based on timestamp (e.g., `INC-2024-0515-143022`)
- Phase-appropriate status (`hypothetical` for premortem, `investigating` for intra-mortem)
- Placeholder content for required sections
- Default filename based on phase (e.g., `premortem-2024-05-15.json`)

### JSON Schema Validation

Validate incidents against the embedded JSON Schema:

```bash
# Full validation (schema + Go type validation)
ilspec validate incident.json

# Schema validation only
ilspec validate --schema-only incident.json

# Quiet mode (exit code only)
ilspec validate -q incident.json
```

Validates:

- Required fields: `incident_id`, `title`, `phase`, `severity`, `created_at`
- Valid enum values for `phase`, `severity`, `status`
- Required fields in `timeline`, `hypotheses`, `action_items`, `evidence` arrays
- JSON structure against schema

Example output:

```
✓ incident.json is valid
  Phase: intra_mortem
  Severity: SEV1
  Status: mitigating
```

### Premortem Template

New template for proactive failure simulation before incidents occur:

```bash
ilspec render premortem.json
# Outputs: premortem-premortem.md
```

The premortem template includes:

- **Failure Scenarios** — Potential failure modes with likelihood scores
- **Early Warning Signals** — What signals would indicate failure is occurring
- **Prevention Gaps** — Why these failures might not be prevented
- **Recommended Mitigations** — Actions to reduce risk
- **Risky Assumptions** — Assumptions that could prove false

### Complete Examples

All lifecycle phases now have complete example JSON files:

| Phase | Example | Description |
|-------|---------|-------------|
| Premortem | `examples/premortem-example.json` | Database failover risk analysis |
| Intra-Mortem | `examples/intra-mortem-example.json` | Payment service outage (active) |
| Postmortem | `examples/postmortem-example.json` | Auth service certificate expiration |

## API

### pkg/schema

New package for JSON Schema validation:

```go
import "github.com/plexusone/incident-lifecycle-spec/pkg/schema"

// Create validator with embedded schema
validator, err := schema.NewValidator()

// Simple validation
err = validator.ValidateBytes(jsonData)

// Detailed validation with paths
errors, err := validator.ValidateBytesDetailed(jsonData)
for _, e := range errors {
    fmt.Printf("%s: %s\n", e.Path, e.Message)
}

// Access raw schema
schemaJSON := schema.IncidentSchemaJSON()
```

## Testing

Comprehensive test coverage for all packages:

**CLI tests:**

- `cmd/ilspec/validate_test.go` — Validation logic tests (required fields, enums, nested structures)

**Package tests:**

- `pkg/schema/schema_test.go` — Schema validator tests
- `pkg/types/types_test.go` — JSON round-trip, enum values
- `pkg/render/render_test.go` — Template rendering, view helpers

Run tests:

```bash
go test -v ./...
```

## Build

GoReleaser configuration supports cross-platform binaries:

- **Platforms:** Linux, macOS, Windows
- **Architectures:** amd64, arm64
- **Distribution:** Homebrew tap support
- **Version injection:** Via ldflags

## CLI Summary

```bash
# Initialize new incident file
ilspec init -p intra_mortem -s SEV1 -t "Incident title"

# Render incident to Markdown (auto-generated filename)
ilspec render incident.json

# Render to stdout
ilspec render incident.json -o -

# Validate incident JSON
ilspec validate incident.json

# Schema-only validation
ilspec validate --schema-only incident.json

# Version
ilspec version
```

## What's Next

- Hypothesis lifecycle tracking visualization
- Timeline rendering with evidence links
- Integration with schemakit's `navigable` profile
- GitHub Actions for automated releases

## Links

- [README](README.md)
- [DESIGN_PRINCIPLES](DESIGN_PRINCIPLES.md)
- [CHANGELOG](CHANGELOG.md)
- [v0.1.0 Release Notes](RELEASE_NOTES_v0.1.0.md)
- [Example: Premortem](examples/premortem-example.json)
- [Example: Intra-Mortem](examples/intra-mortem-example.json)
- [Example: Postmortem](examples/postmortem-example.json)
