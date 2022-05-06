package httphoneypot

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
	Web   WebConfig   `yaml:"web"`
}

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	DB       string `yaml:"db"`
}

type WebConfig struct {
	Port    int       `yaml:"port"`
	Headers []*Header `yaml:"headers"`
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
