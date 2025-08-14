package internal

import (
	"fmt"
	"log"
	"net/textproto"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jordan-wright/email"
	"github.com/rymdport/portal/filechooser"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		// TODO: show this on the bottom for some time, see ansicht
		panic(msg)

	case tea.WindowSizeMsg:
		m.setDimensions(msg.Width, msg.Height)

	case tea.KeyMsg:
		if cmd := m.handleKeyMessage(msg); cmd != nil {
			return m, cmd
		}

	case EditEmail:
		return m, m.editEmail()

	case UpdateEmail:
		m.email.From = msg.From
		m.email.To = msg.To
		m.email.Cc = msg.Cc
		m.email.Bcc = msg.Bcc
		m.email.Subject = msg.Subject
		m.email.Headers = msg.Headers
		m.email.Text = msg.Text
		m.UpdateEmailDisplay()
		return m, nil
	}

	m.bodyViewport.Update(msg)

	return m, nil
}

func (m *Model) handleKeyMessage(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Edit):
		return func() tea.Msg { return EditEmail{} }

	case key.Matches(msg, m.keys.Attach):
		m.attachFile()
		// TODO event based on outcome of selection
		return func() tea.Msg { return "redraw" }

	case key.Matches(msg, m.keys.Send):
		// TODO: call msmtp and quit
		return nil

	case key.Matches(msg, m.keys.Save):
		m.saveToFile()
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

func (m *Model) editEmail() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	tempFile, err := os.CreateTemp("", "absicht-email-*.eml")
	if err != nil {
		return func() tea.Msg { return fmt.Errorf("failed to create temp file: %w", err) }
	}

	err = writeEmailToFile(tempFile, m.email)
	tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return func() tea.Msg { return fmt.Errorf("cannot write temp file: %w", err) }
	}

	cmd := exec.Command(editor, tempFile.Name())
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		defer os.Remove(tempFile.Name())
		if err != nil {
			return err
		}

		editedFile, err := os.Open(tempFile.Name())
		if err != nil {
			return fmt.Errorf("cannot open edited file: %w", err)
		}

		editedEmail, err := email.NewEmailFromReader(editedFile)
		if err != nil {
			return fmt.Errorf("cannot parses edited email: %w", err)
		}

		return UpdateEmail{
			From:    editedEmail.From,
			To:      editedEmail.To,
			Cc:      editedEmail.Cc,
			Bcc:     editedEmail.Bcc,
			Subject: editedEmail.Subject,
			Headers: editedEmail.Headers,
			Text:    editedEmail.Text,
		}

	})
}

func (m *Model) saveToFile() {
	options := filechooser.SaveFileOptions{CurrentName: "mail.eml"}
	files, err := filechooser.SaveFile("absicht", "Save to file", &options)
	if err != nil {
		log.Fatalln(err)
	}

	bytes, err := m.email.Bytes()
	if err != nil {
		log.Fatalln(err)
	}

	for _, filename := range files {
    filename = strings.TrimPrefix(filename, "file://")
		if err := os.WriteFile(filename, bytes, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
}

func writeEmailToFile(file *os.File, email *email.Email) (err error) {
	canonicalHeaders := []struct{ key, value string }{
		{"From", email.From},
		{"To", strings.Join(email.To, ", ")},
		{"Cc", strings.Join(email.Cc, ", ")},
		{"Bcc", strings.Join(email.Bcc, ", ")},
		{"Subject", email.Subject},
	}
	for _, h := range canonicalHeaders {
		_, err = fmt.Fprintf(file, "%s: %s\n", h.key, h.value)
		if err != nil {
			return err
		}
	}

	for key, values := range email.Headers {
		switch key {
		case "Content-Type", "Content-Transfer-Encoding", "Mime-Version", "Message-Id":
			// skip header

		default:
			for _, v := range values {
				_, err = fmt.Fprintf(file, "%s: %s\n", key, v)
				if err != nil {
					return err
				}
			}
		}
	}

	text := strings.ReplaceAll(string(email.Text), "\r\n", "\n")
	_, err = fmt.Fprintf(file, "\n%s\n", strings.Trim(text, "\n"))
	return err
}


type EditEmail struct {}

type UpdateEmail struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Headers textproto.MIMEHeader
	Text    []byte
}
