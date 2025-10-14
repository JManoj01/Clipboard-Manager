package ui

import (
	"clipboard_manager/storage"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(lipgloss.Color("#EE6FF8")).
		Bold(true)
)

type item struct {
	entry storage.ClipboardEntry
}

func (i item) Title() string {
	if i.entry.IsImage {
		return fmt.Sprintf("🖼️  [Image] - ID: %d", i.entry.ID)
	}

	preview := i.entry.Text
	if len(preview) > 60 {
		preview = preview[:60] + "..."
	}
	preview = strings.ReplaceAll(preview, "\n", " ")

	return fmt.Sprintf("[%d] %s", i.entry.ID, preview)
}

func (i item) Description() string {
	return fmt.Sprintf("Category: %s | %s", i.entry.Category, FormatTimeAgo(i.entry.Timestamp))
}

func (i item) FilterValue() string { return i.entry.Text }

type StatusMsg string

type model struct {
	list     list.Model
	viewport viewport.Model
	db       *storage.Database
	viewing  bool
	selected *storage.ClipboardEntry
	status   string
}

func NewBubbleTeaUI(db *storage.Database) *model {
	entries, _ := db.GetRecent(50)

	items := []list.Item{}
	for _, entry := range entries {
		items = append(items, item{entry: entry})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📋 Clipboard Manager"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	vp := viewport.New(80, 20)

	return &model{
		list:     l,
		viewport: vp,
		db:       db,
		viewing:  false,
		status:   "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusMsg:
		m.status = string(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if !m.viewing {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.selected = &i.entry
					m.viewing = true
					m.viewport.SetContent(m.formatEntryView(i.entry))
				}
			}

		case "esc":
			if m.viewing {
				m.viewing = false
				m.selected = nil
			}

		case "d":
			if !m.viewing {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.db.DeleteEntry(i.entry.ID)
					m.refreshList()
					m.status = fmt.Sprintf("Deleted entry #%d", i.entry.ID)
				}
			}

		case "r":
			m.refreshList()
			m.status = "Refreshed"
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 6
	}

	var cmd tea.Cmd
	if m.viewing {
		m.viewport, cmd = m.viewport.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.viewing && m.selected != nil {
		return m.viewport.View() + "\n\n" +
			lipgloss.NewStyle().Faint(true).Render("Press ESC to go back | q to quit")
	}

	footer := lipgloss.NewStyle().Faint(true).Render(m.status + "  |  Enter: View  d: Delete  r: Refresh  q: Quit")
	return m.list.View() + "\n" + footer
}

func (m *model) refreshList() {
	entries, _ := m.db.GetRecent(50)
	items := []list.Item{}
	for _, entry := range entries {
		items = append(items, item{entry: entry})
	}
	m.list.SetItems(items)
}

func (m *model) formatEntryView(entry storage.ClipboardEntry) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf(" Entry #%d ", entry.ID)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Category: %s\n", entry.Category))
	b.WriteString(fmt.Sprintf("Time: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05")))

	if entry.IsImage {
		b.WriteString(fmt.Sprintf("Type: Image\n"))
		b.WriteString(fmt.Sprintf("Path: %s\n", entry.ImagePath))
	} else {
		b.WriteString(fmt.Sprintf("Length: %d characters\n", len(entry.Text)))
		if entry.Language != "" {
			b.WriteString(fmt.Sprintf("Language: %s\n", entry.Language))
		}
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("━", 80))
	b.WriteString("\n\n")

	if !entry.IsImage {
		b.WriteString(entry.Text)
	}

	return b.String()
}

func NewProgram(db *storage.Database) *tea.Program {
	return tea.NewProgram(NewBubbleTeaUI(db), tea.WithAltScreen())
}

func RunBubbleTea(db *storage.Database) error {
	p := NewProgram(db)
	_, err := p.Run()
	return err
}
