package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/bartwalczak/docler-interview/server-go/models"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

func getURL(endpoints ...string) string {
	return fmt.Sprintf("http://%s/%s", apiURI, path.Join(endpoints...))
}

func printHelp() {
	help := `Use the following commands:
	t, today - to list tasks for today
	l, list - to list all tasks
	s, show <id:int> - to show details of a given task
	d, done <id:int> - to mark task as done (not implemented yet)
	P, push <id:int> <days:int> - to push a given task a number of days into the future (not implemented yet)
	p, priority <id:int> <priority:string> - to set task priority (not implemented yet) \n`
	fmt.Print(help)
}

func listTasks(more ...string) {
	cl := resty.New()
	path := append([]string{"tasks"}, more...)
	resp, err := cl.R().Get(getURL(path...))
	if err != nil {
		fmt.Printf("Error getting task list: %s\n", err)
		return
	}
	if resp.StatusCode() != http.StatusOK {
		fmt.Printf("Request returned status code: %d\n", resp.StatusCode())
		return
	}
	bd := resp.Body()
	var tasks []models.Task
	err = json.Unmarshal(bd, &tasks)
	if err != nil {
		fmt.Printf("Error parsing task list response: %s\n", err)
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Title", "Priority", "Status", "Due"})
	for _, v := range tasks {
		t.Append([]string{fmt.Sprintf("%d", v.ID), v.Title, v.Priority, v.Status, v.Due.Format("2006-01-02")})
	}
	t.Render()
}

func showTask(strID string) {
	cl := resty.New()
	resp, err := cl.R().Get(getURL(path.Join("tasks", "id", strID)))
	if err != nil {
		fmt.Printf("Error getting task: %s\n", err)
		return
	}
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNotFound {
			fmt.Printf("Task ID %s not found\n", strID)
			return
		}
		fmt.Printf("Request returned status code: %d\n", resp.StatusCode())
		return
	}
	bd := resp.Body()
	var task models.Task
	err = json.Unmarshal(bd, &task)
	if err != nil {
		fmt.Printf("Error parsing task response: %s\n", err)
		return
	}
	fmt.Printf(`Task %d:
Title: %s
Description: %s
Status: %s
Priority: %s
Due: %s
`, task.ID, task.Title, task.Description, task.Status, task.Priority, task.Due.Format("2006-01-02"))
	return
}

func pushTask(strID, strDays string) {
	// TODO
}

func doneTask(strID string) {
	// TODO
}

func taskPriority(strID, strPriority string) {
	// TODO
}

func routeCommands() {
	switch flag.Arg(0) {
	default:
		printHelp()
	case "t":
		fallthrough
	case "today":
		listTasks("today")
	case "l":
		fallthrough
	case "list":
		listTasks()
	case "s":
		fallthrough
	case "show":
		if flag.NArg() < 2 {
			printHelp()
			return
		}
		showTask(flag.Arg(1))
	case "d":
		fallthrough
	case "done":
		if flag.NArg() < 2 {
			printHelp()
			return
		}
		doneTask(flag.Arg(1))
	case "P":
		fallthrough
	case "push":
		if flag.NArg() < 2 {
			printHelp()
			return
		}
		p := "1"
		if flag.NArg() == 3 {
			p = flag.Arg(2)
		}
		pushTask(flag.Arg(1), p)
	case "p":
		fallthrough
	case "priority":
		if flag.NArg() < 3 {
			printHelp()
			return
		}
		taskPriority(flag.Arg(1), flag.Arg(2))
	}
}

var (
	apiURI string
)

func main() {
	flag.Usage = printHelp
	flag.Parse()

	viper.SetConfigName("cfg")  // name of config file (without extension)
	viper.AddConfigPath(".")    // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Sprintf("fatal error config file: %s \n", err))
	}
	apisub := viper.Sub("api")
	if apisub == nil {
		panic("config file is invalid")
	}

	apiURI = fmt.Sprintf("%s:%s", apisub.GetString("host"), apisub.GetString("port"))

	routeCommands()
}
