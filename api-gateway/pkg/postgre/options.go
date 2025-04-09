package postgre

import (
	"fmt"
	"time"
)

// DBOptions represents database connection options
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

// NewDBOptions creates a new DBOptions with default values
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

// WithHost sets the database host
func (o *DBOptions) WithHost(host string) *DBOptions {
	o.Host = host
	return o
}

// WithPort sets the database port
func (o *DBOptions) WithPort(port int) *DBOptions {
	o.Port = port
	return o
}

// WithDatabase sets the database name
func (o *DBOptions) WithDatabase(database string) *DBOptions {
	o.Database = database
	return o
}

// WithUsername sets the database username
func (o *DBOptions) WithUsername(username string) *DBOptions {
	o.Username = username
	return o
}

// WithPassword sets the database password
func (o *DBOptions) WithPassword(password string) *DBOptions {
	o.Password = password
	return o
}

// WithSSLMode sets the SSL mode
func (o *DBOptions) WithSSLMode(sslMode string) *DBOptions {
	o.SSLMode = sslMode
	return o
}

// WithMaxIdleConnections sets the maximum number of idle connections
func (o *DBOptions) WithMaxIdleConnections(maxIdleConnections int) *DBOptions {
	o.MaxIdleConnections = maxIdleConnections
	return o
}

// WithMaxOpenConnections sets the maximum number of open connections
func (o *DBOptions) WithMaxOpenConnections(maxOpenConnections int) *DBOptions {
	o.MaxOpenConnections = maxOpenConnections
	return o
}

// WithConnectionMaxLifetime sets the maximum lifetime of a connection
func (o *DBOptions) WithConnectionMaxLifetime(connectionMaxLifetime time.Duration) *DBOptions {
	o.ConnectionMaxLifetime = connectionMaxLifetime
	return o
}

// ConnectionString generates a PostgreSQL connection string from the options
func (o *DBOptions) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		o.Host, o.Port, o.Username, o.Password, o.Database, o.SSLMode)
}
