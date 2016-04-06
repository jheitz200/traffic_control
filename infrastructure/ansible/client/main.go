// @APIVersion 1.0
// @APITitle Ansible Client
// @APIDescription Connects to Traffic Ops and grabs any data needed to run Ansible playbooks.
// @Contact _T&P--VADER-IPCDN-ENG_@comcast.com
// @License Apache License, Version 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0
package main

import (
	"os"

	"github.com/cihub/seelog"
	"github.com/jheitz200/ansible_client/traffic_ops/ansible_client/utils"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
)

func main() {
	config, err := utils.LoadConfig(os.Args[1:])
	if err != nil {
		seelog.Critical(err)
		os.Exit(1)
	}
	seelog.Debug(config.String())

	if err := utils.ConfigureSeelog(config.SeelogConfig); err != nil {
		seelog.Critical(err)
		os.Exit(1)
	}

	session, err := client.Login(config.TOServer, config.TOUsername, config.TOPasswd, true)
	if err != nil {
		seelog.Critical(err)
		os.Exit(1)
	}

	StartServer(config, session)
}
