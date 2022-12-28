package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

const AppDir = "goback"

const ConfigFileName = "config.yml"

var (
	defaultFile             = filepath.Join(xdg.Home, ".zsh_history")
	defaultExcludedCommands = []string{
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

type SettingsConfig struct {
	File    string   `yaml:"file_path"`
	Exclude []string `yaml:"excluded_commands"`
}

type Config struct {
	Settings SettingsConfig `yaml:"settings"`
}

type configError struct {
	configDir string
	parser    ConfigParser
	err       error
}

type ConfigParser struct{}

func (parser ConfigParser) getDefaultConfig() Config {
	return Config{
		Settings: SettingsConfig{
			File:    defaultFile,
			Exclude: defaultExcludedCommands,
		},
	}
}

func (parser ConfigParser) getDefaultConfigYamlContents() string {
	defaultConfig := parser.getDefaultConfig()
	yaml, _ := yaml.Marshal(defaultConfig)
	return string(yaml)
}

func (e configError) Error() string {
	return fmt.Sprintf(
		`Couldn't find a config.yml configuration file.
Create one under: %s
press q to exit.
Original error: %v`,
		filepath.Join(xdg.Home, ".config", AppDir, ConfigFileName),
		e.err,
	)
}

func (parser ConfigParser) writeDefaultConfigContents(newConfigFile *os.File) error {
	_, err := newConfigFile.WriteString(parser.getDefaultConfigYamlContents())
	if err != nil {
		return err
	}
	return nil
}

func (parser ConfigParser) createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		newConfigFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer newConfigFile.Close()
		return parser.writeDefaultConfigContents(newConfigFile)
	}

	return nil
}

func (parser ConfigParser) getConfigFileOrCreateIfMissing() (*string, error) {
	var err error
	configDir := filepath.Join(xdg.Home, ".config")
	if configDir == "" {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return nil, configError{parser: parser, configDir: configDir, err: err}
		}
	}
	prsConfigDir := filepath.Join(configDir, AppDir)
	err = os.MkdirAll(prsConfigDir, os.ModePerm)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}
	configFilePath := filepath.Join(prsConfigDir, ConfigFileName)
	err = parser.createConfigFileIfMissing(configFilePath)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}
	return &configFilePath, nil
}

type parsingError struct {
	err error
}

func (e parsingError) Error() string {
	return fmt.Sprintf("failed parsing config.yml: %v", e.err)
}

func (parser ConfigParser) readConfigFile(path string) (Config, error) {
	config := parser.getDefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return config, configError{parser: parser, configDir: path, err: err}
	}
	err = yaml.Unmarshal((data), &config)
	return config, err
}

func initParser() ConfigParser {
	return ConfigParser{}
}

func ParseConfig() (Config, error) {
	var config Config
	var err error
	parser := initParser()
	configFilePath, err := parser.getConfigFileOrCreateIfMissing()
	if err != nil {
		return config, parsingError{err: err}
	}
	config, err = parser.readConfigFile(*configFilePath)
	if err != nil {
		return config, parsingError{err: err}
	}
	return config, nil
}
