package survey

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// InitialSaveModel initializes the Save model with a file picker and text input for the filename.
func InitialSaveModel(content string) Save {
	fp := filepicker.New()
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.CurrentDirectory, _ = os.UserHomeDir()

	ti := textinput.New()
	ti.Placeholder = "filename.txt"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return Save{
		content:    content,
		filepicker: fp,
		filename:   ti,
		status:     "SELECT A DIRECTORY, THEN ENTER FILENAME",
	}
}

func (m Save) Init() tea.Cmd {
	return tea.Batch(
		m.filepicker.Init(),
		textinput.Blink,
	)
}

func (m Save) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			if m.selectedDir == "" {
				if path := m.filepicker.CurrentDirectory; path != "" {
					m.selectedDir = path
					return m, nil
				}
				m.status = "Please select a directory first"
				return m, nil
			}

			if m.filename.Value() == "" {
				m.status = "Please enter a filename"
				return m, nil
			}

			fullPath := filepath.Join(m.selectedDir, m.filename.Value())
			err := os.WriteFile(fullPath, []byte(m.content), 0644)
			if err != nil {
				m.err = err
				m.status = fmt.Sprintf("Error saving to: %s", fullPath)
			} else {
				m.saved = true
				m.status = fmt.Sprintf("Successfully saved to: %s", fullPath)
			}
			return m, tea.Quit

		case tea.KeyBackspace:
			if m.selectedDir != "" && m.filename.Value() == "" {
				// Go back to directory selection if filename is empty and backspace is pressed
				m.selectedDir = ""
				m.status = "Select a directory"
				return m, nil
			}
		}
	}

	if m.selectedDir == "" {
		m.filepicker, cmd = m.filepicker.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.filename, cmd = m.filename.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Save) View() string {
	if m.quitting {
		return ""
	}

	var view strings.Builder

	if m.selectedDir == "" {
		view.WriteString("Browse to select save location:\n")
		view.WriteString(m.filepicker.View())
	} else {
		fullPath := filepath.Join(m.selectedDir, m.filename.Value())
		view.WriteString("Full save path:\n")
		view.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Render(fullPath) + "\n\n")

		view.WriteString("Edit filename if needed:\n")
		view.WriteString(m.filename.View() + "\n")
	}

	statusStyle := lipgloss.NewStyle()
	switch {
	case m.saved:
		statusStyle = statusStyle.Foreground(lipgloss.Color("42"))
	case m.err != nil:
		statusStyle = statusStyle.Foreground(lipgloss.Color("196"))
	default:
		statusStyle = statusStyle.Foreground(lipgloss.Color("241"))
	}
	view.WriteString("\n" + statusStyle.Render(m.status) + "\n\n")

	help := "↑↓: Navigate • Enter: Select • Ctrl+C: Quit"
	if m.selectedDir != "" {
		help = "Enter: Confirm save • Backspace: Go Back • Ctrl+C: Cancel"
	}
	view.WriteString(lipgloss.NewStyle().Faint(true).Render(help))

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(view.String())
}
