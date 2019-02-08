package main

import (
	"github.com/hecomp/session-based-authentication/handlers"
	"github.com/hecomp/session-based-authentication/redisclient"
	"github.com/prometheus/common/log"
	"net/http"
)


func main() {
	redisclient.InitCache()
	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/signin", handlers.Sigin)
	http.HandleFunc("/welcome", handlers.Welcome)
	http.HandleFunc("/refresh", handlers.Refresh)
	// Start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}