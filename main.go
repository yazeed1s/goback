package main

import (
	"fmt"
	// "unsafe"
	//"github.com/charmbracelet/bubbles/list"
	"io"
	"os"
	"regexp"
	"strings"

	table "github.com/calyptia/go-bubble-table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const BUF_SIZE = 300000 // this is the size of the .zsh_history (300,000 bytes)
// (big concern cuz it needs to be updated as the file grows)
// one possible solution would be to directly use Name() & Size() from fs.FileInfo
// this will sync and handle any changes in the size (grow, shrink)
// (will be applied to .bash_history and fish_history as well)
// this is the size of the .zsh_history

var (
	zshfile *os.File
	bytes   int
	err     error
	buffer  string
	non     = []string{
		"ls",
		"cd",
		"cd ..",
		"clear",
		"mkdir",
		"rmdir",
		"rm",
		"mv",
		"cat",
		"clear",
		"pwd",
		"vim",
		"vim .",
		"vi",
		"vi .",
		"nvim",
		"nvim .",
		"code .",
		"codium .",
		"touch"}
)

func OpenFile(path string) (*os.File, error) {
	zshfile, err = os.Open(path)
	if err != nil {
		panic(err)
	}
	return zshfile, err
}

func writeBuffer(file *os.File) string {
	buf := make([]byte, BUF_SIZE)
	for {
		bytes, err = file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if bytes > 0 {
			buffer = string(buf[:bytes])
		}
	}
	return buffer
}

func extractCommands(buf string) []string {
	// a line in .zsh_history = : 1671232234:0;git commit -m"first commit"
	// reg matches the ": 1671232234:0;" of the line,
	reg := regexp.MustCompile(`(?m)^[:;0-9\s]{0,15}`)
	// leaves 'git commit -m"first commit"' alone
	splits := reg.ReplaceAllLiteralString(buf, "")
	commands := strings.Split(splits, "\n")
	return commands
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func filterCommands(str []string) []string {
	var r []string
	for _, str := range str {
		if str != "" && !contains(non, str) {
			r = append(r, str)
		}
	}
	return r
}

var (
	styleDoc = lipgloss.NewStyle().Padding(1)
)

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel(cmds []string) model {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom
	n := len(cmds)
	c := 0
	tbl := table.New([]string{"Nth  ", "Command"}, w, h)
	rows := make([]table.Row, n)
	for _, i := range cmds {
		rows[c] = table.SimpleRow{
			c,
			i,
		}
		c++
	}
	tbl.SetRows(rows)
	return model{table: tbl}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := styleDoc.GetPadding()
		m.table.SetSize(
			msg.Width-left-right,
			msg.Height-top-bottom,
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(
		m.table.View(),
	)
}
func main() {
	filepath := "/Users/yazeed_1/.zsh_history"
	zshfile, err := OpenFile(filepath)
	if err != nil {
		panic(err)
	}
	defer zshfile.Close()
	buffer := writeBuffer(zshfile)
	// //fmt.Println(len(buffer))
	cmds := extractCommands(buffer)
	// //fmt.Println(len(cmds))
	filtered := filterCommands(cmds)
	//fmt.Println(len(filtered))
	// c := 0
	// for _, cmd := range filtered {
	// 	fmt.Printf("n = %d cmd = %v\n", c, cmd)
	// 	c++
	// }
	//fmt.Println(unsafe.Sizeof(BUF_SIZE))
	p := tea.NewProgram(initialModel(filtered), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// // TODO: handle fish_histiry && bash_history
// // TODO: do the tui
// // TODO: finish Configurations options
// // TODO: add a command line option
// // TODO: improve reading performance

// ----------------- style 2 --------------------------------

// const listHeight = 14

// var (
// 	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
// 	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
// 	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
// 	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
// 	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
// 	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
//  docStyle          = lipgloss.NewStyle().Margin(1, 2)
// )

// type item string

// func (i item) FilterValue() string { return "" }

// type itemDelegate struct{}

// func (d itemDelegate) Height() int                               { return 1 }
// func (d itemDelegate) Spacing() int                              { return 0 }
// func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
// func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
// 	i, ok := listItem.(item)
// 	if !ok {
// 		return
// 	}

// 	str := fmt.Sprintf("%d. %s", index+1, i)

// 	fn := itemStyle.Render
// 	if index == m.Index() {
// 		fn = func(s string) string {
// 			return selectedItemStyle.Render("> " + s)
// 		}
// 	}

// 	fmt.Fprint(w, fn(str))
// }

// type model struct {
// 	list     list.Model
// 	choice   string
// 	quitting bool
// }

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.list.SetWidth(msg.Width)
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch keypress := msg.String(); keypress {
// 		case "ctrl+c":
// 			m.quitting = true
// 			return m, tea.Quit

// 		case "enter":
// 			i, ok := m.list.SelectedItem().(item)
// 			if ok {
// 				m.choice = string(i)
// 			}
// 			return m, tea.Quit
// 		}
// 	}
// 	var cmd tea.Cmd
// 	m.list, cmd = m.list.Update(msg)
// 	return m, cmd
// }

// func (m model) View() string {
// 	if m.choice != "" {
// 		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
// 	}
// 	if m.quitting {
// 		return quitTextStyle.Render("Not hungry? That’s cool.")
// 	}
// 	return "\n" + m.list.View()
// }
// func main() {
// 	filepath := "/Users/yazeed_1/.zsh_history"
// 	zshfile, err := OpenFile(filepath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer zshfile.Close()
// 	history := writeBuffer(zshfile)
// 	cmds := extractCommands(history)
// 	filtered := filterCommands(cmds)
// 	items := []list.Item{}
// 	for _, cmd := range filtered {
// 		items = append(items, item(cmd))
// 	}

// 	const defaultWidth = 20

// 	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
// 	l.Title = "What do you want for dinner?"
// 	l.SetShowStatusBar(false)
// 	l.SetFilteringEnabled(false)
// 	l.Styles.Title = titleStyle
// 	l.Styles.PaginationStyle = paginationStyle
// 	l.Styles.HelpStyle = helpStyle

// 	m := model{list: l}

// 	p := tea.NewProgram(m)
// 	if err := p.Start(); err != nil {
// 		panic(err)
// 	}
// }
