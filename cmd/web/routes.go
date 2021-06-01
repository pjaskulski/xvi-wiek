package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/time/rate"
	//"github.com/go-chi/cors"
)

// LimitMiddleware func - limity zapytań API
func LimitMiddleware(next http.Handler) http.Handler {
	// globalny limit - 10 tokenów na zapytania, odświeżanych 5 razy na sekundę
	limiter := rate.NewLimiter(100, 20)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// limity tylko dla zapytań API
		if strings.Contains(r.URL.String(), "/api/") && !limiter.Allow() {
			errorJSON(w, http.StatusTooManyRequests, "przekroczono limit zapytań API")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(LimitMiddleware)

	//r.Use(middleware.Compress(5))

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))

	// procedury obsługi żądań http
	r.Get("/", app.showFacts)
	r.Get("/cytaty", app.showQuotes)
	r.Get("/kalendarz", app.showCalendar)
	r.Get("/ksiazki", app.showBooks)
	r.Get("/informacje", app.showInformation)
	r.Get("/indeksy", app.showIndexes)
	r.Get("/indeksy/chronologia", app.showChronology)
	r.Get("/indeksy/ludzie", app.showPeople)
	r.Get("/indeksy/miejsca", app.showLocation)
	r.Get("/indeksy/slowa", app.showKeyword)
	r.Get("/pdf", app.showPDF)

	r.Get("/dzien/{month}/{day}", app.showFactsByDay)
	// api
	r.Route("/api", func(r chi.Router) {
		r.Get("/dzien/{month}/{day}", app.apiFactsByDay)
		r.Get("/today", app.apiFactsToday)
		r.Get("/short", app.apiFactsShort) // zwraca skrócony opis dla Twittera
	})

	// obługa 404 not found
	r.NotFound(app.showNotFound)

	// obsługa plików statycznych, w katalogu i podkatalogach pusty plik index.html
	FileServer(r, "/static/", http.Dir(dirExecutable+"/ui/static/"))

	return r
}

// FileServer - obsługa plików statycznych
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer nie zezwala na żadne parametry adresu URL.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
