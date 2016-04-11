/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package fixtures

import "github.com/jheitz200/traffic_control/traffic_ops/client"

// Servers returns a default ServerResponse to be used for testing.
func Servers() *client.ServerResponse {
	return &client.ServerResponse{
		Response: []client.Server{
			{
				ID:             "001",
				HostName:       "edge-alb-01",
				DomainName:     "albuquerque.nm.albuq.kabletown.com",
				CDNName:        "CDN-1",
				Type:           "EDGE",
				Profile:        "EDGE1_CDN_520",
				TCPPort:        "80",
				Rack:           "F-4/35",
				PhysLocation:   "Albuquerque",
				Cachegroup:     "albuquerque",
				IP6Address:     "2001:558:fe18:7::2/64",
				IP6Gateway:     "2001:558:fe18:7::1",
				IPAddress:      "10.10.10.10",
				IPGateway:      "10.10.10.10",
				IPNetmask:      "255.255.255.252",
				MgmtIPAddress:  "",
				MgmtIPNetmask:  "",
				MgmtIPGateway:  "",
				Status:         "REPORTED",
				XMPPID:         "edge-alb-01",
				XMPPPasswd:     "**********",
				IloIPAddress:   "10.10.10.10",
				IloUsername:    "vssilo",
				IloPassword:    "password",
				IloIPGateway:   "10.10.10.10",
				IloIPNetmask:   "255.255.255.192",
				InterfaceName:  "bond0",
				InterfaceMtu:   "9000",
				RouterPortName: "TenGigE0/3/0/10\t|\tBundle-Ether1000\t|\tPHY|10G|AGG-MEMBER|dtype:IPCDN-EDGE|rhost:edge-alb-01|rport:eth1|lagg:1000",
				RouterHostName: "ar01.albuquerque.nm.albuq.kabletown.com",
				LastUpdated:    "2015-03-27 17:00:30",
			},
		},
	}
}
