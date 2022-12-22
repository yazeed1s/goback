package app

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	t "goback/tui"

	tea "github.com/charmbracelet/bubbletea"
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

func openFile(path string) (*os.File, error) {
	zshfile, err = os.Open(path)
	if err != nil {
		panic(err)
	}
	return zshfile, err
}

func writeBuffer(file *os.File) string {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	bufSize := stat.Size()
	buf := make([]byte, bufSize+10)
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

func InitTool() {
	filepath := "/Users/yazeed_1/.zsh_history"
	zshfile, err := openFile(filepath)
	if err != nil {
		panic(err)
	}
	defer zshfile.Close()
	buffer := writeBuffer(zshfile)
	cmds := extractCommands(buffer)
	filtered := filterCommands(cmds)
	p := tea.NewProgram(t.InitModel(filtered, zshfile), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
