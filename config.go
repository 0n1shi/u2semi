package u2semi

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
