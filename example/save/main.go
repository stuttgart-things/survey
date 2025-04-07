package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/stuttgart-things/survey"
)

func main() {

	content := `# Example Content
	This will be saved to the selected location.
	Modify this string with your actual content.`

	p := tea.NewProgram(survey.InitialSaveModel(content))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

}
