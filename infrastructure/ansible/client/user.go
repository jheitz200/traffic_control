package main

import (
	"encoding/json"
	"net/http"

	"github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/cihub/seelog"
)

// User ...
func User(to *client.Session, w http.ResponseWriter, r *http.Request) {
	seelog.Debugf("Requesting URL '%s'", r.URL)

	users, err := to.Users()
	if err != nil {
		seelog.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	seelog.Debugf("Retrieved %d users", len(users))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
