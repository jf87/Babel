package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	if _, err := fmt.Fprintf(w, "Welcome to Babel Server"); err != nil {
		return -1, err
	}

	return 200, nil
}

// get actuator types from server
func TypesHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(a.library); err != nil {
		fmt.Println(err)
		return -1, err
	}

	return 200, nil
}

// request to create a link between an actuator and the type/location by the user
func LinkHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var device Device
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048567)) //limit size FIXME how big?
	if err != nil {
		return -1, err
	}
	if err := r.Body.Close(); err != nil {
		return -1, err
	}
fmt.Printf("%s\n\n", body)
	if err := json.Unmarshal(body, &device); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)                                             // not possible to process
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil { //we need err.Error() to access the error string
			return -1, err
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(device); err != nil {
		return -1, err
	}

	go checkForSequence(a, device)

	fmt.Println(device)
	return 200, nil
}
