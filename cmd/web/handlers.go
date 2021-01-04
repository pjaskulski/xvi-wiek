package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

type templateDataInformation struct {
	NumberOfFacts int
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

	ts := app.templateCache["index.page.gohtml"]
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

	ts := app.templateCache["cytaty.page.gohtml"]
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

	ts := app.templateCache["ksiazki.page.gohtml"]
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

	data := &templateDataInformation{NumberOfFacts: numberOfFacts}

	ts := app.templateCache["informacje.page.gohtml"]
	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showFactsByDay
func (app *application) showFactsByDay(w http.ResponseWriter, r *http.Request) {
	month, err := strconv.Atoi(chi.URLParam(r, "month"))
	if err != nil || month < 1 || month > 12 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil || day < 1 || day > 31 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	var isCorrectDate bool = true

	if month == 2 && day > 29 {
		isCorrectDate = false
	}
	if (month == 4 || month == 6 || month == 9 || month == 11) && day > 30 {
		isCorrectDate = false
	}
	if !isCorrectDate {
		app.clientError(w, http.StatusNotFound)
		return
	}

	var data *templateDataFacts

	name := fmt.Sprintf("%02d-%02d", month, day)
	dayMonth := fmt.Sprintf("%d %s", day, monthName[month])

	facts, ok := app.dataCache.Get(name)
	if ok {
		data = &templateDataFacts{Today: dayMonth, Facts: facts.(*[]Fact)}
	} else {
		data = &templateDataFacts{Today: dayMonth, Facts: nil}
	}

	ts := app.templateCache["index.page.gohtml"]
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}
