package postgre

import (
	"fmt"
	"time"
)

type DBOptions struct {
	Host                  string
	Port                  int
	Database              string
	Username              string
	Password              string
	SSLMode               string
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifetime time.Duration
}

func NewDBOptions() *DBOptions {
	return &DBOptions{
		Host:                  "localhost",
		Port:                  5432,
		Database:              "postgres",
		Username:              "postgres",
		Password:              "postgres",
		SSLMode:               "disable",
		MaxIdleConnections:    10,
		MaxOpenConnections:    100,
		ConnectionMaxLifetime: time.Hour,
	}
}

func (o *DBOptions) WithHost(host string) *DBOptions {
	o.Host = host
	return o
}
func (o *DBOptions) WithPort(port int) *DBOptions {
	o.Port = port
	return o
}
func (o *DBOptions) WithDatabase(db string) *DBOptions {
	o.Database = db
	return o
}
func (o *DBOptions) WithUsername(user string) *DBOptions {
	o.Username = user
	return o
}
func (o *DBOptions) WithPassword(pass string) *DBOptions {
	o.Password = pass
	return o
}
func (o *DBOptions) WithSSLMode(mode string) *DBOptions {
	o.SSLMode = mode
	return o
}
func (o *DBOptions) WithMaxIdleConnections(n int) *DBOptions {
	o.MaxIdleConnections = n
	return o
}
func (o *DBOptions) WithMaxOpenConnections(n int) *DBOptions {
	o.MaxOpenConnections = n
	return o
}
func (o *DBOptions) WithConnectionMaxLifetime(d time.Duration) *DBOptions {
	o.ConnectionMaxLifetime = d
	return o
}

func (o *DBOptions) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", o.Host, o.Port, o.Username, o.Password, o.Database, o.SSLMode)
}
