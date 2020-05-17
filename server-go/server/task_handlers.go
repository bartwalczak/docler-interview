package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bartwalczak/docler-interview/server-go/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) getTasksHandler(c echo.Context) error {
	tasks, err := s.getTasks(nil)
	if err != nil {
		logrus.Errorf("error getting tasks: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, tasks, s.indent)
}

func (s *Server) getTasksForTodayHandler(c echo.Context) error {
	tn := time.Now()
	td := time.Date(tn.Year(), tn.Month(), tn.Day(), 0, 0, 0, 0, tn.Location())
	tasks, err := s.getTasks(map[string]interface{}{"due BETWEEN ? AND ?": []interface{}{td, td.AddDate(0, 0, 1)}})
	// Postgres specific query:
	// tasks, err := s.getTasks(map[string]interface{}{"date_trunc('day', due) = current_date": nil})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, tasks, s.indent)
}

func (s *Server) getTasksForADayHandler(c echo.Context) error {
	d := c.Param("date")
	td, err := time.Parse("2006-01-02", d)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	tasks, err := s.getTasks(map[string]interface{}{"due BETWEEN ? AND ?": []interface{}{td, td.AddDate(0, 0, 1)}})
	// Postgres specific query:
	// tasks, err := s.getTasks(map[string]interface{}{"date_trunc('day', due) = ?": d})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, tasks, s.indent)
}

func (s *Server) getTaskByIDHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NotFoundHandler(c)
	}

	task, err := s.getTaskByID(id)
	if err != nil {
		if err == echo.ErrNotFound {
			return echo.NotFoundHandler(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSONPretty(http.StatusOK, task, s.indent)
}

func (s *Server) postTaskHandler(c echo.Context) error {
	task := models.NewTask("")
	due := task.Due

	err := c.Bind(&task)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// time zero fix
	task.DueIfZero(due)

	errs := task.Validate()
	if len(errs) > 0 {
		logrus.Debug("Validation errors")
		logrus.Debug(errs)
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}

	err = s.DB().Create(&task).Error
	if err != nil {
		logrus.Errorf("Error creating task: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c.Response().Header().Set("Location", s.URI(s.getTaskByIDHandler, task.ID))
	return c.NoContent(http.StatusCreated)
}

func (s *Server) putTaskHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	uid := uint(id)

	task, err := s.getTaskByID(id)
	if err != nil {
		if err == echo.ErrNotFound {
			// bind body content
			newTask := models.NewTask("")
			due := newTask.Due
			err := c.Bind(&newTask)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
			// if new id is not given, return 404
			if newTask.ID == 0 {
				return echo.NotFoundHandler(c)
			}

			newTask.DueIfZero(due)

			errs := newTask.Validate()
			if len(errs) > 0 {
				logrus.Debug("Validation errors")
				logrus.Debug(errs)
				return echo.NewHTTPError(http.StatusBadRequest, errs)
			}
			// see if the ids match
			if newTask.ID != uid {
				newTask.ID = uid
			}

			// create if a new id is provided
			err = s.DB().Create(&newTask).Error
			if err != nil {
				logrus.Errorf("Error creating task: %s", err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			c.Response().Header().Set("Location", s.URI(s.getTaskByIDHandler, newTask.ID))
			return c.NoContent(http.StatusCreated)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	due := task.Due
	// bind on the old task to update the values
	err = c.Bind(&task)
	if err != nil {
		logrus.Errorf("error binding new task to old: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// the id must not change for any reason
	task.ID = uid

	// zero date fix
	task.DueIfZero(due)

	errs := task.Validate()
	if len(errs) > 0 {
		logrus.Debug("Validation errors")
		logrus.Debug(errs)
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}

	err = s.DB().Save(&task).Error
	if err != nil {
		logrus.Errorf("Error saving task: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSONPretty(http.StatusOK, task, s.indent)
}
