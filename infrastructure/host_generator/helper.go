package main

import (
	"net/http"

	"github.com/jheitz200/host_generator/infrastructure/host_generator/utils"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
)

// Handle ...
type Handle struct {
	Config     *utils.Config
	TrafficOps *client.Session
	Handler    func(*utils.Config, *client.Session, http.ResponseWriter, *http.Request)
}

// ServeHTTP ...
func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler(h.Config, h.TrafficOps, w, r)
}
