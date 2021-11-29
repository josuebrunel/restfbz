package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	fbz "restfbz/pkg/fizzbuzz"
	"restfbz/pkg/stats"
	"strconv"
	"strings"
)

// LIMIT_LIMIT is the defined limit of limit query params
const LIMIT_LIMIT int = 1000000

func concat(ss ...string) string {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteString(s)
	}
	return sb.String()
}

type responseData struct {
	Count int         `json:"count"`
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

// Config is a structure holding the application settings
type Config struct {
	Dbfile string
	Port   string
}

// New : initialize the application
func New(cf Config) (app Application) {
	app = Application{http.NewServeMux(), cf, nil}
	return
}

// Application is a simple struc holding
// * router: a simple router
type Application struct {
	router *http.ServeMux
	cf     Config
	db     *sql.DB
}

// initdb initialize the database
func (app *Application) initdb() error {
	var file *os.File
	if _, err := os.Stat(app.cf.Dbfile); os.IsNotExist(err) {
		file, err = os.Create(app.cf.Dbfile)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	db, err := sql.Open("sqlite3", app.cf.Dbfile)
	if err != nil {
		return err
	}
	app.db = db
	return err
}

// makeMigrations run all the migration of all the dependency packages
func (app *Application) makeMigrations() error {
	sr := stats.New(app.db)
	return sr.CreateTables()
}

// Run laumches the http server
func (app Application) Run() {
	err := app.initdb()
	if err != nil {
		log.Fatal("Failed to init database")
	}
	defer app.db.Close()
	err = app.makeMigrations()
	if err != nil {
		log.Fatal("Failed to run migrations")
	}
	app.routes()
	log.Print(fmt.Sprintf("Server started and listening on %s", app.cf.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", app.cf.Port), app.router))
}

func (app Application) routes() {
	app.router.HandleFunc("/", withLoggingMiddleware(withStatsMiddleware(app.handleFizzBuzz(), app.db)))
	app.router.HandleFunc("/stats", app.handldeFizzBuzzStats())
}

func (app Application) handldeFizzBuzzStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		sr := stats.New(app.db)
		res, err := sr.GetHighestHit()
		if err != nil {
			http.Error(w, "Error processing the request", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData{Count: 1, Error: err, Data: res})
	}
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
		json.NewEncoder(w).Encode(responseData{Count: len(res), Error: nil, Data: res})
	}
}
