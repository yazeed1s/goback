package main

import (
	"fmt"
	"io"
	// "math"
	"os"
	"regexp"
	// "fyne.io/fyne/v2/cmd/fyne/commands"
	// "strings"
)

const COMMANDS_SIZE = 10000

var (
	zshfile  *os.File
	bytes    int
	err      error
	cmd      string
	splits   [COMMANDS_SIZE]string
	commands [COMMANDS_SIZE]string
)

func OpenFile(path string) (*os.File, error) {
	zshfile, err = os.Open(path)
	if err != nil {
		panic(err)
	}
	return zshfile, err
}

func writeBuffer(file *os.File) string {

	buf := make([]byte, COMMANDS_SIZE)
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

func extractCommands(cmds string) [COMMANDS_SIZE]string {
	a := regexp.MustCompile(`[:;]+`)
	b := regexp.MustCompile(`^[0-9\s]{1,3}`)
	splits := a.Split(cmd, -1)
	for i := 0; i < len(splits); i++ {
		if !b.MatchString(splits[i]) {
			commands[i] = splits[i]
		}
	}
	return commands
}

func main() {
	zshfile, err = OpenFile("/Users/yazeed_1/.zsh_history")
	defer zshfile.Close()
	filepath := "/Users/yazeed_1/.zsh_history"
	file, err := OpenFile(filepath)
	if err != nil {
		panic(err)
	}
	history := writeBuffer(file)
	if err != nil {
		panic(err)
	}
    cmds := extractCommands(history)
	fmt.Println(cmds)
}
