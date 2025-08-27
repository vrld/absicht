package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const borderColorNormal = lipgloss.Color("8")
const borderColorSelected = lipgloss.Color("4")

var styleHeaderKey = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
var styleHeaderValue = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

var styleButton = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0")).Padding(0, 1).Bold(true)
var styleCancelButton = styleButton.Background(lipgloss.Color("9"))

func (m *Model) borderColor(state int) lipgloss.Color {
	if m.inputState == state {
		return borderColorSelected
	}
	return borderColorNormal
}

func (m Model) View() string {
	cardHeaders := m.renderHeaders()
	var cardAttachments string
	if len(m.email.Attachments) > 0 {
		cardAttachments = m.renderAttachments() + "\n"
	}

	bottom := m.renderBottom()

	cardBody := m.renderBody()

	// NOTE: cardAttachments already includes the \n if there are attachments
	return fmt.Sprint(cardHeaders, "\n", cardBody, "\n", cardAttachments, bottom)
}

func (m *Model) renderHeightHeaders() int {
	return len(getAllHeaders(m.email)) + 2
}

func (m *Model) renderHeaders() string {
	headers := getAllHeaders(m.email)
	maxHeaderLength := 0
	for _, h := range headers {
		maxHeaderLength = max(maxHeaderLength, len(h.key))
	}

	lines := []string{}
	for _, header := range headers {
		lines = append(lines, fmt.Sprint(
			styleHeaderKey.Width(maxHeaderLength+2).Render(header.key+":"),
			styleHeaderValue.Render(header.value),
		))
	}

	color := borderColorNormal

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(color).
		Padding(0, 1).
		Width(m.width - 2)

	return fmt.Sprint(
		renderTopLineWithTitle("header", m.width, color),
		"\n",
		style.Render(strings.Join(lines, "\n")),
	)
}

func (m *Model) renderHeightAttachments() int {
	count := len(m.email.Attachments)
	if count == 0 {
		return 0
	}
	return count + 2
}

func (m *Model) renderAttachments() string {
	lines := []string{}
	for i, a := range m.email.Attachments {
		prefix := "- "
		if m.inputState == stateRemoveAttachment {
			rune := attachmentIndexToRune(i)
			prefix = fmt.Sprintf("[%c] ", rune)
		}
		lines = append(lines, fmt.Sprint(prefix, a.Filename, " ", a.ContentType))
	}

	color := m.borderColor(stateRemoveAttachment)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(color).
		Padding(0, 1).
		Width(m.width - 2)

	return fmt.Sprint(
		renderTopLineWithTitle("attachments", m.width, color),
		"\n",
		style.Render(strings.Join(lines, "\n")),
	)
}

func (m *Model) renderHeightBottom() int {
	return 1
}

func (m *Model) renderBottom() string {
	if m.inputState == stateRemoveAttachment {
		return m.renderButtons(
			[]string{styleButton.Render("select attachment to delete")},
			[]string{styleCancelButton.Render("<esc> to cancel")},
		)
	}
	return m.renderButtons(
		[]string{
			styleButton.Render("[e]dit"),
			styleButton.Render("[a]ttach"),
			styleButton.Render("[r]emove attachment"),
			styleButton.Render("[s]ave"),
		},
		[]string{
			styleButton.Render("send (y)"),
			styleCancelButton.Render("[q]uit"),
		},
	)
}

func (m *Model) renderButtons(left, right []string) string {
	leftJoined := strings.Join(left, " ")
	rightJoined := strings.Join(right, " ")

	spacer := ""
	spacerWidth := m.width - (lipgloss.Width(leftJoined) + lipgloss.Width(rightJoined)) - 2
	if spacerWidth > 0 {
		spacer = strings.Repeat(" ", spacerWidth)
	}

	return fmt.Sprint(" ", leftJoined, spacer, rightJoined)
}

func (m *Model) renderBody() string {
	color := m.borderColor(stateReadBody)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(lipgloss.Color(color)).
		Padding(0, 1).
		Width(m.width - 2)

	return fmt.Sprint(
		renderTopLineWithTitle("body", m.width, lipgloss.Color(color)),
		"\n",
		style.Height(m.bodyViewport.Height-1).Render(m.bodyViewport.View()),
	)
}

func renderTopLineWithTitle(title string, width int, color lipgloss.Color) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	switch width {
	case 0:
		return ""
	case 1:
		return borderStyle.Render("╌")
	case 2:
		return borderStyle.Render("├┤")
	}

	titleWidth := lipgloss.Width(title) + 6
	if width < titleWidth {
		return borderStyle.Render(fmt.Sprint("┌", strings.Repeat("─", width-2), "┐"))
	}

	return borderStyle.Render(
		fmt.Sprint(
			"┌",
			strings.Repeat("─", width-titleWidth),
			"─🮤",
			title,
			"🮥─┐",
		),
	)
}

func renderMarkdown(width int, markdown string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
		glamour.WithStylesFromJSONBytes([]byte(markdownStyleJson)),
	)

	if err != nil {
		return "", err
	}

	return renderer.Render(markdown)
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

func attachmentIndexToRune(index int) rune {
	if index >= 0 && index <= 9 {
		return rune('0' + ((index + 1) % 10))
	}

	if index >= 10 && index <= 35 {
		return rune('a' + index - 10)
	}

	return 'ẞ'
}
