package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status, err := ah.h(ah.appContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page:
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

// HTTP handler that maps on / a function that takes the HTTP server response (w)
// and the client HTTP request (r) as its arguments. We then write to the response
// of the server, which then leads to HTTP data being send to the client.
func IndexHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	if _, err := fmt.Fprintf(w, "Hola! Qué tal?"); err != nil {
		return -1, err
	}

	return 200, nil
}

// get actuator types from server
func ActuatorsHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	act := vars["act"]
	i, err := strconv.ParseInt(act, 10, 0)
	if err != nil {
		return -1, err
	}
	fmt.Println(i)

	val := strconv.Itoa(actuators[i])
	fmt.Println(val)

	if _, err := fmt.Fprintf(w, val); err != nil {
		return -1, err
	}
	return 200, nil
}
