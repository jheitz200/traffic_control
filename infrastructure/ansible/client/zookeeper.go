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

// Zookeeper ...
func Zookeeper(to *client.Session, w http.ResponseWriter, r *http.Request) {
	seelog.Debugf("Requesting URL '%s'", r.URL)

	t := "ZOOKEEPER"
	params := make(url.Values)
	params.Add("type", t)

	servers, err := to.ServersByType(params)
	if err != nil {
		seelog.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	seelog.Debugf("Retrieved %d %s servers", len(servers), t)

	cdn := r.URL.Query().Get("cdn")
	var s []string
	for _, value := range servers {
		if value.Type != t {
			continue
		}

		if strings.ToLower(value.CDNName) != cdn {
			continue
		}

		s = append(s, fmt.Sprintf("%s.%s:2181", value.HostName, value.DomainName))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(strings.Join(s, ","))
}
