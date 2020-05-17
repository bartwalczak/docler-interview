package models

import (
	"testing"
	"time"

	"github.com/go-test/deep"
)

func TestTaskNew(t *testing.T) {
	task := NewTask("new title")
	if task.Due.Before(time.Now()) {
		t.Error("Task should not be due before now.")
	}
	exp := Task{
		Title:    "new title",
		Priority: "medium",
		Status:   "to do",
		Due:      task.Due,
	}
	if diff := deep.Equal(task, exp); diff != nil {
		t.Error(diff)
	}
}

func TestTaskValidate(t *testing.T) {
	ts := []struct {
		txt    string
		errlen int
		src    Task
	}{
		{
			txt:    "empty task",
			errlen: 2,
		},
		{
			txt:    "no title",
			errlen: 1,
			src: Task{
				Due: time.Now().Add(10 * time.Minute),
			},
		},
		{
			txt:    "past due",
			errlen: 1,
			src: Task{
				Title: "Some task",
				Due:   time.Now().Add(-10 * time.Minute),
			},
		},
		{
			txt:    "valid task",
			errlen: 0,
			src: Task{
				Title: "Some task",
				Due:   time.Now().Add(10 * time.Minute),
			},
		},
	}
	for _, tc := range ts {
		t.Log(tc.txt)
		act := tc.src.Validate()
		if len(act) != tc.errlen {
			t.Errorf("Error length different. Expected: %d / Actual: %d", tc.errlen, len(act))
		}
	}
}

func TestDueIfZero(t *testing.T) {
	tn := time.Now()
	ts := []struct {
		txt string
		exp time.Time
		due time.Time
		src time.Time
	}{
		{
			txt: "empty",
		},
		{
			txt: "not changed",
			src: tn,
			exp: tn,
			due: time.Now().Add(10 * time.Minute),
		},
		{
			txt: "changed",
			exp: tn,
			due: tn,
		},
	}
	for _, tc := range ts {
		t.Log(tc.txt)
		task := Task{
			Due: tc.src,
		}
		task.DueIfZero(tc.due)
		act := task.Due
		if act != tc.exp {
			t.Errorf("Error: values differ. Expected: %s / Actual: %s", tc.exp, act)
		}
	}

}
