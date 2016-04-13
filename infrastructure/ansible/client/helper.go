package main

import (
	"net/http"

	"github.com/Comcast/traffic_control/traffic_ops/client"
)

// Handle ...
type Handle struct {
	TrafficOps *client.Session
	Handler    func(*client.Session, http.ResponseWriter, *http.Request)
}

// ServeHTTP ...
func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler(h.TrafficOps, w, r)
}
