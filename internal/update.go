package internal

import (
	"bytes"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/google/shlex"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jordan-wright/email"
	"github.com/rymdport/portal/filechooser"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
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
		if m.inRemoveAttachment {
			return m, m.handleRemoveKeyMessage(msg)
		} else if cmd := m.handleKeyMessage(msg); cmd != nil {
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

	case RemoveAttachment:
		m.email.Attachments = slices.Delete(m.email.Attachments, msg.index, msg.index+1)
		m.inRemoveAttachment = false
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
		return nil

	case key.Matches(msg, m.keys.Send):
		sendmail, err := shlex.Split(viper.GetString("sendmail"))
		if err != nil {
			return func() tea.Msg { return err }
		}

		if len(sendmail) == 0 {
			return func() tea.Msg { return fmt.Errorf("No send command given") }
		}

		cmd := exec.Command(sendmail[0], sendmail[1:]...)
		mailBytes, err := m.email.Bytes()
		if err != nil {
			return func() tea.Msg { return err }
		}
		cmd.Stdin = bytes.NewReader(mailBytes)
		return tea.ExecProcess(cmd, func(err error) tea.Msg { return err })

	case key.Matches(msg, m.keys.Save):
		m.saveToFile()
		return nil

	case key.Matches(msg, m.keys.Quit):
		return tea.Quit

	case key.Matches(msg, m.keys.RemoveAttachment) && len(m.email.Attachments) > 0:
		m.inRemoveAttachment = true
		return nil
	}

	return nil
}

func (m *Model) handleRemoveKeyMessage(msg tea.KeyMsg) tea.Cmd {
	if msg.Type == tea.KeyEscape {
		m.inRemoveAttachment = false
		return nil
	}

	attachmentIndex := runeToAttachmentIndex(msg)
	if attachmentIndex >= 0 && attachmentIndex < len(m.email.Attachments) {
		return RemoveAttachmentCmd(attachmentIndex)
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

	m.email.HTML, err = markdownToHtml(m.email.Text)
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

func markdownToHtml(markdown []byte) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf bytes.Buffer
	if err := md.Convert(markdown, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type EditEmail struct{}

type UpdateEmail struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Headers textproto.MIMEHeader
	Text    []byte
}

type RemoveAttachment struct{ index int }

func RemoveAttachmentCmd(index int) tea.Cmd {
	return func() tea.Msg {
		return RemoveAttachment{index}
	}
}

func runeToAttachmentIndex(msg tea.KeyMsg) int {
	rune := msg.Runes[0]
	// 1 -> 0, 2 -> 1, ..., 0 -> 9
	if rune >= '0' && rune <= '9' {
		return (int(rune-'0') + 9) % 10
	}

	// a -> 10, b -> 11, ...
	if rune >= 'a' && rune <= 'z' {
		return int(rune-'a') + 10
	}

	return -1
}
