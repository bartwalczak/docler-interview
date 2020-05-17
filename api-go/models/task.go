package models

import (
	"fmt"
	"time"
)

// Task is the daily task record
type Task struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Due         time.Time `json:"due,omitempty"`
	Priority    string    `json:"priority,omitempty"`
}

// NewTask constructs a new task with default values
func NewTask(title string) Task {
	return Task{
		Title:    title,
		Status:   "to do",
		Due:      time.Now().Add(24 * time.Hour), // set for the next day
		Priority: "medium",
	}
}

// Validate checks the validity of the record
func (t Task) Validate() []error {
	var errs []error
	if t.Title == "" {
		errs = append(errs, fmt.Errorf("Task needs a title"))
	}
	if t.Due.Before(time.Now()) {
		errs = append(errs, fmt.Errorf("Task already past due"))
	}
	return errs
}

// DueIfZero fixes the zero time binding
func (t *Task) DueIfZero(due time.Time) {
	t0 := time.Time{}
	if t.Due == t0 {
		t.Due = due
	}
}
