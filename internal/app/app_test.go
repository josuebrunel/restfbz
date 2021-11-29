package app

import (
	"net/http/httptest"
	"os"
	"testing"
)

func TestApplication(t *testing.T) {
	dbfile := "test.db"
	port := "8989"
	cf := Config{dbfile, port}
	app := New(cf)
	if app.cf.Dbfile != dbfile {
		t.Fatal("Unexpected dbfile name")
	}
	if app.cf.Port != port {
		t.Fatal("Unexpected port number")
	}
	app.initdb()
	defer os.Remove(app.cf.Dbfile)
	app.makeMigrations()
	app.routes()
	// test with invalid method
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	app.router.ServeHTTP(w, req)
	if w.Result().StatusCode != 405 {
		t.Fatal("Unexpected status code returned")
	}
	// test with invalid values
	req = httptest.NewRequest("GET", "/?int1=3&int2=5&limit=toto&str1=fizz&str2=buzz", nil)
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, req)
	if w.Result().StatusCode != 400 {
		t.Fatal("Unexpected status code returned")
	}
	// test with wrong numbers of params
	req = httptest.NewRequest("GET", "/?int1=3&int2=5&limit=25&str1=fizz", nil)
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, req)
	if w.Result().StatusCode != 400 {
		t.Fatal("Unexpected status code returned")
	}
	// test with a limit value greater than 1M
	req = httptest.NewRequest("GET", "/?int1=3&int2=5&limit=1000001&str1=fizz&str2=buzz", nil)
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, req)
	if w.Result().StatusCode != 400 {
		t.Fatal("Unexpected status code returned")
	}
	// test with a correct request
	req = httptest.NewRequest("GET", "/?int1=3&int2=5&limit=25&str1=fizz&str2=buzz", nil)
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, req)
	if w.Result().StatusCode != 200 {
		t.Fatal("Unexpected status code returned")
	}
	// test stats endpoint
	req = httptest.NewRequest("GET", "/stats", nil)
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, req)
	if w.Result().StatusCode != 200 {
		t.Fatal("Unexpected status code returned")
	}
}
