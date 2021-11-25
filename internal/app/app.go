package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	fbz "restfbz/pkg/fizzbuzz"
	"strconv"
	"strings"
)

func concat(ss ...string) string {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteString(s)
	}
	return sb.String()
}

func withLoggingMiddleware(h http.HandlerFunc) http.HandlerFunc {

	type LoggingResponseWriter struct {
		http.ResponseWriter
		statusCode int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		lw := LoggingResponseWriter{w, http.StatusOK}
		log.Printf("%s %s", r.Method, r.URL.String())
		h(lw, r)
		log.Printf("%d", lw.statusCode)
	}
}

// New : initialize the application
func New() (app Application) {
	app = Application{http.NewServeMux()}
	return
}

// Application is a simple struc holding
// * router: a simple router
type Application struct {
	router *http.ServeMux
}

// Run laumches the http server
func (app Application) Run(port string) {
	app.routes()
	log.Printf(fmt.Sprintf("Server started and listening on %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), app.router))
}

func (app Application) routes() {
	app.router.HandleFunc("/", withLoggingMiddleware(app.handleFizzBuzz()))
}

func (app Application) handleFizzBuzz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only GET allowed
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		paramsRequired := 5
		intParams := map[string]int{"int1": -1, "int2": -1, "limit": -1}
		strParams := map[string]string{"str1": "", "str2": ""}

		for key, values := range r.URL.Query() {
			if _, ok := intParams[key]; ok {
				val, err := strconv.Atoi(values[0])
				if err != nil {
					http.Error(w, concat("Invalid value for key: ", key), http.StatusBadRequest)
					return
				}
				intParams[key] = val
			} else {
				strParams[key] = values[0]
			}
			paramsRequired--
		}
		if paramsRequired > 0 {
			http.Error(w, "Invalid number of params", http.StatusBadRequest)
			return
		}

		res := fbz.FizzBuzz(intParams["int1"], intParams["int2"], intParams["limit"],
			strParams["str1"], strParams["str2"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
