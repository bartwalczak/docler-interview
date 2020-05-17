package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/bartwalczak/docler-interview/server-go/models"
	"github.com/bartwalczak/docler-interview/server-go/server"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	apiPortFlag    = flag.String("p", "8690", "API Port")
	sampleDataFlag = flag.String("sample", "", "Sample data to be imported into the db, use only for testing")
)

func envDefaults() {
	viper.SetEnvPrefix("dia")

	viper.SetConfigName("cfg")  // name of config file (without extension)
	viper.AddConfigPath(".")    // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.Fatalf("fatal error config file: %s \n", err)
	}

}

func importSampleData(s *server.Server, fp string) {
	logrus.Infof("Loading sample data: %s", *sampleDataFlag)

	var tasks struct {
		Data []struct {
			Days        int    `json:"days"` // relative due date
			Title       string `json:"title,omitempty"`
			Description string `json:"description,omitempty"`
			Priority    string `json:"priority,omitempty"`
		} `yaml:"data"`
	}

	// read yaml file
	bt, err := ioutil.ReadFile(fp)
	if err != nil {
		logrus.Errorf("Sample data not imported: %s", err)
		return
	}
	// unmarshal to temporary struct
	err = yaml.Unmarshal(bt, &tasks)
	if err != nil {
		logrus.Errorf("Sample data not imported: %s", err)
		return
	}

	// import tasks to db
	logrus.Infof("Found %d tasks", len(tasks.Data))
	tn := time.Now()
	for _, task := range tasks.Data {
		t := models.Task{
			Title:       task.Title,
			Description: task.Description,
			Due:         tn.AddDate(0, 0, task.Days),
			Priority:    task.Priority,
			Status:      "to do",
		}
		err := s.DB().Create(&t).Error
		if err != nil {
			logrus.Errorf("Record not created: %s", err)
			continue
		}
	}
	return
}

func main() {
	flag.Parse()
	envDefaults()

	apiSub := viper.Sub("api")
	dbSub := viper.Sub("db")

	// configure the server
	logrus.Info("configuring the server")
	s, err := server.New(server.Config{
		APIPort: apiSub.GetString("port"),
		APIHost: apiSub.GetString("host"),

		LogLevel: log.Lvl(apiSub.GetUint("log-level")),

		DBConfig: server.DBConfig{
			Name: dbSub.GetString("name"),
			Port: dbSub.GetString("port"),
			Host: dbSub.GetString("host"),
			Pass: dbSub.GetString("pass"),
			User: dbSub.GetString("user"),
		},
	}, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// add sample data if needed
	if *sampleDataFlag != "" {
		importSampleData(s, *sampleDataFlag)
	}

	// handle ^C
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	s.Start()

	<-quit
	s.Stop(context.Background())
}
