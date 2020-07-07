package cmd

import (
	"fmt"
	"strconv"
	"os"

	"gophercises/task/db"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Makred a todo as done",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("You have no tasks to delete 🥶")
			os.Exit(1)
		}

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse argument", arg)
				os.Exit(1)
			}

			err = db.DeleteTask(tasks[id - 1].Key)
			if err != nil {
				fmt.Println("Failed to delete task", id, err)
				os.Exit(1)
			}
			fmt.Println("Completed task 🌸", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
