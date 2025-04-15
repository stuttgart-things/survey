package survey

import (
	"strings"

	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"
)

func (m ListModel) Init() tea.Cmd {
	return nil
}

func InitListModel(defaultVars []string) ListModel {

	input := textinput.New()
	input.Prompt = "> "
	input.Width()
	input.CharLimit = 100

	return ListModel{
		variables: defaultVars,
		input:     input,
	}
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.FinalOutput = strings.Join(m.variables, "\n") + "\n"
			m.shouldQuit = true
			return m, tea.Quit
		case "enter":
			if m.editing {
				if m.addingNew {
					if m.input.Value() != "" {
						m.variables = append(m.variables, m.input.Value())
						m.index = len(m.variables) - 1
					}
					m.addingNew = false
				} else {
					m.variables[m.index] = m.input.Value()
				}
				m.editing = false
				m.input.Blur()
			} else {
				m.editing = true
				m.input.SetValue(m.variables[m.index])
				m.input.Focus()
			}
		case "n":
			if !m.editing {
				m.addingNew = true
				m.editing = true
				m.input.SetValue("")
				m.input.Focus()
				return m, nil
			}
		case "up":
			if m.index > 0 && !m.editing {
				m.index--
			}
		case "down":
			if m.index < len(m.variables)-1 && !m.editing {
				m.index++
			}
		case "backspace":
			if !m.editing && len(m.variables) > 0 {
				m.variables = append(m.variables[:m.index], m.variables[m.index+1:]...)
				if m.index >= len(m.variables) {
					m.index = len(m.variables) - 1
				}
			}
		}
	}

	if m.editing {
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func (m ListModel) View() string {
	if m.shouldQuit {
		return ""
	}

	var sb strings.Builder

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Edit Ansible Variables")
	sb.WriteString(title + "\n\n")

	for i, v := range m.variables {
		if i == m.index {
			if m.editing {
				v = lipgloss.NewStyle().
					Background(lipgloss.Color("#04B575")).
					Render(m.input.View())
			} else {
				v = lipgloss.NewStyle().
					Background(lipgloss.Color("#04B575")).
					Render(v)
			}
		}
		sb.WriteString(v + "\n")
	}

	if m.addingNew && m.editing {
		sb.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("New variable:"))
		sb.WriteString("\n" + m.input.View())
	}

	help := lipgloss.NewStyle().Faint(true).Render(`
↑/↓: Navigate  •  Enter: Edit/Save  •  n: Add New Entry
Backspace: Delete  •  Esc: Save & Exit
`)
	sb.WriteString("\n" + help)

	return sb.String()
}
