// Package types defines the incident lifecycle data structures.
package types

import "time"

// Phase represents the lifecycle phase of an incident artifact.
type Phase string

const (
	PhasePremortem   Phase = "premortem"
	PhaseIntraMortem Phase = "intra_mortem"
	PhasePostmortem  Phase = "postmortem"
)

// Severity represents the incident severity level.
type Severity string

const (
	SeveritySEV0 Severity = "SEV0"
	SeveritySEV1 Severity = "SEV1"
	SeveritySEV2 Severity = "SEV2"
	SeveritySEV3 Severity = "SEV3"
)

// Status represents the current incident status.
type Status string

const (
	StatusHypothetical  Status = "hypothetical"
	StatusInvestigating Status = "investigating"
	StatusIdentified    Status = "identified"
	StatusMitigating    Status = "mitigating"
	StatusResolved      Status = "resolved"
	StatusClosed        Status = "closed"
)

// Incident is the unified lifecycle artifact spanning premortem, intra-mortem, and postmortem phases.
type Incident struct {
	IncidentID string   `json:"incident_id"`
	Title      string   `json:"title"`
	Phase      Phase    `json:"phase"`
	Severity   Severity `json:"severity"`
	Status     Status   `json:"status,omitempty"`

	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	DetectedAt  *time.Time `json:"detected_at,omitempty"`
	MitigatedAt *time.Time `json:"mitigated_at,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`

	Summary          string   `json:"summary,omitempty"`
	ServicesAffected []string `json:"services_affected,omitempty"`

	CustomerImpactSummary string      `json:"customer_impact_summary,omitempty"`
	CustomerImpactScope   ImpactScope `json:"customer_impact_scope,omitempty"`
	InternalImpactSummary string      `json:"internal_impact_summary,omitempty"`

	DetectionMethod  DetectionMethod `json:"detection_method,omitempty"`
	DetectionDetails string          `json:"detection_details,omitempty"`

	RootCause           string   `json:"root_cause,omitempty"`
	ContributingFactors []string `json:"contributing_factors,omitempty"`
	PreventionGaps      []string `json:"prevention_gaps,omitempty"`
	ResolutionSummary   string   `json:"resolution_summary,omitempty"`

	WhatWentWell   []string `json:"what_went_well,omitempty"`
	WhatWentWrong  []string `json:"what_went_wrong,omitempty"`
	LessonsLearned []string `json:"lessons_learned,omitempty"`

	RelatedIncidentIDs []string `json:"related_incident_ids,omitempty"`

	Timeline    []TimelineEvent `json:"timeline,omitempty"`
	Hypotheses  []Hypothesis    `json:"hypotheses,omitempty"`
	ActionItems []ActionItem    `json:"action_items,omitempty"`
	Evidence    []Evidence      `json:"evidence,omitempty"`
}

// ImpactScope represents the scope of customer impact.
type ImpactScope string

const (
	ImpactScopeNone        ImpactScope = "none"
	ImpactScopeMinimal     ImpactScope = "minimal"
	ImpactScopePartial     ImpactScope = "partial"
	ImpactScopeSignificant ImpactScope = "significant"
	ImpactScopeTotal       ImpactScope = "total"
)

// DetectionMethod represents how an incident was detected.
type DetectionMethod string

const (
	DetectionMonitoring     DetectionMethod = "monitoring"
	DetectionCustomerReport DetectionMethod = "customer_report"
	DetectionEmployeeReport DetectionMethod = "employee_report"
	DetectionAutomated      DetectionMethod = "automated"
	DetectionOther          DetectionMethod = "other"
)

// TimelineEvent represents a single event in the incident timeline.
type TimelineEvent struct {
	EventID     string          `json:"event_id"`
	Timestamp   time.Time       `json:"timestamp"`
	Description string          `json:"description"`
	Source      EventSource     `json:"source,omitempty"`
	Confidence  ConfidenceLevel `json:"confidence,omitempty"`
	EvidenceIDs []string        `json:"evidence_ids,omitempty"`
}

// EventSource represents the source of a timeline event.
type EventSource string

const (
	EventSourceMonitoring EventSource = "monitoring"
	EventSourceHuman      EventSource = "human"
	EventSourceCustomer   EventSource = "customer"
	EventSourceAgent      EventSource = "agent"
	EventSourceAutomated  EventSource = "automated"
	EventSourceOther      EventSource = "other"
)

// ConfidenceLevel represents confidence in a fact or hypothesis.
type ConfidenceLevel string

const (
	ConfidenceConfirmed   ConfidenceLevel = "confirmed"
	ConfidenceLikely      ConfidenceLevel = "likely"
	ConfidenceSuspected   ConfidenceLevel = "suspected"
	ConfidenceUnconfirmed ConfidenceLevel = "unconfirmed"
)

// Hypothesis represents a hypothesis about cause or risk.
type Hypothesis struct {
	HypothesisID       string           `json:"hypothesis_id"`
	Description        string           `json:"description"`
	Status             HypothesisStatus `json:"status"`
	Confidence         float64          `json:"confidence,omitempty"`
	ValidatedByEventID string           `json:"validated_by_event_id,omitempty"`
	EvidenceIDs        []string         `json:"evidence_ids,omitempty"`
}

// HypothesisStatus represents the status of a hypothesis.
type HypothesisStatus string

const (
	HypothesisProposed      HypothesisStatus = "proposed"
	HypothesisInvestigating HypothesisStatus = "investigating"
	HypothesisValidated     HypothesisStatus = "validated"
	HypothesisInvalidated   HypothesisStatus = "invalidated"
)

// ActionItem represents a task to prevent recurrence or mitigate risk.
type ActionItem struct {
	ActionID            string       `json:"action_id"`
	Description         string       `json:"description"`
	Owner               string       `json:"owner,omitempty"`
	Priority            Priority     `json:"priority"`
	Status              ActionStatus `json:"status"`
	DueDate             string       `json:"due_date,omitempty"`
	RelatedHypothesisID string       `json:"related_hypothesis_id,omitempty"`
}

// Priority represents task priority.
type Priority string

const (
	PriorityP0 Priority = "P0"
	PriorityP1 Priority = "P1"
	PriorityP2 Priority = "P2"
	PriorityP3 Priority = "P3"
)

// ActionStatus represents the status of an action item.
type ActionStatus string

const (
	ActionStatusOpen       ActionStatus = "open"
	ActionStatusInProgress ActionStatus = "in_progress"
	ActionStatusDone       ActionStatus = "done"
	ActionStatusWontDo     ActionStatus = "wont_do"
)

// Evidence represents supporting evidence for the incident.
type Evidence struct {
	EvidenceID   string       `json:"evidence_id"`
	EvidenceType EvidenceType `json:"evidence_type"`
	Description  string       `json:"description"`
	Source       string       `json:"source,omitempty"`
	URL          string       `json:"url,omitempty"`
	CollectedAt  *time.Time   `json:"collected_at,omitempty"`
}

// EvidenceType represents the type of evidence.
type EvidenceType string

const (
	EvidenceTypeLog        EvidenceType = "log"
	EvidenceTypeTrace      EvidenceType = "trace"
	EvidenceTypeMetric     EvidenceType = "metric"
	EvidenceTypeScreenshot EvidenceType = "screenshot"
	EvidenceTypeDocument   EvidenceType = "document"
	EvidenceTypeOther      EvidenceType = "other"
)
