package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BRO3886/google-tasks-cli/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

// tasksCmd represents the tasks command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "View, create, list and telete tasks in a tasklist",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tasks called")
	},
}

var viewTasksCmd = &cobra.Command{
	Use:   "view",
	Short: "View tasks in a tasklist",
	Long: `
	Use this command to view tasks in a selected 
	tasklist for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ReadCredentials()
		client := getClient(config)

		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		list, err := getTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Println("Choose a Tasklist:")
		for index, i := range list {
			fmt.Printf("[%d] %s\n", index+1, i.Title)
		}
		fmt.Printf("Choose an option: ")
		var option int
		if _, err := fmt.Scan(&option); err != nil {
			log.Fatalf("Unable to read option: %v", err)
		}
		fmt.Println("Tasks in '" + list[option-1].Title + "':\n")

		tasks, err := getTasks(srv, list[option-1].Id)
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		for index, i := range tasks {
			color.Green("[%d] %s\n", index+1, i.Title)
			due, err := time.Parse(time.RFC3339, i.Due)
			if err != nil {
				fmt.Printf("No Due Date\n")
			} else {
				color.Yellow("Due %s\n", due.Format("Mon Jan 2 2006 3:04PM"))
			}
		}

	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create task in a tasklist",
	Long: `
	Use this command to create tasks in a selected 
	tasklist for the currently signed in account
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ReadCredentials()
		client := getClient(config)

		srv, err := tasks.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve tasks Client %v", err)
		}

		list, err := getTaskLists(srv)
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Println("Choose a Tasklist:")
		for index, i := range list {
			fmt.Printf("[%d] %s\n", index+1, i.Title)
		}
		fmt.Printf("Choose an option: ")
		var option int
		if _, err := fmt.Scan(&option); err != nil {
			log.Fatalf("Unable to read option: %v", err)
		}
	},
}

func init() {
	tasksCmd.AddCommand(viewTasksCmd, createTaskCmd)
	rootCmd.AddCommand(tasksCmd)
}

func getTasks(srv *tasks.Service, id string) ([]*tasks.Task, error) {
	r, err := srv.Tasks.List(id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve tasks. %v", err)
	}
	if len(r.Items) == 0 {
		return nil, errors.New("No Tasklist found")
	}
	return r.Items, nil
}
