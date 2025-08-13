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
var goEditNoTimeToWaste bool

func init() {
	flag.StringVar(&emailPath, "r", "-", "Read initial email from this path; `-' means stdin.")
	flag.BoolVar(&goEditNoTimeToWaste, "e", false, "Start editing the email right after staring.")
}

func main() {
	flag.Parse()
	model := internal.InitialModel()

	err := readEmail(&model)

	p := tea.NewProgram(model, tea.WithAltScreen())

  if err != nil {
		go p.Send(err)
	}

	if goEditNoTimeToWaste {
		go p.Send(internal.EditEmail{})
	}

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running the program: %v", err)
	}
}

func readEmail(model *internal.Model) error {
	if emailPath == "-" {
		stat, _ := os.Stdin.Stat()
		hasPipedInput := (stat.Mode() & os.ModeCharDevice) == 0
		if hasPipedInput {
			return model.ReadEmail(bufio.NewReader(os.Stdin))
		}
		return nil
	}

	file, err := os.Open(emailPath)
	if err == nil {
		err = model.ReadEmail(file)
		file.Close()
	}

	return err
}
