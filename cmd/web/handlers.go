package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type templateDataFacts struct {
	Today            string
	TitleOfDay       string
	DescritpionOfDay string
	PrevNext         template.HTML
	Facts            *[]Fact
	KeyFacts         []KeywordFact
	TodayQuote       Quote
	DayUrlPath       string
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

type quoteOfTheDay struct {
	Date         string
	CurrentQuote Quote
}

type templateDataSearchResults struct {
	Query string
	Count int
	Facts *[]KeywordFact
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

func (app *application) getQuote() error {

	today := time.Now()
	day := fmt.Sprintf("%02d-%02d-%04d", today.Day(), int(today.Month()), today.Year())

	if app.TodaysQuote.Date != day {
		tmpQuotes, ok := app.dataCache.Get("quotes")
		if !ok {
			return errors.New("błąd podczas odczytu bazy cytatów")
		}

		quotes := tmpQuotes.(*[]Quote)
		id := randomInt(0, len(*quotes)-1)

		app.TodaysQuote.CurrentQuote = (*quotes)[id]
		app.TodaysQuote.Date = day
	}

	return nil
}

// showFacts func
func (app *application) showFacts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFoundError(w, r)
		return
	}

	var data *templateDataFacts
	var tKeyFacts []KeywordFact

	today := time.Now()
	name := fmt.Sprintf("%02d-%02d", int(today.Month()), today.Day())
	dayMonth := fmt.Sprintf("%d %s", today.Day(), monthName[int(today.Month())])
	facts, ok := app.dataCache.Get(name)
	if ok {
		var tmpFactTitle []string
		tFacts := facts.(*[]Fact)
		// zapis tytułów wydarzeń już prezentowanych na stronie by nie proponować
		// ich jeszcze raz
		for _, item := range *tFacts {
			tmpFactTitle = append(tmpFactTitle, item.Title)
		}

		// uzupełnienie listy propozycji kolejnych ciekawych wydarzeń
		// na podstawie słów kluczowych z już wyświetlancych wydarzeń
		for _, item := range *tFacts {
			if item.Keywords != "" {
				keywords := strings.Split(item.Keywords, ";")
				for _, keyword := range keywords {
					keyword = strings.TrimSpace(keyword)
					if facts, ok := app.FactsByKeyword[keyword]; ok {
						for _, kItem := range facts {
							if !inSlice(tmpFactTitle, kItem.Title) && !inSliceKeywordFact(tKeyFacts, kItem) {
								tKeyFacts = append(tKeyFacts, kItem)
							}
						}
					}
				}
			}
		}

		// jeżeli nie było słów kluczowych, uzupełnienie listy na podstawie
		// listy postaci (o ile jakieś występują w wyświetlanych już wydarzeniach)
		if len(tKeyFacts) == 0 {
			for _, item := range *tFacts {
				if item.People != "" {
					people := strings.Split(item.People, ";")
					for _, person := range people {
						person = strings.TrimSpace(person)
						if facts, ok := app.FactsByPeople[person]; ok {
							for _, kItem := range facts {
								if !inSlice(tmpFactTitle, kItem.Title) && !inSliceKeywordFact(tKeyFacts, KeywordFact(kItem)) {
									tKeyFacts = append(tKeyFacts, KeywordFact(kItem))
								}
							}
						}
					}
				}
			}
		}

		// liczbę podpowiadanych wydarzeń należy ograniczyć do trzech
		// wydarzenia są losowane, mogą więc różnić się przy każdym odświeżeniu strony
		if len(tKeyFacts) > 3 {

			var tmpThree []int
			var tmpKeyFacts []KeywordFact

			for len(tmpThree) < 3 {
				num := randomInt(0, len(tKeyFacts)-1)
				if !inSliceInt(tmpThree, num) {
					tmpThree = append(tmpThree, num)

				}
			}

			for _, n := range tmpThree {
				tmpKeyFacts = append(tmpKeyFacts, tKeyFacts[n])
			}

			tKeyFacts = nil
			tKeyFacts = append(tKeyFacts, tmpKeyFacts...)

			sort.Slice(tKeyFacts, func(i, j int) bool {
				return tKeyFacts[i].Date < tKeyFacts[j].Date
			})
		}

		// Quote Of The Day
		app.getQuote()

		data = &templateDataFacts{
			Today:            dayMonth,
			TitleOfDay:       "",
			DescritpionOfDay: "",
			Facts:            facts.(*[]Fact),
			KeyFacts:         tKeyFacts,
			TodayQuote:       app.TodaysQuote.CurrentQuote,
		}
	} else {
		data = &templateDataFacts{
			Today:            dayMonth,
			TitleOfDay:       "",
			DescritpionOfDay: "",
			Facts:            nil,
			KeyFacts:         nil,
			TodayQuote:       Quote{},
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
	dayUrlPath := fmt.Sprintf("%d/%d", month, day)
	dayMonth := fmt.Sprintf("%d %s", day, monthName[month])

	facts, ok := app.dataCache.Get(name)
	if ok {
		tmpFacts := facts.(*[]Fact)
		titleOfDay := (*tmpFacts)[0].Title
		descriptionOfDay := (*tmpFacts)[0].ContentTwitter
		data = &templateDataFacts{
			Today:            dayMonth,
			TitleOfDay:       titleOfDay,
			DescritpionOfDay: descriptionOfDay,
			PrevNext:         prevnext,
			Facts:            tmpFacts,
			DayUrlPath:       dayUrlPath,
		}
	} else {
		data = &templateDataFacts{
			Today:            dayMonth,
			TitleOfDay:       "",
			DescritpionOfDay: "",
			PrevNext:         prevnext,
			Facts:            nil,
			DayUrlPath:       dayUrlPath,
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
	err := ts.Execute(w, app.SFactsByPeople)
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

func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	ts := app.templateCache["search.page.gohtml"]
	err := ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) resultHandler(w http.ResponseWriter, r *http.Request) {

	var data *templateDataSearchResults

	u, err := url.Parse(r.URL.String())
	if err != nil {
		app.serverError(w, err)
		return
	}

	params := u.Query()
	searchQuery := params.Get("q")

	searchFacts, ok := app.searchInFacts(searchQuery)

	if ok {
		data = &templateDataSearchResults{
			Query: searchQuery,
			Count: len(*searchFacts),
			Facts: searchFacts,
		}
	} else {
		data = &templateDataSearchResults{
			Query: searchQuery,
			Count: 0,
			Facts: nil,
		}
	}

	ts := app.templateCache["wyniki.page.gohtml"]
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}
