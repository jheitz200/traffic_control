package main

import (
	"encoding/json"
	"net/http"

	"github.com/cihub/seelog"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
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

	var u []client.User
	for _, value := range users {
		if value.PublicSSHKey == "" {
			continue
		}

		if (value.Role == "3") || (value.Role == "4") {
			u = append(u, value)
		} else {
			continue
		}
	}
	seelog.Debugf("Retrieved %d users", len(u))

	if len(u) == 1 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u[0])
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u)
	}
}
