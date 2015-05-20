package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	Route{
		"Actuators",
		"GET",
		`/api/actuators/{act:[0-9]+}`,
		ActuatorsHandler,
	},
	Route{
		"ActuatorsSet",
		"POST",
		`/api/actuators/{act:[0-9]+}`,
		SetActuatorsHandler,
	},
}
