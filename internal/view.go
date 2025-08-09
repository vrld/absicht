package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	cardHeaders := m.renderHeaders()
	cardAttachments := m.renderAttachments()

	bodyHeight := m.height - (2 + lipgloss.Height(cardHeaders) + lipgloss.Height(cardAttachments))
	cardBody := m.renderBody(bodyHeight)

	return fmt.Sprintf("%s\n%s\n%s\n%s", cardHeaders, cardBody, cardAttachments, m.help.View(m.keys))
}

var borderColor = lipgloss.Color("4")

func (m Model) renderHeaders() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width - 2)

	lines := []string{
		fmt.Sprint("From: ", m.email.From),
		fmt.Sprint("To: ", strings.Join(m.email.To, ", ")),
	}

	if m.email.Cc != nil {
		lines = append(lines, fmt.Sprint("Cc: ", strings.Join(m.email.Cc, ", ")))
	}

	if m.email.Bcc != nil {
		lines = append(lines, fmt.Sprint("Bcc: ", strings.Join(m.email.Bcc, ", ")))
	}

	lines = append(lines, fmt.Sprint("Subject: ", m.email.Subject))

	return style.Render(strings.Join(lines, "\n"))
}

func (m Model) renderAttachments() string {
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

func (m Model) renderBody(height int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width - 2)

	m.bodyViewport.Height = height

	return style.Height(height).Render(m.bodyViewport.View())
}
