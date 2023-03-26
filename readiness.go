package main

import "net/http"

// add /healthcheckz endpoint
func HandlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK NIGGA"))
}
