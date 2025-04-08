package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/stuttgart-things/survey"
)

func main() {

	// READ YAML CONTENT FROM A FILE (REPLACE WITH ACTUAL PATH)
	initialContent, err := survey.ReadYAMLFile("example.yaml") // Replace "example.yaml" with your actual file
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		os.Exit(1)
	}

	fmt.Println("Initial YAML content:")
	fmt.Println(initialContent)

	// INITIALIZE AND RUN THE TERMINAL EDITOR PROGRAM.
	p := tea.NewProgram(survey.InitialModel(initialContent), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running editor:", err)
		os.Exit(1)
	}

	// IF SUCCESSFUL, PRINT THE FINAL YAML CONTENT.
	if result, ok := m.(survey.Text); ok && result.ErrMsg == "" {
		fmt.Println("\n" + lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render("Final YAML") + "\n")
		fmt.Println(result.Textarea.Value())
	}

}
