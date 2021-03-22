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

func TestInformacje(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/informacje")
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

	fragment := "Główna strona serwisu prezentuje wydarzenia z bieżącego dnia"
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestCytaty(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/cytaty")
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

	fragment := "Teraz porządki francuskie chciał Henryk"
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestIndeksy(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/indeksy")
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

	fragment := "indeks wydarzeń historycznych według lat"
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestEbook(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/pdf")
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

	fragment := `<strong><a href="/static/pdf/xvi-wiek.pdf">xvi-wiek.pdf</a></strong> - zawartość serwisu jako ebook`
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestKalendarz(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/kalendarz")
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

	fragment := `Styczeń`
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestKsiazki(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ksiazki")
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

	fragment := `Uwaga: opisy lub fragmenty opisów książek mogą pochodzić ze stron wydawców`
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestHome(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/")
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

	fragment := `Co wydarzyło się`
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}

func TestDay22Marca(t *testing.T) {

	// środowisko do testów
	appTest := createTestEnv()

	// serwer testowy
	ts := httptest.NewServer(appTest.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/dzien/3/22")
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

	fragment := `Wołogoszcz`
	if !strings.Contains(bodyText, fragment) {
		t.Errorf("brak oczekiwanego w 'body' tekstu: %q", fragment)
	}
}
