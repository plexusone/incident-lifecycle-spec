// Package schema provides embedded JSON Schema and validation.
package schema

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed incident.schema.json
var incidentSchemaJSON []byte

// IncidentSchemaJSON returns the embedded incident schema as bytes.
func IncidentSchemaJSON() []byte {
	return incidentSchemaJSON
}

// Validator validates JSON documents against the incident schema.
type Validator struct {
	schema *jsonschema.Schema
}

// NewValidator creates a new schema validator.
func NewValidator() (*Validator, error) {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("incident.schema.json", bytes.NewReader(incidentSchemaJSON)); err != nil {
		return nil, fmt.Errorf("adding schema resource: %w", err)
	}

	schema, err := compiler.Compile("incident.schema.json")
	if err != nil {
		return nil, fmt.Errorf("compiling schema: %w", err)
	}

	return &Validator{schema: schema}, nil
}

// ValidateBytes validates JSON bytes against the schema.
func (v *Validator) ValidateBytes(data []byte) error {
	var doc interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	return v.Validate(doc)
}

// Validate validates a Go value against the schema.
func (v *Validator) Validate(doc interface{}) error {
	if err := v.schema.Validate(doc); err != nil {
		return err
	}
	return nil
}

// ValidationError represents a schema validation error with details.
type ValidationError struct {
	Path    string
	Message string
}

// ValidateBytesDetailed validates JSON and returns detailed errors.
func (v *Validator) ValidateBytesDetailed(data []byte) ([]ValidationError, error) {
	var doc interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	err := v.schema.Validate(doc)
	if err == nil {
		return nil, nil
	}

	// Extract validation errors
	var errors []ValidationError
	if ve, ok := err.(*jsonschema.ValidationError); ok {
		errors = extractErrors(ve, "")
	} else {
		errors = append(errors, ValidationError{
			Path:    "",
			Message: err.Error(),
		})
	}

	return errors, nil
}

func extractErrors(ve *jsonschema.ValidationError, path string) []ValidationError {
	var errors []ValidationError

	currentPath := path
	if ve.InstanceLocation != "" {
		currentPath = ve.InstanceLocation
	}

	if ve.Message != "" {
		errors = append(errors, ValidationError{
			Path:    currentPath,
			Message: ve.Message,
		})
	}

	for _, cause := range ve.Causes {
		errors = append(errors, extractErrors(cause, currentPath)...)
	}

	return errors
}
