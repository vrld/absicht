package internal

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rymdport/portal/filechooser"
	"log"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setDimensions(msg.Width, msg.Height)

	case tea.KeyMsg:
		if cmd := m.handleKeyMessage(msg); cmd != nil {
			return m, cmd
		}
	}

	m.bodyViewport.Update(msg)

	return m, nil
}

func (m *Model) handleKeyMessage(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Edit):
		// TODO: yield to editor
		return nil

	case key.Matches(msg, m.keys.Attach):
		m.attachFile()
		return nil

	case key.Matches(msg, m.keys.Send):
		// TODO: call msmtp and quit
		return nil

	case key.Matches(msg, m.keys.Send):
		// TODO: prompt for name and save
		return nil

	case key.Matches(msg, m.keys.Quit):
		return tea.Quit
	}

	return nil
}

func (m *Model) attachFile() {
	options := filechooser.OpenFileOptions{Multiple: true}
	files, err := filechooser.OpenFile("absicht", "Select Attachment(s)", &options)
	if err != nil {
		log.Fatalln(err)
	}

	for _, filename := range files {
		m.email.AttachFile(strings.TrimPrefix(filename, "file://"))
	}
}
