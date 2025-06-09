package main

import (
	_ "go.uber.org/automaxprocs"
	"miniblog/cmd/mb-apiserver/app"
	"os"
)

func main() {
	command := app.NewMiniBlogCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

}
