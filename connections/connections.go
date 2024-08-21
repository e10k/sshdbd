package connections

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

type Connections map[string]Connection

func (c Connections) GetConnection(connId string) (*Connection, error) {
	conn, ok := c[connId]
	if !ok {
		return nil, fmt.Errorf("invalid connection id: '%v'", connId)
	}

	return &conn, nil
}
