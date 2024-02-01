package contract

type ServerConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	AbsPath  string `json:"abs-path"`
}
