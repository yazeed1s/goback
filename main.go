package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const BUF_SIZE = 260000

var (
	zshfile *os.File
	bytes   int
	err     error
	cmd     string
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
			cmd = string(buf[:bytes])
		}
	}
	return cmd
}

func extractCommands(cmd string) []string {
	// a line in .zsh_history = : 1671232234:0;git commit -m"first commit"
	// reg matches the ": 1671232234:0;" of the line,
	reg := regexp.MustCompile(`(?m)^[:;0-9\s]{0,15}`)
	// leaves 'git commit -m"first commit"' alone
	splits := reg.ReplaceAllLiteralString(cmd, "")
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

func main() {
	filepath := "/Users/yazeed_1/.zsh_history"
	zshfile, err := OpenFile(filepath)
	if err != nil {
		panic(err)
	}
	defer zshfile.Close()
	history := writeBuffer(zshfile)
	cmds := extractCommands(history)
	filtered := filterCommands(cmds)
	for _, cmd := range filtered {
		fmt.Println("Command = ", cmd)
	}
}

// TODO: handle fish_histiry && bash_history
// TODO: do the tui
// TODO: finish Configurations options
// TODO: add a command line option
// TODO: improve reading performance
