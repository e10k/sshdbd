package config

import (
	"fmt"
)

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

func (c Config) GetConnection(connId string) (*Connection, error) {
	conn, ok := c.Connections[connId]
	if !ok {
		return nil, fmt.Errorf("invalid connection id: '%v'", connId)
	}

	return &conn, nil
}
