package main

import (
	"flag"
	app "restfbz/internal/app"
)

func main() {
	port := flag.String("port", "8999", "Port on which the server will listen")
	dbfile := flag.String("dbfile", "restfbz.db", "Filename of the sqlite db")
	flag.Parse()
	cf := app.Config{Dbfile: *dbfile, Port: *port}
	server := app.New(cf)
	server.Run()
}
