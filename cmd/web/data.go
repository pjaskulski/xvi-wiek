package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// Source type
type Source struct {
	ID      string `yaml:"id"`
	Value   string `yaml:"value"`
	URLName string `yaml:"urlName"`
	URL     string `yaml:"url"`
}

// Fact type
type Fact struct {
	ID             string `yaml:"id" validate:"required"`
	Day            int    `yaml:"day" validate:"required"`
	Month          int    `yaml:"month" validate:"required"`
	Year           int    `yaml:"year" validate:"required"`
	Title          string `yaml:"title" validate:"required"`
	Content        string `yaml:"content" validate:"required"`
	ContentHTML    template.HTML
	ContentTwitter string `yaml:"contentTwitter"`
	Location       string `yaml:"location"`
	Geo            string `yaml:"geo"`
	GeoHTML        template.HTML
	People         string `yaml:"people"`
	Keywords       string `yaml:"keywords"`
	Image          string `yaml:"image"`
	ImageInfo      string `yaml:"imageInfo"`
	ImageHTML      template.HTML
	Sources        []Source `yaml:"sources"`
}

// Quote type
type Quote struct {
	Content string `yaml:"content"`
	Source  string `yaml:"source"`
}

// Book type
type Book struct {
	Author      string `yaml:"author"`
	Title       string `yaml:"title"`
	Year        string `yaml:"year"`
	Pubhause    string `yaml:"pubhause"`
	Where       string `yaml:"where"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	ISBN        string `yaml:"ISBN"`
	URL         string `yaml:"URL"`
	URLHTML     template.HTML
	URLName     string `yaml:"URLName"`
	Image       string `yaml:"image"`
	ImageHTML   template.HTML
	Pages       int `yaml:"pages"`
}

// YearFact type
type YearFact struct {
	Date      string
	DateMonth string
	Title     string
	URLHTML   template.HTML
}

// PeopleFact type
type PeopleFact struct {
	Date      string
	DateMonth string
	Title     string
	URLHTML   template.HTML
}

// LocationFact type
type LocationFact struct {
	Date      string
	DateMonth string
	Title     string
	URLHTML   template.HTML
}

// KeywordFact type
type KeywordFact struct {
	Date      string
	DateMonth string
	Title     string
	URLHTML   template.HTML
}

// DayFactTable map
var DayFactTable map[string]bool

// Validate func
func (f *Fact) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func createDataCache() *cache.Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return c
}

// readFact func
func (app *application) readFact(filename string) {
	var result []Fact
	var fact Fact

	defer waitgroup.Done()

	name := filenameWithoutExtension(filepath.Base(filename))

	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	r := bytes.NewReader(fileBuf)
	yamlDec := yaml.NewDecoder(r)

	yamlErr := yamlDec.Decode(&fact)

	for yamlErr == nil {
		err = fact.Validate()
		if err != nil {
			app.errorLog.Println("file:", filepath.Base(filename)+",", "error:", err)
		}

		fact.ContentHTML = template.HTML(prepareFactHTML(fact.Content, fact.ID, fact.Sources))
		fact.ImageHTML = template.HTML(prepareImageHTML(fact.Image, fact.ImageInfo))
		if fact.Geo != "" {
			fact.GeoHTML = template.HTML(prepareGeoHTML(fact.Geo))
		}

		// uzupełnienie indeksu lat FactsByYear
		tmpYear := &YearFact{}
		tmpYear.Date = fmt.Sprintf("%04d-%02d-%02d", fact.Year, fact.Month, fact.Day)
		tmpYear.DateMonth = fmt.Sprintf("%d %s", fact.Day, monthName[fact.Month])
		tmpYear.Title = fact.Title
		tmpYear.URLHTML = template.HTML(prepareFactLinkHTML(fact.Month, fact.Day, fact.ID))

		lock.Lock()
		if facts, ok := app.FactsByYear[fact.Year]; ok {
			facts = append(facts, *tmpYear)
			app.FactsByYear[fact.Year] = facts
		} else {
			facts := make([]YearFact, 0)
			facts = append(facts, *tmpYear)
			app.FactsByYear[fact.Year] = facts
		} // FactsByYear
		lock.Unlock()

		// uzupełnienie indeksu postaci FactsByPeople
		if fact.People != "" {
			tmpPeople := &PeopleFact{}
			tmpPeople.Date = fmt.Sprintf("%04d-%02d-%02d", fact.Year, fact.Month, fact.Day)
			tmpPeople.DateMonth = fmt.Sprintf("%d %s %d", fact.Day, monthName[fact.Month], fact.Year)
			tmpPeople.Title = fact.Title
			tmpPeople.URLHTML = template.HTML(prepareFactLinkHTML(fact.Month, fact.Day, fact.ID))
			persons := strings.Split(fact.People, ";")
			for _, person := range persons {
				person = strings.TrimSpace(person)
				lock.Lock()
				if facts, ok := app.FactsByPeople[person]; ok {
					facts = append(facts, *tmpPeople)
					app.FactsByPeople[person] = facts
				} else {
					facts := make([]PeopleFact, 0)
					facts = append(facts, *tmpPeople)
					app.FactsByPeople[person] = facts
				}
				lock.Unlock()
			}
		} // FactsByPeople

		// uzupełnienie indeksu lat FactsByLocation
		tmpLocation := &LocationFact{}
		tmpLocation.Date = fmt.Sprintf("%04d-%02d-%02d", fact.Year, fact.Month, fact.Day)
		tmpLocation.DateMonth = fmt.Sprintf("%d %s %d", fact.Day, monthName[fact.Month], fact.Year)
		tmpLocation.Title = fact.Title
		tmpLocation.URLHTML = template.HTML(prepareFactLinkHTML(fact.Month, fact.Day, fact.ID))
		location := strings.TrimSpace(fact.Location)
		if location != "" {
			lock.Lock()
			if facts, ok := app.FactsByLocation[location]; ok {
				facts = append(facts, *tmpLocation)
				app.FactsByLocation[location] = facts
			} else {
				facts := make([]LocationFact, 0)
				facts = append(facts, *tmpLocation)
				app.FactsByLocation[location] = facts
			}
			lock.Unlock()
		} // FactsByLocation

		// uzupełnienie indeksu postaci FactsByKeyword
		if fact.Keywords != "" {
			tmpKeyword := &KeywordFact{}
			tmpKeyword.Date = fmt.Sprintf("%04d-%02d-%02d", fact.Year, fact.Month, fact.Day)
			tmpKeyword.DateMonth = fmt.Sprintf("%d %s %d", fact.Day, monthName[fact.Month], fact.Year)
			tmpKeyword.Title = fact.Title
			tmpKeyword.URLHTML = template.HTML(prepareFactLinkHTML(fact.Month, fact.Day, fact.ID))
			keywords := strings.Split(fact.Keywords, ";")
			for _, keyword := range keywords {
				keyword = strings.TrimSpace(keyword)
				lock.Lock()
				if facts, ok := app.FactsByKeyword[keyword]; ok {
					facts = append(facts, *tmpKeyword)
					app.FactsByKeyword[keyword] = facts
				} else {
					facts := make([]KeywordFact, 0)
					facts = append(facts, *tmpKeyword)
					app.FactsByKeyword[keyword] = facts
				}
				lock.Unlock()
			}
		} // FactsByKeyword

		result = append(result, fact)

		yamlErr = yamlDec.Decode(&fact)
	}
	// jeżeli był błądw pliku yaml, inny niż koniec pliku to zapis w logu
	if yamlErr != nil && yamlErr.Error() != "EOF" {
		app.errorLog.Println("file:", filepath.Base(filename)+",", "error:", yamlErr)
	}

	lock.Lock()
	numberOfFacts += len(result)
	DayFactTable[name] = true
	app.dataCache.Add(name, &result, cache.NoExpiration)
	lock.Unlock()
}

// readQuote func
func (app *application) readQuote() (*[]Quote, error) {
	var result []Quote
	var quote Quote

	filename, _ := filepath.Abs(dirExecutable + "/data/quotes.yaml")

	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(fileBuf)
	yamlDec := yaml.NewDecoder(r)

	for yamlDec.Decode(&quote) == nil {
		result = append(result, quote)
	}

	return &result, nil
}

// readBook func
func (app *application) readBook() (*[]Book, error) {
	var result []Book
	var book Book

	filename := dirExecutable + "/data/books.yaml"

	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(fileBuf)
	yamlDec := yaml.NewDecoder(r)

	for yamlDec.Decode(&book) == nil {
		book.ImageHTML = template.HTML(prepareBookHTML(book.Image))
		book.URLHTML = template.HTML(prepareBookURLHTML(book.URL, book.URLName))
		result = append(result, book)
	}

	return &result, nil
}

// loadData - wczytuje podczas startu serwera dane do struktur w pamięci operacyjnej
func (app *application) loadData(path string) error {
	// wydarzenia
	app.infoLog.Printf("Wczytywanie bazy wydarzeń...")
	start := time.Now()

	// mapa z listą dni - czy dla danego dnia istnieją wydarzenia w bazie
	DayFactTable = make(map[string]bool)

	// mapa dla indeksu lat
	app.FactsByYear = make(map[int][]YearFact)
	// mapa dla indeksu postaci
	app.FactsByPeople = make(map[string][]PeopleFact)
	// mapa dla indeksu miejsc
	app.FactsByLocation = make(map[string][]LocationFact)
	// mapa dla indeksu słów kluczowych
	app.FactsByKeyword = make(map[string][]KeywordFact)

	dataFiles, _ := filepath.Glob(filepath.Join(path, "*-*.yaml"))
	if len(dataFiles) > 0 {
		for _, tFile := range dataFiles {
			waitgroup.Add(1)
			go app.readFact(tFile)
		}
		waitgroup.Wait()
	}

	// sortowanie wydarzeń historycznych dla postaci
	for person, facts := range app.FactsByPeople {
		sort.Slice(facts, func(i, j int) bool {
			return facts[i].Date < facts[j].Date
		})
		app.FactsByPeople[person] = facts
	}

	// sortowanie wydarzeń historycznych dla miejsc
	for location, facts := range app.FactsByLocation {
		sort.Slice(facts, func(i, j int) bool {
			return facts[i].Date < facts[j].Date
		})
		app.FactsByLocation[location] = facts
	}

	// sortowanie wydarzeń historycznych dla słów kluczowych
	for keyword, facts := range app.FactsByKeyword {
		sort.Slice(facts, func(i, j int) bool {
			return facts[i].Date < facts[j].Date
		})
		app.FactsByKeyword[keyword] = facts
	}

	// dodatkowy slice dla szablonu
	for key, facts := range app.FactsByKeyword {
		temp := SliceFactsByKeyword{Keyword: key, FactsByKeyword: facts}
		app.SFactsByKeyword = append(app.SFactsByKeyword, temp)
	}

	cl := collate.New(language.Polish)
	sort.SliceStable(app.SFactsByKeyword, func(i, j int) bool {
		return cl.CompareString(app.SFactsByKeyword[i].Keyword, app.SFactsByKeyword[j].Keyword) == -1
	})

	// cytaty
	app.infoLog.Printf("Wczytywanie bazy cytatów...")

	quotes, err := app.readQuote()
	if err != nil {
		return err
	}
	app.dataCache.Add("quotes", quotes, cache.NoExpiration)

	// książki
	app.infoLog.Printf("Wczytywanie bazy książek...")

	books, err := app.readBook()
	if err != nil {
		return err
	}
	app.dataCache.Add("books", books, cache.NoExpiration)

	elapsed := time.Since(start)
	app.infoLog.Printf("Czas wczytywania danych: %s", elapsed)

	return nil
}

/* dayFact - funkcja zwraca fragment html z linkiem jeżeli dla danego dnia są wydarzenia
   historyczne w bazie, lub sam numer dnia (o szarym kolorze) jeżeli ich nie ma.
   Wykorzystywana w kalendarzu.
*/
func dayFact(month int, day int) template.HTML {
	name := fmt.Sprintf("%02d-%02d", month, day)

	if DayFactTable[name] {
		result := fmt.Sprintf(`<a href="/dzien/%d/%d">%d</a>`, month, day, day)
		return template.HTML(result)
	}

	return template.HTML(fmt.Sprintf(`<span style="color: DarkGrey;">%d</span>`, day))
}
