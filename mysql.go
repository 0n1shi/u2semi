package httphoneypot

type RequestMySQLModel struct {
	Method string            `json:"method"`
	URL    string            `json:"url"`
	Proto  string            `json:"version"`
	Header map[string]string `json:"headers"`
	Body   string            `json:"body"`
}

type MySQLClient struct {

	Username string 
	Password string 
	Hostname string 
	DB       string 
}

func 