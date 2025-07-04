package models

// TemplatePhase represents a phase in a job template
type TemplatePhase struct {
	ID          string `json:"id" db:"id"`
	TemplateID  string `json:"templateId" db:"template_id"`
	Name        string `json:"name" db:"name"`
	Order       int    `json:"order" db:"phase_order"`
	Description string `json:"description,omitempty" db:"description"`
}
