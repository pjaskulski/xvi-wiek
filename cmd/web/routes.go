package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// procedury obsługi żądań http
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.showFacts)
	mux.HandleFunc("/cytaty", app.showQuotes)
	mux.HandleFunc("/ksiazki", app.showBooks)
	mux.HandleFunc("/informacje", app.showInformation)

	// obsługa plików statycznych, w katalogu i podkatalogach pusty plik index.html
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
