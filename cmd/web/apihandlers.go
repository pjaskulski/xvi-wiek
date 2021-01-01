package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

// SourceJSON type
type SourceJSON struct {
	Name string `json:"value"`
	URL  string `json:"url"`
}

// FactJSON type
type FactJSON struct {
	Date     string       `json:"date"`
	Title    string       `json:"title"`
	Content  string       `json:"content"`
	Location string       `json:"location"`
	Geo      string       `json:"geo"`
	People   string       `json:"people"`
	Keywords string       `json:"keywords"`
	Sources  []SourceJSON `json:"sources"`
}

func clearField(value string) string {
	value = prepareTextStyle(value, false)
	return value
}

func errorJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, _ := json.Marshal(map[string]string{"message": msg})
	w.Write(response)
}

func toStructJSON(data interface{}) []FactJSON {
	factStruct := data.(*[]Fact)

	factsJSON := []FactJSON{}
	factJSON := FactJSON{}
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
			sourceJSON.Name = itemSource.Value
			sourceJSON.URL = itemSource.URL
			factJSON.Sources = append(factJSON.Sources, sourceJSON)
			factJSON.Content = strings.Replace(factJSON.Content, "["+itemSource.ID+"]", "", -1)
		}

		factsJSON = append(factsJSON, factJSON)
	}

	return factsJSON
}

func factResponseJSON(w http.ResponseWriter, code int, data interface{}) {
	factsJSON := toStructJSON(data)
	response, _ := json.Marshal(factsJSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// apiFactsByDay
func (app *application) apiFactsByDay(w http.ResponseWriter, r *http.Request) {
	month, err := strconv.Atoi(chi.URLParam(r, "month"))
	if err != nil || month < 1 || month > 12 {
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		return
	}

	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil || day < 1 || day > 31 {
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
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
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
		return
	}

	name := fmt.Sprintf("%02d-%02d", month, day)
	facts, ok := app.dataCache.Get(name)
	if ok {
		factResponseJSON(w, 200, facts)
	} else {
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
	}
}

// apiFactsToday
func (app *application) apiFactsToday(w http.ResponseWriter, r *http.Request) {
	name := fmt.Sprintf("%02d-%02d", int(time.Now().Month()), time.Now().Day())
	facts, ok := app.dataCache.Get(name)
	if ok {
		factResponseJSON(w, 200, facts)
	} else {
		errorJSON(w, 404, "Błędne zapytanie lub brak danych")
	}
}
