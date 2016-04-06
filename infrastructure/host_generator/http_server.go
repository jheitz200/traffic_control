package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cihub/seelog"
	"github.com/jheitz200/host_generator/infrastructure/host_generator/utils"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
)

// Alerts ...
type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

// Alert ...
type Alert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

// StartServer the HTTP Server
func StartServer(c *utils.Config, to *client.Session) {
	version := "1.0"

	seelog.Debugf("Server available at %s", c.BindAddr)

	http.Handle(fmt.Sprintf("/api/%s/servers", version), Handle{TrafficOps: to, Handler: Server})

	seelog.Critical(http.ListenAndServe(c.BindAddr, nil))
	os.Exit(1)
}
