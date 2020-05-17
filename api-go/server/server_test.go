package server

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bartwalczak/docler-interview/api-go/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

const (
	serverURL = "http://localhost:8788"
)

func TestMain(m *testing.M) {
	// setup temp dir
	tmpDir, err := ioutil.TempDir(".", ".temp")
	if err != nil {
		logrus.Fatalf("Can't create the temp folder: %s", err)
	}
	// clean up after
	defer func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			logrus.Fatalf("Can't clean up after testing: %s", err)
		}
	}()

	// setup mock store
	db, err := gorm.Open("sqlite3", filepath.Join(tmpDir, "gorm.db"))
	if err != nil {
		logrus.Fatalf("Can't start the db store: %s", err)
	}
	db.AutoMigrate(models.Models()...)

	// create the server
	s, err := New(Config{
		APIPort:  "8788",
		LogLevel: 6,
	}, db)
	if err != nil {
		logrus.Fatalf("Can't start the server: %s", err)
	}

	logrus.Info("starting server")
	s.Start()
	// run tests
	logrus.Info("running tests")
	m.Run()

	// stop the server
	logrus.Info("stopping server")
	s.Stop(context.Background())
}
