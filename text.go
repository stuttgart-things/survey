package survey

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"
	"gopkg.in/yaml.v3"
)

func InitialModel(content string) Text {
	ta := textarea.New()
	ta.SetValue(content)
	ta.Focus()

	// Set reasonable defaults for a scrollable editor
	ta.ShowLineNumbers = true
	ta.CharLimit = 0
	ta.Prompt = "│ "
	ta.SetHeight(100) // Fixed visible lines, content will scroll

	// Optimized styling for better visibility
	ta.Styles.Blurred.Base = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	ta.Styles.Blurred.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("236"))

	ta.Styles.Blurred.LineNumber = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	return Text{
		Textarea: ta,
	}
}

func (m Text) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		tea.EnterAltScreen,
	)
}

func (m Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.WindowSize = msg
		// Leave some space for header/footer
		height := msg.Height - 3
		if height < 3 {
			height = 3 // Minimum reasonable height
		}
		m.Textarea.SetWidth(msg.Width - 4)
		m.Textarea.SetHeight(20)

		return m, nil

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
			if m.ErrMsg != "" {
				m.ErrMsg = ""
			} else {
				m.Quitting = true
				return m, tea.Quit
			}
		case "ctrl+n":
			m.Textarea.SetValue(m.Textarea.Value() + "\n")
		}
	}

	var cmd tea.Cmd
	m.Textarea, cmd = m.Textarea.Update(msg)
	return m, cmd
}

func (m Text) View() string {
	if m.Quitting {
		return ""
	}

	// Header with line count info
	totalLines := len(strings.Split(m.Textarea.Value(), "\n"))
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true).
		Render(fmt.Sprintf("YAML Editor (Lines: %d)", totalLines))

	// Error message
	errorMsg := ""
	if m.ErrMsg != "" {
		errorMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render("Error: "+m.ErrMsg) + "\n"
	}

	// Footer with key bindings
	footer := lipgloss.NewStyle().
		Faint(true).
		Render("↑↓: Scroll • ^S: Save • ^N: New Line • Esc: Quit")

	// Combine all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.Textarea.View(),
		errorMsg,
		footer,
	)
}

func ReadYAMLFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

func validateYAML(content string) error {
	var dummy interface{}
	return yaml.Unmarshal([]byte(content), &dummy)
}
