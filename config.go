package u2semi

type Conf struct {
	Repo RepoConf `yaml:"repo"`
	Web  WebConf  `yaml:"web"`
}

type RepoConf struct {
	Type  RepoType      `yaml:"type"`
	MySQL MySQLRepoConf `yaml:"mysql"`
}

type RepoType string

const (
	RepoTypeNone  RepoType = "none"
	RepoTypeMySQL RepoType = "mysql"
)

type MySQLRepoConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	DB       string `yaml:"db"`
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
