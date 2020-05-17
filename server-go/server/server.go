package server

import (
	"context"
	"fmt"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// Server is the instance of the api server
type Server struct {
	Config

	baseURL string
	indent  string

	server *echo.Echo
	db     *gorm.DB
}

// New constructs a new server
func New(cfg Config, db *gorm.DB) (*Server, error) {
	var err error

	// if database is not supplied, use config
	if db == nil {
		// open db connection
		db, err = cfg.DBConfig.Open()
		if err != nil {
			return nil, err
		}
	}

	if cfg.APIHost == "" {
		cfg.APIHost = "localhost"
	}

	// return server
	return &Server{
		Config: cfg,

		indent:  "   ",
		baseURL: fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort),

		db: db,
	}, nil
}

// DB returns the server database
func (s *Server) DB() *gorm.DB {
	return s.db
}

// URI builds the URI for a particular handler
func (s *Server) URI(handler echo.HandlerFunc, params ...interface{}) string {
	return fmt.Sprintf("http://%s", path.Join(s.baseURL, s.server.URI(handler, params...)))
}

// Start starts the server
func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true

	// setup logger
	logrus.SetLevel(logrus.Level(s.LogLevel))
	e.Logger.SetLevel(s.LogLevel)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status} latency=${latency_human} req=${bytes_in} resp=${bytes_out}\n",
	}))

	e.File("/", "assets/html/index.html")
	e.Static("/js", "assets/js")
	e.Static("/css", "assets/css")

	// add routes and handlers
	e.GET("/check", s.healthCheckHandler)

	e.GET("/tasks", s.getTasksHandler)
	e.GET("/tasks/today", s.getTasksForTodayHandler)
	e.GET("/tasks/date/:date", s.getTasksForADayHandler)
	e.GET("/tasks/id/:id", s.getTaskByIDHandler)

	e.POST("/tasks", s.postTaskHandler)
	e.PUT("/tasks/id/:id", s.putTaskHandler)

	s.server = e

	// Start server
	go func() {
		logrus.Infof("starting server at %s", s.APIPort)
		if err := s.server.Start(fmt.Sprintf(":%s", s.APIPort)); err != nil {
			e.Logger.Errorf("shutting down the server: %s", err)
		}
	}()

	return nil
}

// Stop stops the server
func (s *Server) Stop(c context.Context) {
	s.server.Shutdown(c)
	if s.db != nil {
		s.db.Close()
	}
}
