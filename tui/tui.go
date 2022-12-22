package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 0)

type item struct {
	c string
}

func (i item) Title() string { 
    return i.c 
}

func (i item) Description() string {
	return ""
}

func (i item) FilterValue() string { 
    return i.c 
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func InitModel(cmds []string, file *os.File) model {
	items := make([]list.Item, len(cmds))
	for i, cmd := range cmds {
		items[i] = item{c: cmd}
	}

	d := list.NewDefaultDelegate()
	d.SetSpacing(1)
	d.ShowDescription = false
	m := model{list: list.New(items, d, 0, 0)}
    stat, err := file.Stat()
    if err != nil {
        panic(err)
    }
	m.list.Title = stat.Name()
	m.list.SetStatusBarItemName("command", "commands")
	m.list.Styles.PaginationStyle.Margin(0)
	return m
}
