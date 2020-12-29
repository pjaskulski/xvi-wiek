package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type templateDataFacts struct {
	Today string
	Facts *[]Fact
}

type templateDataBooks struct {
	Books *[]Book
}

type templateDataQuotes struct {
	Quotes *[]Quote
}

var monthName = map[int]string{
	1:  "stycznia",
	2:  "lutego",
	3:  "marca",
	4:  "kwietnia",
	5:  "maja",
	6:  "czerwca",
	7:  "lipca",
	8:  "sierpnia",
	9:  "września",
	10: "października",
	11: "listopada",
	12: "grudnia",
}

// showFacts func
func (app *application) showFacts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFoundError(w)
		return
	}

	var data *templateDataFacts

	today := time.Now()
	name := fmt.Sprintf("%02d-%02d", int(today.Month()), today.Day())
	dayMonth := fmt.Sprintf("%d %s", today.Day(), monthName[int(today.Month())])
	facts, ok := app.dataCache.Get(name)
	if ok {
		data = &templateDataFacts{Today: dayMonth, Facts: facts.(*[]Fact)}
	} else {
		data = &templateDataFacts{Today: dayMonth, Facts: nil}
	}

	ts := app.templateCache["index.page.tmpl.html"]
	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showQuotes func
func (app *application) showQuotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	quotes, ok := app.dataCache.Get("quotes")
	if !ok {
		app.serverError(w, errors.New("błąd podczas odczytu bazy cytatów"))
		return
	}

	data := &templateDataQuotes{Quotes: quotes.(*[]Quote)}

	ts := app.templateCache["cytaty.page.tmpl.html"]
	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showBooks func
func (app *application) showBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	books, ok := app.dataCache.Get("books")
	if !ok {
		app.serverError(w, errors.New("błąd podczas odczytu bazy książek"))
		return
	}

	data := &templateDataBooks{Books: books.(*[]Book)}

	ts := app.templateCache["ksiazki.page.tmpl.html"]
	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showInformation func
func (app *application) showInformation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	ts := app.templateCache["informacje.page.tmpl.html"]
	err := ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
