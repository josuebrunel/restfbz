package main

import (
	"os"
	app "restfbz/internal/app"
)

func main() {
	var port string
	if len(os.Args) < 2 {
		port = "8080"
	} else {
		port = os.Args[1]
	}

	server := app.New()
	server.Run(port)
}
