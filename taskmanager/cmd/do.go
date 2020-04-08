package cmd

import (
	"fmt"
	"strconv"

	"github.com/fuskovic/gophercises/taskmanager/pkg/store"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "mark a task as completed",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := store.New()
		if err != nil {
			fmt.Printf("failed to init store : %s\n", err)
			return
		}
		defer s.Close()

		if len(args) > 1 {
			fmt.Printf("the do command only accepts one argument but received %d", len(args))
			return
		}

		ID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("failed to convert ID string input to int : %s\n", err)
			return
		}

		if err := s.Do(ID); err != nil {
			fmt.Printf("failed to mark task ID : %d as completed : %s\n", ID, err)
			return
		}
		fmt.Printf("successfully marked task ID %d as completed!\n", ID)

		completedTasks, err := s.Completed()
		if err != nil {
			fmt.Printf("failed to list completed tasks : %s\n", err)
			return
		}

		fmt.Println("**********COMPLETED TASKS**********")

		for _, task := range completedTasks {
			fmt.Printf("%d - %s\n", task.ID, task.Description)
		}
	},
}

func init() {
	taskCmd.AddCommand(doCmd)
}
