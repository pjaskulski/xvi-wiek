package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
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
	Day       int      `yaml:"day"`
	Month     int      `yaml:"month"`
	Year      int      `yaml:"year"`
	Title     string   `yaml:"title"`
	Content   string   `yaml:"content"`
	Location  string   `yaml:"location"`
	Geo       string   `yaml:"geo"`
	People    string   `yaml:"people"`
	Keywords  string   `yaml:"keywords"`
	Image     string   `yaml:"image"`
	ImageInfo string   `yaml:"imageInfo"`
	Sources   []Source `yaml:"sources"`
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
	Image       string `yaml:"image"`
	Pages       int    `yaml:"pages"`
}

func createDataCache() *cache.Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return c
}

// readFact func
func readFact(filename string) (*[]Fact, error) {
	var result []Fact
	var fact Fact

	//filename, _ := filepath.Abs(fmt.Sprintf("./data/%02d-%02d.yaml", month, day))

	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(fileBuf)
	yamlDec := yaml.NewDecoder(r)

	for yamlDec.Decode(&fact) == nil {
		result = append(result, fact)
	}

	return &result, nil
}

// readQuote func
func readQuote() (*[]Quote, error) {
	var result []Quote
	var quote Quote

	filename, _ := filepath.Abs(fmt.Sprintf("./data/quotes.yaml"))

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
func readBook() (*[]Book, error) {
	var result []Book
	var book Book

	filename, _ := filepath.Abs(fmt.Sprintf("./data/books.yaml"))

	fileBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(fileBuf)
	yamlDec := yaml.NewDecoder(r)

	for yamlDec.Decode(&book) == nil {
		result = append(result, book)
	}

	return &result, nil
}

// loadData
func (app *application) loadData(path string) error {
	// wydarzenia
	fmt.Println("Wczytywanie bazy wydarzeń...")

	dataFiles, _ := filepath.Glob(filepath.Join(path, "*-*.yaml"))
	for _, tFile := range dataFiles {
		name := filenameWithoutExtension(filepath.Base(tFile))
		facts, err := readFact(tFile)
		if err != nil {
			return err
		}
		app.dataCache.Add(name, facts, cache.NoExpiration)
	}

	// cytaty
	fmt.Println("Wczytywanie bazy cytatów...")

	quotes, err := readQuote()
	if err != nil {
		return err
	}
	app.dataCache.Add("quotes", quotes, cache.NoExpiration)

	// książki
	fmt.Println("Wczytywanie bazy książek...")

	books, err := readBook()
	if err != nil {
		return err
	}
	app.dataCache.Add("books", books, cache.NoExpiration)

	return nil
}