package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestApiHandlers(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	// definicje testów
	tests := []struct {
		route       string
		status      int
		contentType string
	}{
		{
			route:       "/api/dzien/1/1",
			status:      200,
			contentType: "application/json",
		},
		{
			route:       "/api/dzien/3/24",
			status:      200,
			contentType: "application/json",
		},
		{
			route:       "/api/today",
			status:      200,
			contentType: "application/json",
		},
		{
			route:       "/api/short",
			status:      200,
			contentType: "application/json",
		},
		{
			route:       "/api/dzien/2/30",
			status:      404,
			contentType: "application/json",
		},
	}

	for _, test := range tests {
		appTest.infoLog.Println("API, GET: ", test.route)

		rs, err := ts.Client().Get(ts.URL + test.route)
		if err != nil {
			t.Fatal(err)
		}

		contentType := rs.Header.Get("Content-type")
		if contentType != test.contentType {
			t.Errorf("content-type: oczekiwano %s; otrzymano %s", test.contentType, contentType)
		}

		if rs.StatusCode != test.status {
			t.Errorf("http status: oczekiwano %d; otrzymano %d", test.status, rs.StatusCode)
		}
		defer rs.Body.Close()

		body, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !json.Valid(body) {
			t.Errorf("GET %q otrzymano niepoprawną odpowiedź json", test.route)
		}
	}
}
