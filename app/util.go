package app

import (
	"fmt"
	"goback/config"
	t "goback/tui"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// const BUF_SIZE = 300000 // this is the size of the .zsh_history (300,000)
// (big concern cuz it needs to be updated as the file grows)
// one possible solution would be to directly use Name() & Size() from fs.FileInfo
// this will sync and handle any changes in the size (grow, shrink)
// (will be applied to .bash_history and fish_history as well)

var (
	zshfile  *os.File
	bytes    int
	err      error
	buffer   string
	cfg      = getCfg()
	excluded = cfg.Settings.Exclude
)

func getCfg() config.Config {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

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

func extractCommandsZsh(buf string) []string {
	// single line in .zsh_history = ': 1671232234:0;git commit -m"first commit"'
	// reg matches the ": 1671232234:0;" part of the line,
	reg := regexp.MustCompile(`(?m)^[:;0-9\s]{0,15}`)
	// extracts 'git commit -m"first commit"' alone
	splits := reg.ReplaceAllLiteralString(buf, "")
	commands := strings.Split(splits, "\n")
	return commands
}

func extractCommandsBash(buf string) []string {
	// a line in .bash_history = 'chsh -s /bin/zsh'
	commands := strings.Split(buf, "\n")
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
		if str != "" && !contains(excluded, str) {
			r = append(r, str)
		}
	}
	return r
}

func reverseCmds[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func InitTool() {
	filepath := cfg.Settings.File
	file, err := openFile(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buffer := writeBuffer(file)
	var (
		cmds     []string
		filtered []string
	)
	strings.Contains(filepath, ".zsh_history")
	if strings.Contains(filepath, ".zsh_history") {
		cmds = extractCommandsZsh(buffer)
	} else if strings.Contains(filepath, ".bash_history") {
		cmds = extractCommandsBash(buffer)
	} // TODO: fish history
	filtered = filterCommands(cmds)
	reverseCmds(filtered)
	p := tea.NewProgram(t.InitModel(filtered, file), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
