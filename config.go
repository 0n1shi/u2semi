package u2semi

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Repo RepoConf `yaml:"repo"`
	Web  WebConf  `yaml:"web"`
}

type RepoConf struct {
	DSN string `yaml:"dsn"`
}

type WebConf struct {
	Port            int                    `yaml:"port"`
	Headers         []*Header              `yaml:"headers"`
	ContentDir      string                 `yaml:"content_directory"`
	DirListTemplate string                 `yaml:"directory_listing_template"`
	Contents        map[string]*WebContent `yaml:"contents"`
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type WebContent struct {
	Body string `yaml:"body"`
}

func LoadConf(path string) (*Conf, error) {
	var conf Conf
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
