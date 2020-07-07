package cmd

import (
	"fmt"
	"os"

	"gophercises/task/db"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your todos",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Failed to list todos", err)
			os.Exit(1)
		}
		
		if len(tasks) == 0 {
			fmt.Println("You have no tasks 🥳")
			return
		}

		fmt.Println("You have the 👇 tasks 🙈:")
		for i, t := range tasks {
			fmt.Println(i+1, t.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
