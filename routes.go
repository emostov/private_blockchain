package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes ...
type Routes []Route

//
// func NewRouter() *mux.Router {

// 	router := mux.NewRouter().StrictSlash(true)
// 	for _, route := range routes {
// 		router.
// 			Methods(route.Method).
// 			Path(route.Pattern).
// 			Name(route.Name).
// 			Handler(route.HandlerFunc)
// 	}

// 	return router
// }

// NewRouter ...
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}

var routes = Routes{
	Route{
		"Upload",
		"GET",
		"/Upload",
		Upload,
	},
	Route{
		"AskForBlock",
		"GET",
		"/block/{height}/{hash}",
		AskForBlock,
	},
	Route{
		"Show",
		"GET",
		"/show",
		ShowHandler,
	},
	Route{
		"Show",
		"POST",
		"/heartbeat/recieve",
		HeartBeatRecieve,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"peer",
		"POST",
		"/peer",
		Register,
	},
	Route{
		"Startserver",
		"GET",
		"/startserver",
		StartServer,
	},
}
