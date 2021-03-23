package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
		route               string
		headerContentType   string
		status              int
		responseContentType string
	}{
		{
			route:               "/api/dzien/1/1",
			headerContentType:   "application/json",
			status:              200,
			responseContentType: "application/json",
		},
		{
			route:               "/api/dzien/1/1",
			headerContentType:   "application/xml",
			status:              200,
			responseContentType: "application/xml",
		},
		{
			route:               "/api/dzien/3/24",
			headerContentType:   "application/json",
			status:              200,
			responseContentType: "application/json",
		},
		{
			route:               "/api/today",
			headerContentType:   "application/json",
			status:              200,
			responseContentType: "application/json",
		},
		{
			route:               "/api/today",
			headerContentType:   "application/xml",
			status:              200,
			responseContentType: "application/xml",
		},
		{
			route:               "/api/short",
			headerContentType:   "application/json",
			status:              200,
			responseContentType: "application/json",
		},
		{
			route:               "/api/short",
			headerContentType:   "application/xml",
			status:              200,
			responseContentType: "application/xml",
		},
		{
			route:               "/api/dzien/2/30",
			headerContentType:   "application/json",
			status:              404,
			responseContentType: "application/json",
		},
		{
			route:               "/api/dzien/2/30",
			headerContentType:   "application/xml",
			status:              404,
			responseContentType: "application/xml",
		},
		{
			route:               "/api/dzien/100",
			headerContentType:   "application/json",
			status:              404,
			responseContentType: "application/json",
		},
		{
			route:               "/api/dzien/100",
			headerContentType:   "application/xml",
			status:              404,
			responseContentType: "application/xml",
		},
		{
			route:               "/api/dzien/1/400",
			headerContentType:   "application/json",
			status:              404,
			responseContentType: "application/json",
		},
		{
			route:               "/api/dzien/1/400",
			headerContentType:   "application/xml",
			status:              404,
			responseContentType: "application/xml",
		},
	}

	for _, test := range tests {
		appTest.infoLog.Println("API, GET: ", test.route, "Content-Type: ", test.headerContentType)

		client := ts.Client()
		req, _ := http.NewRequest("GET", ts.URL+test.route, nil)
		req.Header.Set("Content-Type", test.headerContentType)
		rs, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		contentType := rs.Header.Get("Content-type")
		if contentType != test.responseContentType {
			t.Errorf("content-type: oczekiwano %s; otrzymano %s", test.responseContentType, contentType)
		}

		if rs.StatusCode != test.status {
			t.Errorf("http status: oczekiwano %d; otrzymano %d", test.status, rs.StatusCode)
		}
		defer rs.Body.Close()

		body, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			t.Fatal(err)
		}

		if test.headerContentType == "application/json" && !json.Valid(body) {
			t.Errorf("GET %q otrzymano niepoprawną odpowiedź json: %q", test.route, string(body))
		}

		if test.headerContentType == "application/xml" && !IsValidXML(string(body)) {
			t.Errorf("GET %q otrzymano niepoprawną odpowiedź xml: %q", test.route, string(body))
		}
	}
}
