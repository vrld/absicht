package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/absicht/internal"
)

func main() {
	model := internal.InitialModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}


