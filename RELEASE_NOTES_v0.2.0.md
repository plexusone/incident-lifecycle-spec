# Release Notes - v0.2.0

## Overview

This release adds the **validate command**, **premortem template**, and **unit tests** to incident-lifecycle-spec. It also includes GoReleaser configuration for cross-platform binary distribution.

## Highlights

- `ilspec validate` command for validating incident JSON files
- Premortem template for proactive failure simulation
- Comprehensive unit tests for types and render packages
- GoReleaser configuration for binary releases

## New Features

### Validate Command

Validate incident JSON files against the schema:

```bash
# Validate an incident file
ilspec validate incident.json

# Quiet mode (exit code only)
ilspec validate incident.json -q
```

Validates:

- Required fields: `incident_id`, `title`, `phase`, `severity`, `created_at`
- Valid enum values for `phase`, `severity`, `status`
- Required fields in `timeline`, `hypotheses`, `action_items`, `evidence` arrays

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

### Unit Tests

Comprehensive test coverage for core packages:

**Types package:**

- JSON round-trip serialization
- Enum value correctness (Phase, Severity, Status, etc.)

**Render package:**

- Renderer initialization
- Template rendering (intra-mortem, postmortem)
- View helper methods (`ConfirmedFacts()`, `ActiveHypotheses()`, etc.)
- `LoadIncidentFromReader()`

Run tests:

```bash
go test -v ./...
```

### GoReleaser Configuration

Cross-platform binary releases:

- **Platforms:** Linux, macOS, Windows
- **Architectures:** amd64, arm64
- **Distribution:** Homebrew tap support
- **Version injection:** Via ldflags

## CLI Summary

```bash
# Render incident to Markdown (auto-generated filename)
ilspec render incident.json

# Render to stdout
ilspec render incident.json -o -

# Validate incident JSON
ilspec validate incident.json

# Version
ilspec version
```

## Links

- [README](README.md)
- [CHANGELOG](CHANGELOG.md)
- [v0.1.0 Release Notes](RELEASE_NOTES_v0.1.0.md)
