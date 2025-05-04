package prompts

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed *.tmpl
var promptFS embed.FS

var (
	SystemPromptTemplate string
	UserPromptTemplate   string
)

// PromptLoader loads and executes prompt templates.
type PromptLoader struct {
	templates *template.Template
}

// NewPromptLoader parses all .tmpl files in the package and returns a loader.
func NewPromptLoader() (*PromptLoader, error) {
	tmpl, err := template.New("").ParseFS(promptFS, "*.tmpl")
	if err != nil {
		return nil, err
	}
	return &PromptLoader{templates: tmpl}, nil
}

// GetSystemPrompt returns the system prompt constant from gaming_prompts.go.
func (pl *PromptLoader) GetSystemPrompt() string {
	return SystemPromptTemplate
}

// GetUserPrompt returns the user prompt constant from gaming_prompts.go.
func (pl *PromptLoader) GetUserPrompt() string {
	return UserPromptTemplate
}

// GetDetailedGamingPrompt executes the DetailedGamingPrompt template with the given data.
func (pl *PromptLoader) GetDetailedGamingPrompt(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := pl.templates.ExecuteTemplate(&buf, "DetailedGamingPrompt", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetImprovementPlanPrompt executes the ImprovementPlanPrompt template with the given data.
func (pl *PromptLoader) GetImprovementPlanPrompt(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := pl.templates.ExecuteTemplate(&buf, "ImprovementPlanPrompt", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
