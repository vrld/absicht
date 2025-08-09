package internal

import (
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
	// TODO: viewport for body
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

	// mock data
	// TODO: read from email
	model.email.From = "s@nd.er"
	model.email.To = []string{"rec@iv.er"}
	model.email.Bcc = []string{"s@cr.et"}
	model.email.Subject = "check this out"

	model.email.Text = []byte("This is a test\nThere are lines\nmany lines\nso\nmany\nlines\n\nk, bye\n\n-- \nand a signature")

	model.email.AttachFile("README.md")
	model.email.AttachFile("main.go")
	model.email.AttachFile("flake.nix")

	model.setDimensions(1, 1)
	model.bodyViewport.KeyMap = viewport.DefaultKeyMap()
	model.bodyViewport.SetContent(string(model.email.Text))

	return model
}

func (m *Model) setDimensions(width, height int) {
		m.width = width
		m.height = height
		m.help.Width = width

		m.bodyViewport.Width = width
		m.bodyViewport.Height = height
}
