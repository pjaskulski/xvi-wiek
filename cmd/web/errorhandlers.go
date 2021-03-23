package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	clue := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, clue)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	if status == http.StatusNotFound {
		app.showNotFound(w, r)
		return
	}
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, r, http.StatusNotFound)
}
