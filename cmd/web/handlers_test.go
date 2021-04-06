package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// funkcja pomocnicza, ustawia środowisko testów, wczytuje dane
func createTestEnv() *application {
	// ścieżka w przypadku uruchamiania testów przez go test ./cmd/web
	dirExecutable = "../../"

	//bufor dla szablonów stron html
	templateCache, err := createTemplateCache(dirExecutable + "/ui/html/")
	if err != nil {
		log.Fatal(err)
	}

	// główna struktura aplikacji
	app := &application{
		errorLog:      log.New(os.Stdout, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:       log.New(os.Stdout, "INFO: \t", log.Ldate|log.Ltime),
		templateCache: templateCache,
		dataCache:     createDataCache(),
	}

	// wczytanie danych z plików yaml
	err = app.loadData(dirExecutable + "/data/")
	if err != nil {
		app.errorLog.Fatal(err)
	}

	return app
}

func TestHandlersCorrectPath(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	// definicje testów
	tests := []struct {
		route    string
		expected string
	}{
		{
			route:    "/",
			expected: `Co wydarzyło się`,
		},
		{
			route:    "/informacje",
			expected: "Główna strona serwisu prezentuje wydarzenia z bieżącego dnia",
		},
		{
			route:    "/cytaty",
			expected: "Teraz porządki francuskie chciał Henryk",
		},
		{
			route:    "/indeksy",
			expected: "indeks wydarzeń historycznych według lat",
		},
		{
			route:    "/indeksy/chronologia",
			expected: "1490",
		},
		{
			route:    "/indeksy/ludzie",
			expected: "Albrecht Fryderyk Hohenzollern",
		},
		{
			route:    "/indeksy/miejsca",
			expected: "Ansbach",
		},
		{
			route:    "/indeksy/slowa",
			expected: "dyplomacja",
		},
		{
			route:    "/pdf",
			expected: `<strong><a href="/static/pdf/xvi-wiek.pdf">xvi-wiek.pdf</a></strong> - zawartość serwisu jako ebook`,
		},
		{
			route:    "/kalendarz",
			expected: `Styczeń`,
		},
		{
			route:    "/ksiazki",
			expected: `Uwaga: opisy lub fragmenty opisów książek mogą pochodzić ze stron wydawców`,
		},
		{
			route:    "/dzien/3/22",
			expected: `Wołogoszcz`,
		},
		{
			route:    "/dzien/12/1",
			expected: "Brak danych", // dane na ten dzień jeszcze nie wprowadzone, zwraca 200 i informację
		},
		{
			route:    "/dzien/12/31",
			expected: "Brak danych", // dane na ten dzień jeszcze nie wprowadzone, zwraca 200 i informację
		},
		{
			route:    "/dzien/03/31",
			expected: "Jan Benedykt Solfa",
		},
	}

	for _, test := range tests {
		appTest.infoLog.Println("RUN handler: ", test.route)

		rs, err := ts.Client().Get(ts.URL + test.route)
		if err != nil {
			t.Fatal(err)
		}

		if rs.StatusCode != http.StatusOK {
			t.Errorf("oczekiwano %d; otrzymano %d", http.StatusOK, rs.StatusCode)
		}
		defer rs.Body.Close()

		body, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			t.Fatal(err)
		}

		bodyText := string(body)

		if !strings.Contains(bodyText, test.expected) {
			t.Errorf("handler %q brak oczekiwanego w 'body' tekstu: %q", test.route, test.expected)
		}
	}
}

func TestHandlersInvalidPath(t *testing.T) {
	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	// definicje testów
	tests := []struct {
		route  string
		status int
	}{
		{
			route:  "/about",
			status: 404,
		},
		{
			route:  "/dzien/2/30",
			status: 404,
		},
		{
			route:  "/dzien/4/31",
			status: 404,
		},
		{
			route:  "/dzien/15/19",
			status: 404,
		},
		{
			route:  "/dzien/1/3a",
			status: 404,
		},
		{
			route:  "/dzien/aaa/bbb",
			status: 404,
		},
		{
			route:  "/dzien/10",
			status: 404,
		},
	}

	for _, test := range tests {
		appTest.infoLog.Println("RUN handler: ", test.route)

		rs, err := ts.Client().Get(ts.URL + test.route)
		if err != nil {
			t.Fatal(err)
		}

		if rs.StatusCode != test.status {
			t.Errorf("http status, oczekiwano %d; otrzymano %d", test.status, rs.StatusCode)
		}
		defer rs.Body.Close()
	}
}
