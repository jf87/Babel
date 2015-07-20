package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// provides the different types of devices
func TypesHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(a.library); err != nil {
		fmt.Println(err)
		return -1, err
	}

	return 200, nil
}

// request to create a link between a device and the type/location by the user
func LinkHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var device Device
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048567)) //limit size FIXME how big?
	if err != nil {
		return -1, err
	}
	if err := r.Body.Close(); err != nil {
		return -1, err
	}
	fmt.Printf("Body:\n%s\n\n", body)
	if err := json.Unmarshal(body, &device); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)                                             // not possible to process
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil { //we need err.Error() to access the error string
			return -1, err
		}
	}
	//reduce point list based on type of device that should be linked
	if len(device.Bacnet_types) > 0 && a.points_reduced == nil {
		fmt.Println("reducing points by type")
		var pr Points
		for _, v := range *a.points {
			var o Objects
			for _, va := range v.Objs {
				if contains(device.Bacnet_types, va.Props.Type) {
					o = append(o, va)
				}
			}
			if len(o) > 0 {
				v.Objs = o
				pr = append(pr, v)
			}
		}
		a.points_reduced = pr
		fmt.Printf("points reduced: %v\n", pr)
	} else if a.points_reduced == nil {
		a.points_reduced = *a.points
	}
	if len(a.points_reduced) > 0 {
		fmt.Println(device)
		fmt.Println(device.Value)

		go monitorBMS(a, device)

		var res Result
		i := <-su

		if i == 1 {
			res.Result = "Matched Point. Thank you :-)"
			a.points_reduced = *a.points
		} else if i == 0 {
			res.Result = "Could not find your input, please try again."
		} else {
			res.Result = "Reduced Points to" + strconv.Itoa(i)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			return -1, err
		}

		//go fakeActuation(a, device)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		var res Result
		res.Result = "error"
		if err := json.NewEncoder(w).Encode(res); err != nil {
			return -1, err
		}
	}
	return 200, nil
}

// tells the client if matching was a success, long polling
func SuccessHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var res Result
	fmt.Println("SuccessHandler")

	i := <-su

	if i == 1 {
		res.Result = "one"
		a.points_reduced = *a.points

	} else if i == 2 {
		res.Result = "multiple"
	} else if i == 3 {
		res.Result = "none"
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return -1, err
	}

	return 200, nil
}

// provides BMS points to smap driver. Points can dynamically change based on what happens
// in checkForSequence goroutine
func PointsHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {

	if active {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(a.points_reduced); err != nil {
			fmt.Println(err)
			return -1, err
		}
	} else {

		i := <-sync
		if i == 1 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(a.points_reduced); err != nil {
				fmt.Println(err)
				return -1, err
			}
		}
	}
	return 200, nil
}

func SyncHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {

	sync_smap <- true
	return 200, nil
}
