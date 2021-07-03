// Udostępnianie danych wydarzeń historycznych z serwisu XVI-wiek.pl
//
// Dokumentacja API
//
// Schemes: http
// BasePath: /api
// Version: 1.0.0
//
// Consumes:
// -
// Produces:
// - application/json
// swagger:meta
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi"
)

// SourceJSON type
type SourceJSON struct {
	Name string `json:"value" xml:"value"`
	URL  string `json:"url" xml:"url"`
}

// HistoricalEvent type
// swagger:response factsResponse
type HistoricalEvent struct {
	Date     string       `json:"date" xml:"date"`
	Title    string       `json:"title" xml:"title"`
	Content  string       `json:"content" xml:"content"`
	Location string       `json:"location" xml:"location"`
	Geo      string       `json:"geo" xml:"geo"`
	People   string       `json:"people" xml:"people"`
	Keywords string       `json:"keywords" xml:"keywords"`
	Sources  []SourceJSON `json:"sources" xml:"sources"`
}

// type SearchHistoricalEvent
type SearchHistoricalEvent struct {
	ID             string `json:"id" xml:"id"`
	Date           string `json:"date" xml:"date"`
	Day            string `json:"day" xml:"day"`
	Month          string `json:"month" xml:"month"`
	Title          string `json:"title" xml:"title"`
	ContentTwitter string `json:"content" xml:"content"`
}

//  type ShortHistoricalEvent
// swagger:response factsShortResponse
type ShortHistoricalEvent struct {
	Date           string `json:"date" xml:"date"`
	ContentTwitter string `json:"content" xml:"content"`
}

// type HealthcheckEvent
type HealthcheckEvent struct {
	Status  string `json:"status" xml:"status"`
	Version string `json:"version" xml:"version"`
}

//func clearField(value string) string {
//	value = prepareTextStyle(value, false)
//	return value
//}

func errorJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, _ := json.Marshal(map[string]string{"message": msg})
	w.Write(response)
}

func errorXML(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	response := fmt.Sprintf("<error><title>%d</title><content>%s</content></error>", code, msg)
	w.Write([]byte(response))
}

func toStructJSON(data interface{}) []HistoricalEvent {
	factStruct := data.(*[]Fact)

	factsJSON := []HistoricalEvent{}
	factJSON := HistoricalEvent{}
	sourceJSON := SourceJSON{}

	for _, item := range *factStruct {
		factJSON.Date = fmt.Sprintf("%02d-%02d-%04d", item.Day, item.Month, item.Year)
		factJSON.Title = item.Title
		factJSON.Content = prepareTextStyle(item.Content, true)
		factJSON.Location = item.Location
		factJSON.People = item.People
		factJSON.Geo = item.Geo
		factJSON.Keywords = item.Keywords

		for _, itemSource := range item.Sources {
			// podmiana id źródła na pełną informację o źródle
			// dla źródeł typu reference
			if itemSource.Type == "reference" {
				newValue, found := ReferenceMap[itemSource.Value]
				if found {
					itemSource.Value = newValue
				}
			}

			sourceJSON.Name = itemSource.Value

			if itemSource.Type == "reference" && itemSource.Page != "" {
				sourceJSON.Name += ", " + itemSource.Page
			}
			sourceJSON.URL = itemSource.URL
			factJSON.Sources = append(factJSON.Sources, sourceJSON)
			factJSON.Content = strings.Replace(factJSON.Content, "["+itemSource.ID+"]", "", -1)
		}

		factsJSON = append(factsJSON, factJSON)
	}

	return factsJSON
}

// toSearchStructJSON
func toSearchStructJSON(data interface{}) []SearchHistoricalEvent {
	factStruct := data.(*[]KeywordFact)

	factsJSON := []SearchHistoricalEvent{}
	factJSON := SearchHistoricalEvent{}

	for _, item := range *factStruct {
		factJSON.ID = item.ID
		factJSON.Date = item.Date
		// pola Day i Month uzupełniane na podstawie daty, wartości typu "05"
		// muszą być zamienione na "5"
		if len(item.Date) == 10 {
			factJSON.Day = item.Date[8:]
			tmp, err := strconv.Atoi(factJSON.Day)
			if err == nil {
				factJSON.Day = strconv.Itoa(tmp)
			}
			factJSON.Month = item.Date[5:7]
			tmp, err = strconv.Atoi(factJSON.Month)
			if err == nil {
				factJSON.Month = strconv.Itoa(tmp)
			}
		}
		factJSON.Title = item.Title
		factJSON.ContentTwitter = item.ContentTwitter

		factsJSON = append(factsJSON, factJSON)
	}

	return factsJSON
}

