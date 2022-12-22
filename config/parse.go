package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

const DashDir = "gh-dash"
const ConfigYmlFileName = "config.yml"
const ConfigYamlFileName = "config.yaml"
const DEFAULT_XDG_CONFIG_DIRNAME = ".config"

var validate *validator.Validate

type Defaults struct {
	Zsh  string `yaml:"zsh_history_path" validate:"string"`
	Fish string `yaml:"fish_history_path" validate:"string"`
	Bash string `yaml:"bash_history_path" validate:"string"`

	// Preview                PreviewConfig `yaml:"preview"`
	// PrsLimit               int           `yaml:"prsLimit"`
	// IssuesLimit            int           `yaml:"issuesLimit"`
	// View                   ViewType      `yaml:"view"`
	// Layout                 LayoutConfig  `yaml:"layout,omitempty"`
	// RefetchIntervalMinutes int           `yaml:"refetchIntervalMinutes,omitempty"`
}

type Config struct {
	Default Defaults `yaml:"defaults"`
}

type configError struct {
	configDir string
	parser    ConfigParser
	err       error
}

type ConfigParser struct{}

func (parser ConfigParser) getDefaultConfig() Config {
	return Config{
		Default: Defaults{
			Zsh:  "~/.zsh_history",
			Fish: "~/.fish_history",
			Bash: "~/.bash_history",
		},
	}
}

func (parser ConfigParser) getDefaultConfigYamlContents() string {
	defaultConfig := parser.getDefaultConfig()
	yaml, _ := yaml.Marshal(defaultConfig)

	return string(yaml)
}

func (parser ConfigParser) writeDefaultConfig(configFile *os.File) error {
	_, err := configFile.WriteString(parser.getDefaultConfigYamlContents())

	if err != nil {
		return err
	}

	return nil
}

func (parser ConfigParser) createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		newConfigFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer newConfigFile.Close()
		return parser.writeDefaultConfig(newConfigFile)
	}

	return nil
}
