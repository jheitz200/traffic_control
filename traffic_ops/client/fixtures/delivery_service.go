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

import "github.com/Comcast/traffic_control/traffic_ops/client"

// DeliveryServices returns a default DeliveryServiceResponse to be used for testing.
func DeliveryServices() *client.DeliveryServiceResponse {
	return &client.DeliveryServiceResponse{
		Response: []client.DeliveryService{
			client.DeliveryService{
				ID:                   "001",
				XMLID:                "ds-test",
				Active:               true,
				DSCP:                 "40",
				Signed:               false,
				QStringIgnore:        "1",
				GeoLimit:             "0",
				GeoProvider:          "0",
				DNSBypassTTL:         "30",
				Type:                 "HTTP",
				ProfileName:          "ds-123",
				CDNName:              "test-cdn",
				CCRDNSTTL:            "3600",
				GlobalMaxTPS:         "0",
				MaxDNSAnswers:        "0",
				MissLat:              "44.654321",
				MissLong:             "-99.123456",
				Protocol:             "0",
				IPv6RoutingEnabled:   true,
				RangeRequestHandling: "0",
				TRResponseHeaders:    "Access-Control-Allow-Origin: *",
				MultiSiteOrigin:      "0",
				DisplayName:          "Testing",
				InitialDispersion:    "1",
			},
		},
	}
}
