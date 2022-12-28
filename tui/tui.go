package tui

import (
	"goback/config"
	"log"
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Margin(0, 0)

type KeyMap struct {
	toggleTitleBar key.Binding
	Copy           key.Binding
	// Execute        key.Binding
}

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Copy: key.NewBinding(
			key.WithKeys("c", "copy"),
			key.WithHelp("c", "copy command"),
		),
		// Execute: key.NewBinding(
		// 	key.WithKeys("r", "run"),
		// 	key.WithHelp("r", "execute command"),
		// ),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle title"),
		),
	}
}

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
	list   list.Model
	keys   *KeyMap
	choice string
	config config.Config
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
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil
		case key.Matches(msg, m.keys.Copy):
			m.choice = m.list.SelectedItem().FilterValue()
			m.list.NewStatusMessage(m.choice + " is copied!")
			err := clipboard.WriteAll(m.choice)
			if err != nil {
				m.list.NewStatusMessage("Command " + m.choice + " cannot be copied!!")

			}
			return m, nil
			// this is a bad idea, coping the command alone would be a better solution
			// case key.Matches(msg, m.keys.Execute):
			// 	m.choice = m.list.SelectedItem().FilterValue()
			// TODO: execute m.choice
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmds []tea.Cmd

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func InitModel(cmds []string, file *os.File) *model {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	keys := DefaultKeyMap()
	items := make([]list.Item, len(cmds))
	for i, cmd := range cmds {
		items[i] = item{c: cmd}
	}
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	l := list.New(items, d, 0, 0)
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	l.Title = stat.Name()
	l.SetStatusBarItemName("command", "commands")
	l.Styles.PaginationStyle.Margin(0)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Copy,
			// keys.Execute,
			keys.toggleTitleBar,
		}
	}
	return &model{
		list:   l,
		keys:   keys,
		choice: "",
		config: cfg,
	}

}
