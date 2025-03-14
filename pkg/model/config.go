package model

import "time"

// Config type holds all configuration for the application
type Config struct {
	Server   Server
	Database Database
	App      App
}

// Server type holds the server-specific configuration
type Server struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Database type holds the database-specific configuration
type Database struct {
	Host               string
	Port               string
	User               string
	Password           string
	DBName             string
	SSLMode            string
	MaxOpenConnections int
	MaxIdleConnections int
}

// App type holds application-specific configuration
type App struct {
	Port        string
	Environment string
}
