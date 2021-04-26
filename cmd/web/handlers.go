package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type templateDataFacts struct {
	Today      string
	TitleOfDay string
	PrevNext   template.HTML
	Facts      *[]Fact
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
		app.notFoundError(w, r)
		return
	}

	var data *templateDataFacts

	today := time.Now()
	name := fmt.Sprintf("%02d-%02d", int(today.Month()), today.Day())
	dayMonth := fmt.Sprintf("%d %s", today.Day(), monthName[int(today.Month())])
	facts, ok := app.dataCache.Get(name)
	if ok {
		data = &templateDataFacts{
			Today:      dayMonth,
			TitleOfDay: "",
			Facts:      facts.(*[]Fact),
		}
	} else {
		data = &templateDataFacts{
			Today:      dayMonth,
			TitleOfDay: "",
			Facts:      nil,
		}
	}

	ts := app.templateCache["index.page.gohtml"]
	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showCalendar func
func (app *application) showCalendar(w http.ResponseWriter, r *http.Request) {

	ts := app.templateCache["kalendarz.page.gohtml"]
	err := ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// showQuotes func
func (app *application) showQuotes(w http.ResponseWriter, r *http.Request) {

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
		app.clientError(w, r, http.StatusNotFound)
		return
	}

	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil || day < 1 || day > 31 {
		app.clientError(w, r, http.StatusNotFound)
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
		app.showNotFound(w, r)
		return
	}

	prevnext := template.HTML(getPrevNextHTML(month, day))

	var data *templateDataFacts

	name := fmt.Sprintf("%02d-%02d", month, day)
	dayMonth := fmt.Sprintf("%d %s", day, monthName[month])

	facts, ok := app.dataCache.Get(name)
	if ok {
		tmpFacts := facts.(*[]Fact)
		titleOfDay := (*tmpFacts)[0].Title
		data = &templateDataFacts{
			Today:      dayMonth,
			TitleOfDay: titleOfDay,
			PrevNext:   prevnext,
			Facts:      tmpFacts,
		}
	} else {
		data = &templateDataFacts{
			Today:      dayMonth,
			TitleOfDay: "",
			PrevNext:   prevnext,
			Facts:      nil,
		}
	}

	ts := app.templateCache["day.page.gohtml"]
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showIndexes func
func (app *application) showIndexes(w http.ResponseWriter, r *http.Request) {

	ts := app.templateCache["indeksy.page.gohtml"]
	err := ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// showChronology func
func (app *application) showChronology(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html")

	ts := app.templateCache["chronologia.page.gohtml"]
	err := ts.Execute(w, app.FactsByYear)
	if err != nil {
		app.serverError(w, err)
	}
}

// showPeople func
func (app *application) showPeople(w http.ResponseWriter, r *http.Request) {

	ts := app.templateCache["ludzie.page.gohtml"]
	err := ts.Execute(w, app.FactsByPeople)
	if err != nil {
		app.serverError(w, err)
	}
}

// showLocation func
func (app *application) showLocation(w http.ResponseWriter, r *http.Request) {

	ts := app.templateCache["miejsca.page.gohtml"]
	err := ts.Execute(w, app.SFactsByLocation)
	if err != nil {
		app.serverError(w, err)
	}
}

// showKeyword func
func (app *application) showKeyword(w http.ResponseWriter, r *http.Request) {

	ts := app.templateCache["slowa.page.gohtml"]
	err := ts.Execute(w, app.SFactsByKeyword)
	if err != nil {
		app.serverError(w, err)
	}
}

// showPDF func
func (app *application) showPDF(w http.ResponseWriter, r *http.Request) {

	ts := app.templateCache["pdf.page.gohtml"]
	err := ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// showNotFound func
func (app *application) showNotFound(w http.ResponseWriter, r *http.Request) {
	var isAPI bool = false

	if strings.Contains(r.URL.String(), `/api/`) {
		isAPI = true
	}

	if r.Header.Get("Content-Type") == "application/xml" {
		errorXML(w, 404, "Błędne zapytanie lub brak danych")
	} else if r.Header.Get("Content-Type") == "application/json" || isAPI {
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
	} else {
		ts := app.templateCache["notfound.page.gohtml"]
		w.WriteHeader(404)
		err := ts.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
		}
	}
}
