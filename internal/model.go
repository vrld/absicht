package internal

import (
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/jordan-wright/email"
)

type KeyMap struct {
	Edit             key.Binding
	Attach           key.Binding
	Send             key.Binding
	Save             key.Binding
	Quit             key.Binding
	RemoveAttachment key.Binding
	ScrollUp         key.Binding
	ScrollDown       key.Binding
	ScrollLeft       key.Binding
	ScrollRight      key.Binding
	HalfPageUp       key.Binding
	HalfPageDown     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Edit, k.Attach, k.RemoveAttachment, k.Send, k.Save, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Edit, k.Attach, k.RemoveAttachment},
		{k.Send, k.Save, k.Quit},
	}
}

const (
	stateReadBody = iota
	stateRemoveAttachment
)

type Model struct {
	width, height int
	keys          KeyMap
	help          help.Model
	email         *email.Email
	bodyViewport  viewport.Model
	inputState    int
}

func InitialModel() Model {
	model := Model{
		keys: KeyMap{
			Edit:             key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit email")),
			Attach:           key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "attach file")),
			Send:             key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "send")),
			Save:             key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
			Quit:             key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
			RemoveAttachment: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "remove attachment")),
			ScrollUp:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
			ScrollDown:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
			ScrollLeft:       key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "scroll left")),
			ScrollRight:      key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "scroll right")),
			HalfPageUp:       key.NewBinding(key.WithKeys("u", "ctrl+u"), key.WithHelp("u", "½ page up")),
			HalfPageDown:     key.NewBinding(key.WithKeys("d", "ctrl+d"), key.WithHelp("d", "½ page down")),
		},
		email:      email.NewEmail(),
		inputState: stateReadBody,
	}

	model.bodyViewport = viewport.New(1, 1)
	model.setDimensions(1, 1)
	model.UpdateEmailDisplay()

	return model
}

func (m *Model) ReadEmail(reader io.Reader) error {
	mail, err := email.NewEmailFromReader(reader)
	if err != nil {
		return err
	}
	m.email = mail
	m.UpdateEmailDisplay()
	return nil
}

func (m *Model) UpdateEmailDisplay() {
	emailText := string(m.email.Text)
	rendered, err := renderMarkdown(m.width, emailText)

	if err == nil {
		m.bodyViewport.SetContent(rendered)
	} else {
		m.bodyViewport.SetContent(emailText)
	}

	m.setViewportHeight()
}

func (m *Model) setDimensions(width, height int) {
	m.width = width
	m.height = height
	m.help.Width = width

	m.bodyViewport.Width = width
	m.setViewportHeight()
}

func (m *Model) setViewportHeight() {
	m.bodyViewport.Height = m.height - m.renderHeightAttachments() - m.renderHeightHeaders() - m.renderHeightBottom() - 2
}
