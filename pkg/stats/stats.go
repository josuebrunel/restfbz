package stats

import (
	"context"
	"database/sql"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

// HighestHitResult holds the result of the query returning the highest hit
type HighestHitResult struct {
	QS  string `json:"qs"`
	Hit int    `json:"hit"`
}

// StatRecorder records
type StatRecorder struct {
	db *sql.DB
}

// New initialize a new StatRecorder
func New(db *sql.DB) StatRecorder {
	return StatRecorder{db}
}

// execute run an sql statement within a transaction
func (sr StatRecorder) execute(statement string, values []interface{}) (sql.Result, error) {
	ctx := context.Background()
	tx, err := sr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(statement)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()
	tx.Commit()
	return res, err
}

// upsert inserts a mew record if it doesn't exist or update it
func (sr StatRecorder) upsert(uri string) (sql.Result, error) {
	stmt := `
		INSERT INTO stats (queryparams, hits) VALUES (@query, 1)
		ON CONFLICT (queryparams)
		DO UPDATE SET hits = hits + 1
	`
	return sr.execute(stmt, []interface{}{sql.Named("query", uri)})
}

// CreateTables creates the necessary tables for StatRecorder
func (sr StatRecorder) CreateTables() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			queryparams TEXT NOT NULL UNIQUE,
			hits INTEGER NOT NULL
		)
	`
	_, err := sr.execute(stmt, nil)
	return err
}

// CreateRecord insert or update a new record into the stats table
func (sr StatRecorder) CreateRecord(uri string) error {
	u, err := url.Parse(uri)
	qs := u.Query()
	if err != nil {
		return err
	}
	params := []string{"int1", "int2", "limit", "str1", "str2"}
	var sb strings.Builder
	for _, p := range params {
		if v, ok := qs[p]; ok {
			if sb.Len() > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(p)
			sb.WriteString("=")
			sb.WriteString(v[0])
		}
	}
	_, err = sr.upsert(sb.String())
	return err
}

// GetHighestHit returns the query params with the higest hist
func (sr StatRecorder) GetHighestHit() (HighestHitResult, error) {
	stmt := "SELECT queryparams, max(hits) FROM stats LIMIT 1"
	var res = HighestHitResult{}
	err := sr.db.QueryRow(stmt, 1).Scan(&res.QS, &res.Hit)
	return res, err
}
