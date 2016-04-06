package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cihub/seelog"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
)


// @Title Logstash
// @Description queries Traffic Ops for a list of logstash servers per CDN
// @Accept  json
// @Success 200 {array}  string
// @Resource /logstash
// @Router /api/1.0/logstash	[get]
func Logstash(to *client.Session, w http.ResponseWriter, r *http.Request) {
	seelog.Debugf("Requesting URL '%s'", r.URL)

	t := "LOGSTASH"
	params := make(url.Values)
	params.Add("type", t)

	servers, err := to.ServersByType(params)
	if err != nil {
		seelog.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	cdn := r.URL.Query().Get("cdn")
	var s []string
	for _, value := range servers {
		if value.Type != t {
			continue
		}

		if (cdn != "") && (value.CDNName != cdn) {
			continue
		}

		if (value.HostName != "") && (value.DomainName != "") {
			s = append(s, fmt.Sprintf("%s.%s", value.HostName, value.DomainName))
		} else {
			s = append(s, value.IPAddress)
		}
	}
	seelog.Debugf("Retrieved %d %s servers", len(servers), t)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
