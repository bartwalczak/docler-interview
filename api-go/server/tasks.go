package server

import (
	"github.com/bartwalczak/docler-interview/api-go/models"
	"github.com/labstack/echo/v4"
)

func (s *Server) getTaskByID(id int) (models.Task, error) {
	var task models.Task
	r := s.DB().Where("id = ?", id).Limit(1).Find(&task)
	if r.Error != nil {
		if r.RecordNotFound() {
			return task, echo.ErrNotFound
		}
		return task, r.Error
	}
	return task, nil
}

func (s *Server) getTasks(where map[string]interface{}) ([]models.Task, error) {
	var tasks []models.Task
	db := s.DB()
	// add where conditions
	for cond, val := range where {
		if val == nil {
			db = db.Where(cond)
		} else {
			if vals, ok := val.([]interface{}); ok {
				db = db.Where(cond, vals...)
			} else {
				db = db.Where(cond, val)
			}
		}
	}
	// find tasks
	err := db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	// output
	return tasks, nil
}
