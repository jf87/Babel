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
		"Types",
		"GET",
		"/api/types",
		TypesHandler,
	},
	Route{
		"Link",
		"POST",
		"/api/link",
		LinkHandler,
	},
}
