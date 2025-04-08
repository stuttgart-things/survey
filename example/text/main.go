package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	lipgloss "github.com/charmbracelet/lipgloss/v2"

	"github.com/stuttgart-things/survey"
)

func main() {
	// READ YAML CONTENT FROM A FILE
	initialContent, err := survey.ReadYAMLFile("example.yaml")
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		os.Exit(1)
	}

	// Initialize with full terminal size handling
	p := tea.NewProgram(
		survey.InitialModel(initialContent),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse scrolling
	)

	// Run the program
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running editor:", err)
		os.Exit(1)
	}

	// Handle the final model
	if result, ok := m.(survey.Text); ok {
		if result.ErrMsg != "" {
			fmt.Println("Error:", result.ErrMsg)
			os.Exit(1)
		}

		fmt.Println("\n" + lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render("Final YAML") + "\n")
		fmt.Println(result.Textarea.Value())
	}
}
