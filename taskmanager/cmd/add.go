package cmd

import (
	"fmt"
	"strings"

	"github.com/fuskovic/gophercises/taskmanager/pkg/store"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.New()
		if err != nil {
			fmt.Printf("failed to init store : %s\n", err)
			return
		}
		defer s.Close()

		task := store.Task{
			Description: strings.Join(args, " "),
			IsCompleted: false,
		}

		if err := s.Add(task); err != nil {
			fmt.Printf("failed to add task : %s\n", err)
			return
		}

		fmt.Println("successfully added task!")

		incompleteTasks, err := s.List()
		if err != nil {
			fmt.Printf("failed to list incomplete tasks : %s\n", err)
			return
		}

		fmt.Println("**********INCOMPLETE TASKS**********")

		for _, task := range incompleteTasks {
			fmt.Printf("%d - %s\n", task.ID, task.Description)
		}
	},
}

func init() {
	taskCmd.AddCommand(addCmd)
}
