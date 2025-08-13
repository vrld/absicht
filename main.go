package main

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/absicht/internal"
)

func main() {
	model := internal.InitialModel()

	stat, _ := os.Stdin.Stat()
	hasPipedInput := (stat.Mode() & os.ModeCharDevice) == 0
	if hasPipedInput {
		model.ReadEmail(bufio.NewReader(os.Stdin))
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}


