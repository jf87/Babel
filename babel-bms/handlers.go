package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	if _, err := fmt.Fprintf(w, "Babel BMS Test Server"); err != nil {
		return -1, err
	}

	return 200, nil
}

// get actuator types from server
func ActuatorsHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	act := vars["act"]
	i, err := strconv.Atoi(act)
	if err != nil {
		return -1, err
	}
	fmt.Println(i)

	if i >= len(a.actuators) {
		err := fmt.Errorf("Actuatorindex is out of range.")
		return 404, err
	}

	val := strconv.Itoa(a.actuators[i])
	fmt.Println(val)

	if _, err := fmt.Fprintf(w, val); err != nil {
		return -1, err
	}
	return 200, nil
}

// set the value of an actuator, e.g., by using the Anroid app
func SetActuatorsHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var s Setpoint
	vars := mux.Vars(r)
	act := vars["act"]
	i, err := strconv.Atoi(act)
	if err != nil {
		return -1, err
	}

	if i >= len(a.actuators) {
		err := fmt.Errorf("Actuatorindex is out of range.")
		return -1, err
	}

	a.useractuators[i] = i

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 102400))
	if err != nil {
		return -1, err
	}

	if err := r.Body.Close(); err != nil {
		return -1, err
	}
	fmt.Printf("%s", body)
	if err := json.Unmarshal(body, &s); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)                                             // not possible to process
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil { //we need err.Error() to access the error string
			return -1, err
		}
	}
	fmt.Println(s)

	a.actuators[i] = int(s.Value)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		return -1, nil
	}

	return 200, nil
}
