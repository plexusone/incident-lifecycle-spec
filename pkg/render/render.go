// Package render provides Markdown rendering for incident artifacts.
package render

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/plexusone/incident-lifecycle-spec/pkg/types"
)

//go:embed templates/*.tmpl
var embeddedTemplates embed.FS

// Renderer renders incident artifacts to Markdown.
type Renderer struct {
	templates *template.Template
}

// New creates a new Renderer with embedded templates.
func New() (*Renderer, error) {
	funcMap := template.FuncMap{
		"mulf": func(a, b float64) float64 { return a * b },
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(embeddedTemplates, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing embedded templates: %w", err)
	}

	return &Renderer{templates: tmpl}, nil
}

// NewFromDir creates a new Renderer loading templates from a directory.
func NewFromDir(dir string) (*Renderer, error) {
	funcMap := template.FuncMap{
		"mulf": func(a, b float64) float64 { return a * b },
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(dir + "/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates from %s: %w", dir, err)
	}

	return &Renderer{templates: tmpl}, nil
}

// LoadIncident loads an incident from a JSON file.
func LoadIncident(path string) (*types.Incident, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	return LoadIncidentFromReader(f)
}

// LoadIncidentFromReader loads an incident from a reader.
func LoadIncidentFromReader(r io.Reader) (*types.Incident, error) {
	var incident types.Incident
	if err := json.NewDecoder(r).Decode(&incident); err != nil {
		return nil, fmt.Errorf("decoding incident: %w", err)
	}
	return &incident, nil
}

// RenderIntraMortem renders an incident as an intra-mortem update.
func (r *Renderer) RenderIntraMortem(incident *types.Incident) (string, error) {
	return r.render("intra-mortem.md.tmpl", incident)
}

// RenderPostmortem renders an incident as a postmortem report.
func (r *Renderer) RenderPostmortem(incident *types.Incident) (string, error) {
	return r.render("postmortem.md.tmpl", incident)
}

// Render renders an incident using the specified template name.
func (r *Renderer) Render(templateName string, incident *types.Incident) (string, error) {
	return r.render(templateName, incident)
}

func (r *Renderer) render(templateName string, incident *types.Incident) (string, error) {
	// Create a view model with helper methods
	view := newIncidentView(incident)

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, templateName, view); err != nil {
		return "", fmt.Errorf("executing template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// incidentView wraps an Incident with helper methods for templates.
type incidentView struct {
	*types.Incident
}

func newIncidentView(i *types.Incident) *incidentView {
	return &incidentView{Incident: i}
}

// ConfirmedFacts returns timeline events with confirmed confidence.
func (v *incidentView) ConfirmedFacts() []string {
	var facts []string
	for _, e := range v.Timeline {
		if e.Confidence == types.ConfidenceConfirmed {
			facts = append(facts, e.Description)
		}
	}
	return facts
}

// ActiveHypotheses returns hypotheses that are being investigated.
func (v *incidentView) ActiveHypotheses() []types.Hypothesis {
	var active []types.Hypothesis
	for _, h := range v.Hypotheses {
		if h.Status == types.HypothesisInvestigating {
			active = append(active, h)
		}
	}
	return active
}

// ProposedHypotheses returns hypotheses that are proposed but not yet investigated.
func (v *incidentView) ProposedHypotheses() []types.Hypothesis {
	var proposed []types.Hypothesis
	for _, h := range v.Hypotheses {
		if h.Status == types.HypothesisProposed {
			proposed = append(proposed, h)
		}
	}
	return proposed
}

// InProgressActions returns action items that are in progress.
func (v *incidentView) InProgressActions() []types.ActionItem {
	var inProgress []types.ActionItem
	for _, a := range v.ActionItems {
		if a.Status == types.ActionStatusInProgress {
			inProgress = append(inProgress, a)
		}
	}
	return inProgress
}

// FormattedCreatedAt returns a formatted created_at timestamp.
func (v *incidentView) FormattedCreatedAt() string {
	if v.CreatedAt == nil {
		return ""
	}
	return v.CreatedAt.Format("2006-01-02 15:04:05 UTC")
}

// FormattedUpdatedAt returns a formatted updated_at timestamp.
func (v *incidentView) FormattedUpdatedAt() string {
	if v.UpdatedAt == nil {
		return ""
	}
	return v.UpdatedAt.Format("2006-01-02 15:04:05 UTC")
}

// FormattedStartedAt returns a formatted started_at timestamp.
func (v *incidentView) FormattedStartedAt() string {
	if v.StartedAt == nil {
		return ""
	}
	return v.StartedAt.Format("2006-01-02 15:04:05 UTC")
}

// FormattedResolvedAt returns a formatted resolved_at timestamp.
func (v *incidentView) FormattedResolvedAt() string {
	if v.ResolvedAt == nil {
		return ""
	}
	return v.ResolvedAt.Format("2006-01-02 15:04:05 UTC")
}