// toStructOneJSON
func toStructOneJSON(data interface{}) HistoricalEvent {
	factStruct := data.(Fact)
	factJSON := HistoricalEvent{}
	sourceJSON := SourceJSON{}

	factJSON.Date = fmt.Sprintf("%02d-%02d-%04d", factStruct.Day, factStruct.Month, factStruct.Year)
	factJSON.Title = factStruct.Title
	factJSON.Content = prepareTextStyle(factStruct.Content, true)
	factJSON.Location = factStruct.Location
	factJSON.People = factStruct.People
	factJSON.Geo = factStruct.Geo
	factJSON.Keywords = factStruct.Keywords

	for _, itemSource := range factStruct.Sources {
		// podmiana id źródła na pełną informację o źródle
		// dla źródeł typu reference
		if itemSource.Type == "reference" {
			newValue, found := ReferenceMap[itemSource.Value]
			if found {
				itemSource.Value = newValue
			}
		}
		sourceJSON.Name = itemSource.Value
		if itemSource.Type == "reference" && itemSource.Page != "" {
			sourceJSON.Name += ", " + itemSource.Page
		}
		sourceJSON.URL = itemSource.URL
		factJSON.Sources = append(factJSON.Sources, sourceJSON)
		factJSON.Content = strings.Replace(factJSON.Content, "["+itemSource.ID+"]", "", -1)
	}

	return factJSON
}

// toShortStructJSON
func toShortStructJSON(data interface{}) ShortHistoricalEvent {
	factStruct := data.(*[]Fact)
	factJSON := ShortHistoricalEvent{}

	choice := randomInt(0, len(*factStruct)-1)

	factJSON.Date = fmt.Sprintf("%02d-%02d-%04d", (*factStruct)[choice].Day, (*factStruct)[choice].Month, (*factStruct)[choice].Year)
	factJSON.ContentTwitter = (*factStruct)[choice].ContentTwitter

	return factJSON
}

// toHealthcheckStructJSON
func toHealthcheckStructJSON() HealthcheckEvent {
	factJSON := HealthcheckEvent{}

	factJSON.Status = "dostępny"
	factJSON.Version = "1.0.0"

	return factJSON
}

