package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vrld/absicht/internal"
)

var emailPath string

func init() {
	flag.StringVar(&emailPath, "r", "-", "Read initial email from this path; `-' means stdin.")
}

func main() {
	flag.Parse()
	model := internal.InitialModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if emailPath == "-" {
		stat, _ := os.Stdin.Stat()
		hasPipedInput := (stat.Mode() & os.ModeCharDevice) == 0
		if hasPipedInput {
			model.ReadEmail(bufio.NewReader(os.Stdin))
		}
	} else if file, err := os.Open(emailPath); err == nil {
		model.ReadEmail(file)
	} else {
		go p.Send(err)
	}

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}


