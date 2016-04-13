package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/Comcast/traffic_control/traffic_ops/client/utils"
	"github.com/cihub/seelog"
)

// StartServer the HTTP Server
func StartServer(c *utils.Config, to *client.Session) {
	version := "1.0"

	seelog.Debugf("Server available at %s", c.BindAddr)

	http.Handle(fmt.Sprintf("/api/%s/elasticsearch", version), Handle{TrafficOps: to, Handler: Elasticsearch})
	http.Handle(fmt.Sprintf("/api/%s/logstash", version), Handle{TrafficOps: to, Handler: Logstash})
	http.Handle(fmt.Sprintf("/api/%s/servers", version), Handle{TrafficOps: to, Handler: Server})
	http.Handle(fmt.Sprintf("/api/%s/users", version), Handle{TrafficOps: to, Handler: User})
	http.Handle(fmt.Sprintf("/api/%s/zookeeper", version), Handle{TrafficOps: to, Handler: Zookeeper})

	seelog.Critical(http.ListenAndServe(c.BindAddr, nil))
	os.Exit(1)
}
