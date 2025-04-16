package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/stuttgart-things/survey"
)

var (
	output      string
	defaultVars = []string{
		"manage_filesystem+-true",
		"update_packages+-true",
		"install_requirements+-true",
		"install_motd+-true",
		"username+-sthings",
		"lvm_home_sizing+-'15%'",
		"lvm_root_sizing+-'35%'",
		"lvm_var_sizing+-'50%'",
	}
)

func main() {
	p := tea.NewProgram(survey.InitListModel(defaultVars))
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	if model, ok := m.(survey.ListModel); ok && model.FinalOutput != "" {
		output = model.FinalOutput
	}

	fmt.Println("\nFinal output:", output)
}
