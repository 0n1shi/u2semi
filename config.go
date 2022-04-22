package httphoneypot

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
}

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	DB       string `yaml:"db"`
}
