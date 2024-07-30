package config

type Connection struct {
	Host     string
	Port     int
	Dbname   string
	Username string
	Password string
}

type Config struct {
	Connections map[string]Connection
}
