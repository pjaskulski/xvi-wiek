package main

import (
	"expvar"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/time/rate"
	//"github.com/go-chi/cors"
)

// LimitMiddleware func - limity zapytań API
func LimitMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// limity tylko dla zapytań API
		if strings.Contains(r.URL.String(), "/api/") {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				errorJSON(w, http.StatusBadRequest, err.Error())
				return
			}

			// Lock the mutex to prevent this code from being executed concurrently
			lock.Lock()

			// sprawdza czy ip clienta jest już w mapie, jeżeli nie tworzy nowy limiter
			if _, found := clients[ip]; !found {
				// 10 tokenów na zapytania, odświeżanych 5 razy na sekundę
				clients[ip] = &client{limiter: rate.NewLimiter(5, 10)}
			}

			// data i czas, kiedy ostatnio widziano kienta z danego ip
			clients[ip].lastSeen = time.Now()

			// jeżeli limit nie pozwala na obsługę kolejnego zapytania zwracany jest błąd 429
			if !clients[ip].limiter.Allow() {
				lock.Unlock()
				errorJSON(w, http.StatusTooManyRequests, "przekroczono limit zapytań API")
				return
			}

			// unlock the mutex before calling the next handler in the chain
			lock.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// LimitCleaner - funkcja czyści mapę klientów ze starych wpisów (starszych
// niż dwie godziny)
func LimitCleaner() {
	for {
		time.Sleep(time.Hour)

		lock.Lock()

		for ip, client := range clients {
			if time.Since(client.lastSeen) > 2*time.Hour {
				delete(clients, ip)
			}
		}

		lock.Unlock()
	}
}

// func enableCORS - middleware ustawia w nagłówku Access-Control-Allow-Origin
// (tylko dla zapytań API)
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/api/") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
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
	r.Use(enableCORS)
	r.Use(app.appMetrics)
	// bez limitów podczas uruchamiania testów
	if !isTesting {
		r.Use(LimitMiddleware)
	}

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
	r.Get("/szukaj", app.searchHandler)
	r.Get("/wyniki", app.resultHandler)

	r.Get("/dzien/{month}/{day}", app.showFactsByDay)

	if isRunByRun() {
		expvar.Publish("goroutines", expvar.Func(func() interface{} {
			return runtime.NumGoroutine()
		}))
		r.Handle("/monitor", expvar.Handler())
	}

	// api
	r.Route("/api", func(r chi.Router) {
		r.Get("/dzien/{month}/{day}", app.apiFactsByDay)
		r.Get("/today", app.apiFactsToday)
		r.Get("/short", app.apiFactsShort)        // zwraca skrócony opis dla Twittera
		r.Get("/healthcheck", app.apiHealthcheck) // testowy endpoint - status api
		r.Get("/fact/{month}/{day}/{id}", app.apiFactByDayAndID)
		// r.Get("/find/{searchQuery}")
		r.Get("/find", app.apiSearchFacts)
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

// appMetric func - metryki dla aplikacji
func (app *application) appMetrics(next http.Handler) http.Handler {
	totalRequestIn := expvar.NewInt("total_request_in")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestIn.Add(1)

		next.ServeHTTP(w, r)
	})
}
