package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gophercises/task/cmd"
	"gophercises/task/db"

	"github.com/mitchellh/go-homedir"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, ".tasks.db")
	err := db.Init(dbPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cmd.RootCmd.Execute()
}
