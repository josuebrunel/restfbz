package app

import (
	"net/http/httptest"
	"testing"
)

func TestApplication(t *testing.T) {
	app := New()
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
}
