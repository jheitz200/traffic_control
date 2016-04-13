package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/cihub/seelog"
)

// Server formats Traffic Ops server information for use in puppet.
func Server(to *client.Session, w http.ResponseWriter, r *http.Request) {
	seelog.Debugf("Requesting URL '%s'", r.URL)

	switch r.Method {
	case "GET":
		get(to, w, r)
	default:
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(fmt.Errorf("Invalid method: %s", r.Method))
	}
}

func get(to *client.Session, w http.ResponseWriter, r *http.Request) {
	t := strings.ToUpper(r.URL.Query().Get("type"))
	servers, err := getServers(to, t)
	if err != nil {
		seelog.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	seelog.Debugf("Retrieved %d servers", len(servers))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servers)
}

func getServers(to *client.Session, t string) ([]client.Server, error) {
	params := make(url.Values)
	if t != "" {
		params.Add("type", t)
	}

	servers, err := to.ServersByType(params)
	if err != nil {
		return nil, err
	}

	return servers, nil
}
