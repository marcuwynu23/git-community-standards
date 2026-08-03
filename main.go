package main

import (
	"os"

	"git-community-standards/internal/app"
)

var version = "dev"

func main() {
	os.Exit(app.NewApp(version).Run(os.Args[1:]))
}
