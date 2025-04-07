package survey

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

// INITIALText INITIALIZES THE Text WITH EITHER PROVIDED CONTENT OR A FILE.
func InitialModel(content string) Text {
	ta := textarea.New()
	ta.SetValue(content)
	ta.Focus()
	ta.SetWidth(120)
	ta.ShowLineNumbers = true
	ta.CharLimit = 0
	ta.Prompt = "│ "

	return Text{
		Textarea: ta,
	}
}

// INIT INITIALIZES THE EDITOR WITH BLINKING CURSOR BEHAVIOR.
func (m Text) Init() tea.Cmd {
	return textarea.Blink
}

// UPDATE HANDLES USER INPUT AND UPDATES THE EDITOR STATE.
func (m Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			if err := validateYAML(m.Textarea.Value()); err != nil {
				m.ErrMsg = "YAML Error: " + err.Error()
			} else {
				m.Quitting = true
				return m, tea.Quit
			}
		case "esc":
			m.Quitting = true
			return m, tea.Quit
		case "ctrl+n":
			// Add a newline at the end (quick add)
			m.Textarea.SetValue(m.Textarea.Value() + "\n")
		}
	}

	var cmd tea.Cmd
	m.Textarea, cmd = m.Textarea.Update(msg)
	return m, cmd
}

// VIEW RETURNS THE CURRENT VIEW FOR RENDERING THE EDITOR IN THE TERMINAL.
func (m Text) View() string {
	if m.Quitting {
		return ""
	}

	errorMsg := ""
	if m.ErrMsg != "" {
		errorMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.ErrMsg) + "\n"
	}

	info := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(
		"Lines: %d | Chars: %d | Ctrl+S: Save • Ctrl+N: New Line • Esc: Quit",
		len(strings.Split(m.Textarea.Value(), "\n")),
		len(m.Textarea.Value()),
	))

	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("Edit YAML:"),
		m.Textarea.View(),
		errorMsg+info,
	)
}

// READYAMLFILE READS A YAML FILE FROM THE FILESYSTEM AND RETURNS ITS CONTENT.
func ReadYAMLFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// VALIDATEYAML CHECKS IF THE INPUT CONTENT IS VALID YAML.
func validateYAML(content string) error {
	var dummy interface{}
	return yaml.Unmarshal([]byte(content), &dummy)
}
