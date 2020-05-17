package server

import "github.com/labstack/gommon/log"

// Config contains server config
type Config struct {
	APIHost  string
	APIPort  string
	LogLevel log.Lvl

	DBConfig
}
