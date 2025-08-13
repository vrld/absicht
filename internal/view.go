package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	cardHeaders := m.renderHeaders()
	var cardAttachments string
	if len(m.email.Attachments) > 0 {
		cardAttachments = m.renderAttachments() + "\n"
	}

	attachmentHeight := lipgloss.Height(cardAttachments)
	if len(m.email.Attachments) == 0 {
		attachmentHeight = 1  // TODO: figure out why this needs to be 1 instead of 0
	}
	bodyHeight := m.height - (1 + lipgloss.Height(cardHeaders) + attachmentHeight)
	cardBody := m.renderBody(bodyHeight)
	// NOTE: cardAttachments already includes the \n if there are attchments
	return fmt.Sprint(cardHeaders, "\n", cardBody, "\n", cardAttachments, m.help.View(m.keys))
}

var borderColor = lipgloss.Color("4")

func (m *Model) renderHeaders() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width - 2)

	lines := []string{
		fmt.Sprint("From: ", m.email.From),
		fmt.Sprint("To: ", strings.Join(m.email.To, ", ")),
	}

	if cc := strings.Join(m.email.Cc, ", "); cc != "" {
		lines = append(lines, fmt.Sprint("Cc: ", cc))
	}

	if bcc := strings.Join(m.email.Bcc, ", "); bcc != "" {
		lines = append(lines, fmt.Sprint("Bcc: ", bcc))
	}

	lines = append(lines, fmt.Sprint("Subject: ", m.email.Subject))

	for key, values := range m.email.Headers {
		switch key {
		case "Content-Type", "Content-Transfer-Encoding", "Mime-Version", "Message-Id":
			// skip header

		default:
			for _, v := range values {
				lines = append(lines, fmt.Sprintf("%s: %s", key, v))
			}
		}
	}

	return style.Render(strings.Join(lines, "\n"))
}

func (m *Model) renderAttachments() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width - 2)

	lines := []string{}
	for _, a := range m.email.Attachments {
		lines = append(lines, fmt.Sprint("- ", a.Filename, " ", a.ContentType))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func (m *Model) renderBody(height int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width - 2)

	m.bodyViewport.Height = height

	return style.Height(height).Render(m.bodyViewport.View())
}
