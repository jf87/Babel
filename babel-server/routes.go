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
	Route{
		"Success",
		"GET",
		"/api/result",
		SuccessHandler,
	},
	Route{
		"Points",
		"GET",
		"/api/points",
		PointsHandler,
	},
	Route{
		"PointsInit",
		"GET",
		"/api/pointsinit",
		PointsInitHandler,
	},
	Route{
		"Sync",
		"GET",
		"/api/sync",
		SyncHandler,
	},
}
