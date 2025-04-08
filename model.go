package survey

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FUNCTION MAPPING
var defaultFunctions = map[string]func(params map[string]interface{}) string{}

// QUESTION STRUCT TO HOLD THE QUESTION DATA FROM YAML
type Question struct {
	Prompt          string                 `yaml:"prompt"`
	Name            string                 `yaml:"name"`
	Default         string                 `yaml:"default,omitempty"`
	DefaultFunction string                 `yaml:"default_function,omitempty"`
	DefaultParams   map[string]interface{} `yaml:"default_params,omitempty"`
	Options         []string               `yaml:"options"`
	Kind            string                 `yaml:"kind,omitempty"` // "function" instead of "text"
	MinLength       int                    `yaml:"minLength,omitempty"`
	MaxLength       int                    `yaml:"maxLength,omitempty"`
	Type            string                 `yaml:"type,omitempty"` // Updated field to match the YAML
}

// MODEL HOLDS THE STATE FOR THE TERMINAL UI.
type Text struct {
	Textarea   textarea.Model
	ErrMsg     string
	Quitting   bool
	WindowSize tea.WindowSizeMsg // Track window size for responsive layout
}

// SAVE MODEL HOLDS THE STATE FOR THE SAVE FUNCTIONALITY.
type Save struct {
	content     string
	filepicker  filepicker.Model
	filename    textinput.Model
	selectedDir string
	quitting    bool
	status      string
	saved       bool
	err         error
}
