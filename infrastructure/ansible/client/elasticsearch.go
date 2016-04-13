package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/cihub/seelog"
)

// Elasticsearch ...
func Elasticsearch(to *client.Session, w http.ResponseWriter, r *http.Request) {
	seelog.Debugf("Requesting URL '%s'", r.URL)

	t := "ELASTICSEARCH"
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

	var s []string
	for _, value := range servers {
		if value.Type != t {
			continue
		}

		if (value.HostName != "") && (value.DomainName != "") {
			s = append(s, fmt.Sprintf("%s.%s", value.HostName, value.DomainName))
		} else {
			s = append(s, value.IPAddress)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
