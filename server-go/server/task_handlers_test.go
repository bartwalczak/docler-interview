package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bartwalczak/docler-interview/server-go/models"
	"github.com/gavv/httpexpect"
)

func TestTaskHandlers(t *testing.T) {
	// Setup
	e := httpexpect.New(t, serverURL)

	t.Log("get task list - should be empty")
	e.GET("/tasks").Expect().Status(http.StatusOK).
		JSON().Array().Length().Equal(0)

	t.Log("post invalid task")
	e.POST("/tasks").WithJSON(models.Task{
		Description: "No title",
	}).Expect().Status(http.StatusBadRequest)

	t.Log("post valid task")
	e.POST("/tasks").WithJSON(models.Task{
		Title: "New task",
	}).Expect().Status(http.StatusCreated).
		Header("Location").Contains("tasks/id/1")

	t.Log("get task list - should be length 1")
	e.GET("/tasks").Expect().Status(http.StatusOK).
		JSON().Array().Length().Equal(1)

	t.Log("get the first task")
	e.GET("/tasks/id/1").Expect().Status(http.StatusOK).
		JSON().Object().Path("$.title").String().Equal("New task")

	t.Log("check non-existent task")
	e.GET("/tasks/id/2").Expect().Status(http.StatusNotFound)

	t.Log("check invalid id task")
	e.GET("/tasks/id/abc").Expect().Status(http.StatusNotFound)

	t.Log("try to update valid task with yestarday's date")
	tn := time.Now()
	e.PUT("/tasks/id/1").WithJSON(models.Task{
		Due: tn.Add(-24 * time.Hour),
	}).Expect().Status(http.StatusBadRequest)

	t.Log("update valid task")
	obj := e.PUT("/tasks/id/1").WithJSON(models.Task{
		Description: "Added description",
		Status:      "done",
	}).Expect().Status(http.StatusOK).
		JSON().Object()

	obj.Path("$.description").Equal("Added description")
	obj.Path("$.status").Equal("done")

	t.Log("get the updated task")
	obj = e.GET("/tasks/id/1").Expect().Status(http.StatusOK).
		JSON().Object()

	obj.Path("$.description").Equal("Added description")
	obj.Path("$.status").Equal("done")

	t.Log("create an invalid task with put")
	e.PUT("/tasks/id/2").WithJSON(models.Task{
		ID:          2,
		Description: "No title",
	}).Expect().Status(http.StatusBadRequest)

	t.Log("create a new task with put")
	e.PUT("/tasks/id/2").WithJSON(models.Task{
		ID:    2,
		Title: "Second Task",
	}).Expect().Status(http.StatusCreated).
		Header("Location").Contains("tasks/id/2")

	t.Log("get task list - should be length 2 now")
	e.GET("/tasks").Expect().Status(http.StatusOK).
		JSON().Array().Length().Equal(2)

	t.Log("get the created task")
	obj = e.GET("/tasks/id/2").Expect().Status(http.StatusOK).
		JSON().Object()

	obj.Path("$.title").Equal("Second Task")
	obj.Path("$.status").Equal("to do")

	t.Log("get tasks for today")
	e.GET("/tasks/today").Expect().Status(http.StatusOK).
		JSON().Array().Length().Equal(0)

	t.Log("get tasks for tomorrow")
	d := time.Now().Add(24 * time.Hour)
	e.GET(fmt.Sprintf("/tasks/date/%s", d.Format("2006-01-02"))).Expect().
		Status(http.StatusOK).
		JSON().Array().Length().Equal(2)

}
