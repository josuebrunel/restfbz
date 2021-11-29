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

const LIMIT_LIMIT int = 1000000

func concat(ss ...string) string {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteString(s)
	}
	return sb.String()
}

// Config is a structure holding the application settings
type Config struct {
	Dbfile string
	Port   string
}

// New : initialize the application
func New(cf Config) (app Application) {
	app = Application{http.NewServeMux(), cf}
	return
}

// Application is a simple struc holding
// * router: a simple router
type Application struct {
	router *http.ServeMux
	cf     Config
}

// Run laumches the http server
func (app Application) Run() {
	app.routes()
	log.Printf(fmt.Sprintf("Server started and listening on %s", app.cf.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", app.cf.Port), app.router))
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
				// check limit params is less than 1M
				if key == "limit" && val > LIMIT_LIMIT {
					http.Error(w, "Invalid value for limit param (value > 1000000)", http.StatusBadRequest)
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
