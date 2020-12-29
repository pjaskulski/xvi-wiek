package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/patrickmn/go-cache"
)

// Config struct
type Config struct {
	Port string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	dataCache     *cache.Cache
}

func main() {
	// konfiguracja przez parametr z linii komend
	cfg := new(Config)
	flag.StringVar(&cfg.Port, "port", "8080", "port HTTP")
	flag.Parse()

	// logi z informacjami (->konsola) i błędami (->plik)
	infoLog := log.New(os.Stdout, "INFO: \t", log.Ldate|log.Ltime)
	fErr, err := os.OpenFile("./log/errors.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fErr.Close()
	errorLog := log.New(fErr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)

	//bufor dla szablonów stron html
	templateCache, err := createTemplateCache("./ui/html/")
	if err != nil {
		log.Fatal()
	}

	// aplikacja
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		dataCache:     createDataCache(),
	}

	// wczytane danych do pamięci podręcznej
	err = app.loadData("./data/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// start serwera http
	serwer := &http.Server{
		Addr:     ":" + cfg.Port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start serwera, port :%s", cfg.Port)
	err = serwer.ListenAndServe()
	errorLog.Fatal(err)
}
