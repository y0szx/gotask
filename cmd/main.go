package main

import (
	"fmt"
	"os"

	"phone_book/internal/cli"
	"phone_book/internal/module"
	"phone_book/internal/storage"
)

const filename = "orders.json"

func main() {
	storageJSON := storage.NewStorage(filename)
	phoneBookService := module.NewModule(module.Deps{Storage: storageJSON})
	commands := cli.NewCli(cli.Deps{Module: phoneBookService})

	if err := commands.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