// factResponseJSON
func factResponseJSON(w http.ResponseWriter, code int, contentType string, data interface{}) {
	factsJSON := toStructJSON(data)

	if contentType == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(factsJSON)
	} else {
		response, _ := json.Marshal(factsJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// factSearchResponseJSON
func factSearchResponseJSON(w http.ResponseWriter, code int, contentType string, data interface{}) {
	factsJSON := toSearchStructJSON(data)

	if contentType == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(factsJSON)
	} else {
		response, _ := json.Marshal(factsJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// factResponseOneJSON
func factResponseOneJSON(w http.ResponseWriter, code int, contentType string, data interface{}) {
	factJSON := toStructOneJSON(data)

	if contentType == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(factJSON)
	} else {
		response, _ := json.Marshal(factJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// factShortResponseJSON
func factShortResponseJSON(w http.ResponseWriter, code int, contentType string, data interface{}) {
	factShortJSON := toShortStructJSON(data)

	if contentType == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(factShortJSON)
	} else {
		response, _ := json.Marshal(factShortJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// healthcheckResponseJSON
func healthcheckResponseJSON(w http.ResponseWriter, code int, contentType string) {
	healthcheckJSON := toHealthcheckStructJSON()

	if contentType == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(healthcheckJSON)
	} else {
		response, _ := json.Marshal(healthcheckJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// getFactsByDay
func (app *application) getFactsByDay(month, day int) (interface{}, bool) {

	var isCorrectDate bool = true

	if month == 2 && day > 29 {
		isCorrectDate = false
	}
	if (month == 4 || month == 6 || month == 9 || month == 11) && day > 30 {
		isCorrectDate = false
	}

	if !isCorrectDate {
		return nil, false
	}

	name := fmt.Sprintf("%02d-%02d", month, day)
	facts, ok := app.dataCache.Get(name)
	if ok {
		return facts, true
	}

	return nil, false
}

// swagger:route GET /day/{month}/{day} dzien listaWydarzen
// zwraca wydarzenia historyczne dla wskazanego dnia
// responses:
//   200: factsResponse

// apiFactsByDay
func (app *application) apiFactsByDay(w http.ResponseWriter, r *http.Request) {

	var isContentXML bool = r.Header.Get("Content-Type") == "application/xml"

	month, err := strconv.Atoi(chi.URLParam(r, "month"))
	if err != nil || month < 1 || month > 12 {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
		return
	}

	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil || day < 1 || day > 31 {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
		return
	}

	facts, ok := app.getFactsByDay(month, day)

	if ok {
		factResponseJSON(w, 200, r.Header.Get("Content-Type"), facts)
	} else {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
	}
}

// apiFactByDayAndID
func (app *application) apiFactByDayAndID(w http.ResponseWriter, r *http.Request) {

	var isContentXML bool = r.Header.Get("Content-Type") == "application/xml"

	month, err := strconv.Atoi(chi.URLParam(r, "month"))
	if err != nil || month < 1 || month > 12 {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
		return
	}

	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil || day < 1 || day > 31 {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
		return
	}

	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
		return
	}

	facts, ok := app.getFactsByDay(month, day)

	if ok {
		itemFacts := facts.(*[]Fact)
		for _, item := range *itemFacts {
			if item.ID == id {
				factResponseOneJSON(w, 200, r.Header.Get("Content-Type"), item)
				return
			}
		}
	}

	if isContentXML {
		errorXML(w, 404, "Błędne zapytanie lub brak danych")
	} else {
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
	}

}

// swagger:route GET /today today listaWydarzen
// zwraca wydarzenia historyczne dla bieżącego dnia
// responses:
//   200: factsResponse

// apiFactsToday
func (app *application) apiFactsToday(w http.ResponseWriter, r *http.Request) {
	name := fmt.Sprintf("%02d-%02d", int(time.Now().Month()), time.Now().Day())
	facts, ok := app.dataCache.Get(name)
	if ok {
		factResponseJSON(w, 200, r.Header.Get("Content-Type"), facts)
	} else {
		if r.Header.Get("Content-Type") == "application/xml" {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
	}
}

// swagger:route GET /short short listaWydarzen
// zwraca skrócony opis wydarzenia historyczngo dla bieżąceo dnia
// responses:
//   200: factsShortResponse

// apiFactsShort
func (app *application) apiFactsShort(w http.ResponseWriter, r *http.Request) {
	name := fmt.Sprintf("%02d-%02d", int(time.Now().Month()), time.Now().Day())
	facts, ok := app.dataCache.Get(name)
	if ok {
		factShortResponseJSON(w, 200, r.Header.Get("Content-Type"), facts)
	} else {
		if r.Header.Get("Content-Type") == "application/xml" {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
	}
}

// swagger:route GET /healthcheck healthcheck HealthcheckEvent
// zwraca status serwisu
// responses:
//   200: healthcheckResponse

// apiHealthcheck
func (app *application) apiHealthcheck(w http.ResponseWriter, r *http.Request) {
	healthcheckResponseJSON(w, 200, r.Header.Get("Content-Type"))
}

// apiSearchFacts
func (app *application) apiSearchFacts(w http.ResponseWriter, r *http.Request) {

	var isContentXML bool = r.Header.Get("Content-Type") == "application/xml"

	// searchQuery := chi.URLParam(r, "searchQuery")
	searchQuery := r.URL.Query().Get("searchQuery")

	if utf8.RuneCountInString(searchQuery) < 3 {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie, tekst do wyszukiwania powinien zawierać co najmniej 3 znaki")
		} else {
			errorJSON(w, 404, "Błędne zapytanie, tekst do wyszukiwania powinien zawierać co najmniej 3 znaki")
		}
		return
	}

	searchFacts, ok := app.searchInFacts(searchQuery)
	if ok {
		factSearchResponseJSON(w, 200, r.Header.Get("Content-Type"), searchFacts)
	} else {
		if isContentXML {
			errorXML(w, 404, "Błędne zapytanie lub brak danych")
		} else {
			errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		}
	}

}
