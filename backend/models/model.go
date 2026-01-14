package models

// Provider represents the provider type enum
type Provider string

const (
	ProviderChat   Provider = "chat"
	ProviderClaude Provider = "claude"
	ProviderCodex  Provider = "codex"
	ProviderGemini Provider = "gemini"
)

// Model represents a model entity
type Model struct {
	ID          int      `json:"id"`
	ModelID     string   `json:"model_id"`
	AssModelIDs []string `json:"ass_model_ids,omitempty"` // Associated model IDs
	Provider    Provider `json:"provider"`
	Enabled     bool     `json:"enabled"`
	UpdatedAt   int64    `json:"updated_at"`
}
