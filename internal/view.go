package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
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
		attachmentHeight = 1 // TODO: figure out why this needs to be 1 instead of 0
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

func RenderMarkdown(width int, markdown string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
		glamour.WithStylesFromJSONBytes([]byte(markdownStyleJson)),
	)

	if err != nil {
		return "", err
	}

	result, err := renderer.Render(markdown)
	if err != nil {
		return "", err
	}

	return result, nil
}

var markdownStyleJson = `{
  "document": {
    "block_prefix": "\n",
    "block_suffix": "\n",
    "color": "15",
    "margin": 0
  },
  "block_quote": {
    "indent": 1,
    "indent_token": "│ "
  },
  "paragraph": {},
  "list": {
    "level_indent": 2
  },
  "heading": {
    "block_suffix": "\n",
    "color": "4",
    "bold": true
  },
  "h1": {
    "prefix": " ",
    "suffix": " ",
    "color": "0",
    "background_color": "12",
    "bold": true
  },
  "h2": {
    "prefix": "## "
  },
  "h3": {
    "prefix": "### "
  },
  "h4": {
    "prefix": "#### "
  },
  "h5": {
    "prefix": "##### "
  },
  "h6": {
    "prefix": "###### ",
    "color": "14",
    "bold": false
  },
  "text": {},
  "strikethrough": {
    "crossed_out": true
  },
  "emph": {
    "italic": true
  },
  "strong": {
    "bold": true
  },
  "hr": {
    "color": "8",
    "format": "\n--------\n"
  },
  "item": {
    "block_prefix": "• "
  },
  "enumeration": {
    "block_prefix": ". "
  },
  "task": {
    "ticked": "[✓] ",
    "unticked": "[ ] "
  },
  "link": {
    "color": "6",
    "underline": true
  },
  "link_text": {
    "color": "2",
    "bold": true
  },
  "image": {
    "color": "13",
    "underline": true
  },
  "image_text": {
    "color": "7",
    "format": "Image: {{.text}} →"
  },
  "code": {
    "prefix": " ",
    "suffix": " ",
    "color": "0",
    "background_color": "7"
  },
  "code_block": {
    "color": "8",
    "margin": 2,
    "chroma": {
      "text": {
        "color": "#C4C4C4"
      },
      "error": {
        "color": "#F1F1F1",
        "background_color": "#F05B5B"
      },
      "comment": {
        "color": "#676767"
      },
      "comment_preproc": {
        "color": "#FF875F"
      },
      "keyword": {
        "color": "#00AAFF"
      },
      "keyword_reserved": {
        "color": "#FF5FD2"
      },
      "keyword_namespace": {
        "color": "#FF5F87"
      },
      "keyword_type": {
        "color": "#6E6ED8"
      },
      "operator": {
        "color": "#EF8080"
      },
      "punctuation": {
        "color": "#E8E8A8"
      },
      "name": {
        "color": "#C4C4C4"
      },
      "name_builtin": {
        "color": "#FF8EC7"
      },
      "name_tag": {
        "color": "#B083EA"
      },
      "name_attribute": {
        "color": "#7A7AE6"
      },
      "name_class": {
        "color": "#F1F1F1",
        "underline": true,
        "bold": true
      },
      "name_constant": {},
      "name_decorator": {
        "color": "#FFFF87"
      },
      "name_exception": {},
      "name_function": {
        "color": "#00D787"
      },
      "name_other": {},
      "literal": {},
      "literal_number": {
        "color": "#6EEFC0"
      },
      "literal_date": {},
      "literal_string": {
        "color": "#C69669"
      },
      "literal_string_escape": {
        "color": "#AFFFD7"
      },
      "generic_deleted": {
        "color": "#FD5B5B"
      },
      "generic_emph": {
        "italic": true
      },
      "generic_inserted": {
        "color": "#00D787"
      },
      "generic_strong": {
        "bold": true
      },
      "generic_subheading": {
        "color": "#777777"
      },
      "background": {
        "background_color": "#373737"
      }
    }
  },
  "table": {},
  "definition_list": {},
  "definition_term": {},
  "definition_description": {
    "block_prefix": "\n🠶 "
  },
  "html_block": {},
  "html_span": {}
}`
