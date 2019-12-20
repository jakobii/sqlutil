package sqlutil

import (
	"fmt"
	_ "github.com/lib/pq" // comment
)

// Connection contains the variables required to create a sql server connection.
type Connection struct {
	Server   string
	Database string
	Username string
	Password string
	Port     int
}

// NewConnection creates a new Connection with some default values.
func NewConnection() Connection {
	var c Connection
	c.Port = 1433
	return c
}

func (c Connection) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.Username, c.Password, c.Server, c.Database)
}