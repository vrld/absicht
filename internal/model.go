package internal

import (
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/jordan-wright/email"
)

type KeyMap struct {
	Edit   key.Binding
	Attach key.Binding
	Send   key.Binding
	Save   key.Binding
	Quit   key.Binding
	// TODO: bindings to manage attachments
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Edit, k.Attach, k.Send, k.Save, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Edit, k.Attach},
		{k.Send, k.Save, k.Quit},
	}
}

type Model struct {
	width, height int
	keys          KeyMap
	help          help.Model
	email         *email.Email
	bodyViewport  viewport.Model
}

func InitialModel() Model {
	model := Model{
		keys: KeyMap{
			Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
			Attach: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "attach")),
			Send:   key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "send")),
			Save:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
			Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		},
		email: email.NewEmail(),
	}

	model.setDimensions(1, 1)
	model.bodyViewport.KeyMap = viewport.DefaultKeyMap()
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
	rendered, err := RenderMarkdown(m.width, emailText)
	if err == nil {
		m.bodyViewport.SetContent(rendered)
	} else {
		m.bodyViewport.SetContent(emailText)
	}
}

func (m *Model) setDimensions(width, height int) {
		m.width = width
		m.height = height
		m.help.Width = width

		m.bodyViewport.Width = width
		m.bodyViewport.Height = height
}
