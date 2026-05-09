package schema

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}
	if v == nil {
		t.Fatal("NewValidator() returned nil")
	}
}

func TestIncidentSchemaJSON(t *testing.T) {
	data := IncidentSchemaJSON()
	if len(data) == 0 {
		t.Fatal("IncidentSchemaJSON() returned empty data")
	}
}

func TestValidateBytes_ValidIncident(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	validJSON := `{
		"incident_id": "INC-001",
		"title": "Test incident",
		"phase": "intra_mortem",
		"severity": "SEV1",
		"created_at": "2024-01-15T10:00:00Z"
	}`

	err = v.ValidateBytes([]byte(validJSON))
	if err != nil {
		t.Errorf("ValidateBytes() error = %v, want nil", err)
	}
}

func TestValidateBytes_InvalidJSON(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	invalidJSON := `{not valid json`

	err = v.ValidateBytes([]byte(invalidJSON))
	if err == nil {
		t.Error("ValidateBytes() error = nil, want error for invalid JSON")
	}
}

func TestValidateBytes_MissingRequired(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	// Missing incident_id
	invalidJSON := `{
		"title": "Test incident",
		"phase": "intra_mortem",
		"severity": "SEV1",
		"created_at": "2024-01-15T10:00:00Z"
	}`

	err = v.ValidateBytes([]byte(invalidJSON))
	if err == nil {
		t.Error("ValidateBytes() error = nil, want error for missing required field")
	}
}

func TestValidateBytesDetailed_MissingRequired(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	// Missing incident_id
	invalidJSON := `{
		"title": "Test incident",
		"phase": "intra_mortem",
		"severity": "SEV1",
		"created_at": "2024-01-15T10:00:00Z"
	}`

	errors, err := v.ValidateBytesDetailed([]byte(invalidJSON))
	if err != nil {
		t.Fatalf("ValidateBytesDetailed() error = %v", err)
	}
	if len(errors) == 0 {
		t.Error("ValidateBytesDetailed() returned no errors, want errors for missing required field")
	}
}

func TestValidateBytesDetailed_InvalidEnum(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	// Invalid phase value
	invalidJSON := `{
		"incident_id": "INC-001",
		"title": "Test incident",
		"phase": "invalid_phase",
		"severity": "SEV1",
		"created_at": "2024-01-15T10:00:00Z"
	}`

	errors, err := v.ValidateBytesDetailed([]byte(invalidJSON))
	if err != nil {
		t.Fatalf("ValidateBytesDetailed() error = %v", err)
	}
	if len(errors) == 0 {
		t.Error("ValidateBytesDetailed() returned no errors, want errors for invalid enum")
	}
}

func TestValidateBytesDetailed_Valid(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	validJSON := `{
		"incident_id": "INC-001",
		"title": "Test incident",
		"phase": "intra_mortem",
		"severity": "SEV1",
		"created_at": "2024-01-15T10:00:00Z"
	}`

	errors, err := v.ValidateBytesDetailed([]byte(validJSON))
	if err != nil {
		t.Fatalf("ValidateBytesDetailed() error = %v", err)
	}
	if len(errors) != 0 {
		t.Errorf("ValidateBytesDetailed() returned %d errors, want 0: %v", len(errors), errors)
	}
}
