package app

import (
	"fmt"
	"log"
	"net/http"
)

// Application is a simple struc holding
// * router: a simple router
type Application struct {
	router *http.ServeMux
}

// New : initialize the application
func New() (app Application) {
	app = Application{http.NewServeMux()}
	return
}

// Run laumches the http server
func (app Application) Run(port string) {
	log.Printf(fmt.Sprintf("Server started and listening on %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), app.router))
}
