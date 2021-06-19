package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

// Config struct
type Config struct {
	Port string
}

// SliceFactsByKeyword type
type SliceFactsByKeyword struct {
	Keyword        string
	FactsByKeyword []KeywordFact
}

// SliceFactsByLocation type
type SliceFactsByLocation struct {
	Location        string
	FactsByLocation []LocationFact
}

// SliceFactsByPeople type
type SliceFactsByPeople struct {
	People        string
	FactsByPeople []PeopleFact
}

type application struct {
	errorLog         *log.Logger
	infoLog          *log.Logger
	templateCache    map[string]*template.Template
	dataCache        *cache.Cache
	FactsByYear      map[int][]YearFact
	FactsByPeople    map[string][]PeopleFact
	FactsByLocation  map[string][]LocationFact
	FactsByKeyword   map[string][]KeywordFact
	SFactsByKeyword  []SliceFactsByKeyword
	SFactsByLocation []SliceFactsByLocation
	SFactsByPeople   []SliceFactsByPeople
	TodaysQuote      quoteOfTheDay
	FactsForSearch   []SearchFact
}

// client type - struktura do zapisywania danych klienta po ip (limit api)
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	numberOfFacts int
	dirExecutable string
	waitgroup     = sync.WaitGroup{}
	lock          = sync.Mutex{}
	clients       = make(map[string]*client)
	isTesting     bool
)

func main() {
	// konfiguracja przez parametr z linii komend
	cfg := new(Config)
	flag.StringVar(&cfg.Port, "port", "8080", "port HTTP")
	flag.Parse()

	// ścieżka do pliku wykonywalnego
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	if isRunByRun() {
		dirExecutable = "."
	} else {
		dirExecutable = path.Dir(ex)
	}

	// logi z informacjami (->konsola) i błędami (->plik)
	infoLog := log.New(os.Stdout, "INFO: \t", log.Ldate|log.Ltime)
	fErr, err := os.OpenFile(dirExecutable+"/log/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fErr.Close()
	errorLog := log.New(fErr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)

	//bufor dla szablonów stron html
	templateCache, err := createTemplateCache(dirExecutable + "/ui/html/")
	if err != nil {
		log.Fatal(err)
	}

	// aplikacja
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		dataCache:     createDataCache(),
	}

	// wczytane danych do pamięci podręcznej
	err = app.loadData(dirExecutable + "/data/")
	if err != nil {
		app.errorLog.Fatal(err)
	}

	// start serwera http
	serwer := &http.Server{
		Addr:         ":" + cfg.Port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	// uruchomienie goroutine z funkcją czyszczenia mapy danych klientów (limitowanie api)
	if !isTesting {
		go LimitCleaner()
	}

	app.infoLog.Printf("Start serwera, port :%s", cfg.Port)

	go func() {
		err = serwer.ListenAndServe()
		if err != nil {
			app.errorLog.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	app.infoLog.Println("Otrzymano sygnał zatrzymania programu", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serwer.Shutdown(tc)
}
