package stats_test

import (
	"database/sql"
	"os"
	"restfbz/pkg/stats"
	"testing"
)

func TestStatRecorder(t *testing.T) {
	dbname := "test.db"
	file, err := os.Create(dbname)
	defer os.Remove(dbname)
	if err != nil {
		t.Fatal("Failed to create test db")
	}
	defer file.Close()
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		t.Fatal("Failt to open db file")
	}
	sr := stats.New(db)
	sr.CreateTables()
	if _, err := os.Stat(dbname); os.IsNotExist(err) {
		t.Fatalf("DB File %s has not been created", dbname)
	}
	sr.CreateRecord("?int1=3&int2=5&limit=25&str1=foo&str2=bar")
	sr.CreateRecord("?int1=3&int2=5&limit=25&str1=foo&str2=bar")
	sr.CreateRecord("?int1=4&int2=8&limit=50&str1=buzz&str2=fizz")
	sr.CreateRecord("?int1=3&int2=5&limit=25&str1=foo&str2=bar")
	var count int
	db.QueryRow("select count(*) from stats", 1).Scan(&count)
	if count != 2 {
		t.Fatal("Number of records not equal to 2")
	}
	res, err := sr.GetHighestHit()
	if err != nil {
		t.Fatal("Failed to retrieve highest count")
	}
	if res.QS != "int1=3&int2=5&limit=25&str1=foo&str2=bar" || res.Hit != 3 {
		t.Fatalf("Unexpected returned value %v", res)
	}
}
