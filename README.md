# goback

## About The Tool

A terminal based command history browser, can be handy & usefull sometimes!

###  ⚡️ Built With

- [Go](https://golang.org/)

- [bubbletea](https://github.com/charmbracelet/bubbletea)

- [bubbles](https://github.com/charmbracelet/bubbles)

- [lipgloss](https://github.com/charmbracelet/lipgloss)

- [Cobra](https://github.com/spf13/cobra)

##  ✨ Features

- Works with both .zsh_history & .bash_history

- Fast starting time

- Search for a given command

- Copy selected command to the clipboard

- Ability to exclude a given set of commands from the result (like cd , ls, rm , cp, vim, etc)

## 📌 Note:

The tool can only work on Unix-based systems (macOS, linux), windows is not supported!  

## ⭐️ Screenshots
<p align="left">
  <img src="https://github.com/Yazeed1s/goback/blob/main/screenshots/s1.png" width="750">
</p>
<p align="left">
  <img src="https://github.com/Yazeed1s/goback/blob/main/screenshots/s2.png" width="750">
</p>
<p align="left">
  <img src="https://github.com/Yazeed1s/goback/blob/main/screenshots/s3.png" width="750">
</p>

## 📦 Installation

### Git

```
// clone the repo
git clone https://github.com/Yazeed1s/goback.git
// cd into it
cd goback
// build the project
go build 
// add the binary to /go/bin/
go install

```



## 🚀 Usage

- Just run `goback` and the list should be displayed!

## 👩🏻‍🦯 Navigation

| Key            | Description |
| ---------------| ----------- |
| `↓ or j`       | Scroll down |
| `↑ or k`       | Scroll up                            |
| `→ or l`       | Move to the next page |
| `← or h `      | Move to the previous page |
| `g`            | Jump to the start |
| `G`            | Jump to the end |
| `c`            | Copy command to the clipboard |
| `/`            | Toggle filtering (sreaching for a command) |
| `esc`          | Clear filter (clear serach results) |
| `enter`        | Apply filter (or copy command) |
| `?`            | Show help |
| `?`            | Same key used to close help window |
| `t`            | Toggle title |
| `q or ctr + c` | Exit |

## ⚙️ Configuration

A config file will be generated when you first run `goback`. The file can be found in the following locations:

* macOS: ~/.config/goback/config.yml

* Linux: ~/.config/goback/config.yml

`.zsh_history`  is the default file. If you use bash shell, you can simply change the value to `.bash_history`

The config file will include the following default values:

```yml

settings:
  file_path: /Path/to/.zsh_history
  excluded_commands: 
  - ls
  - cd
  - clr
  - cd ..
  - clear
  - mkdir
  - rmdir
  - rm
  - mv
  - goback
  - cat
  - clear
  - pwd
  - vim
  - vim .
  - vi
  - vi .
  - nvim 
  - nvim .
  - code .
  - codium .
  - touch

```


## ✅  TODO: 
- [x] Add support for bash shell 
- [ ] Add support for fish shell
- [ ] Add support for windows (if possible)

## 🔥 Contributing

If you would like to add a feature or to fix a bug please feel free to send a PR.
