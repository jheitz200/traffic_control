package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jheitz200/test_helper"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
)

func TestServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/1.0/server", nil)
	if err != nil {
		testHelper.Error(t, "%s", err)
	}

	r := client.ServerResponse{
		Response: []client.Server{
			client.Server{
				DomainName:    "town.net",
				HostName:      "kable",
				InterfaceMtu:  "9000",
				InterfaceName: "bond0",
				IPAddress:     "127.0.0.1",
				IPGateway:     "127.0.0.1",
				IPNetmask:     "255.255.255.252",
				TCPPort:       "80",

				Cachegroup:   "mid-northeast",
				CDNName:      "cdn-1",
				ID:           "555",
				LastUpdated:  "2016-01-22 08:14:30",
				PhysLocation: "Denver",
				Profile:      "logstash",
				Rack:         "RR 119.02",
				Status:       "ONLINE",
				Type:         "LOGSTASH",
			},
		},
	}

	server := testHelper.ValidHTTPServer(r)
	resp := httptest.NewRecorder()

	var httpClient http.Client
	to := client.Session{
		UserName:  "tester",
		Password:  "password",
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Get /api/1.0/server")
	{
		Server(&to, resp, req)

		if resp.Code != 200 {
			testHelper.Error(t, "\t Should get back \"200\" for HTTPStatusCode, got: \"%d\"", resp.Code)
		} else {
			testHelper.Success(t, "\t Should get back \"200\" for HTTPStatusCode")
		}

		if !strings.Contains(resp.Body.String(), "\"domainName\":\"town.net\"") {
			testHelper.Error(t, "\t Should get back \"town.net\" for \"domainName\", got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"town.net\" for \"domainName\"")
		}

		if !strings.Contains(resp.Body.String(), "\"hostName\":\"kable\"") {
			testHelper.Error(t, "\t Should get back \"kable\" for \"hostName\", got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"kable\" for \"hostName\"")
		}

		if !strings.Contains(resp.Body.String(), "\"ipAddress\":\"127.0.0.1\"") {
			testHelper.Error(t, "\t Should get back \"127.0.0.1\" for \"ipAddress\", got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"127.0.0.1\" for \"ipAddress\"")
		}
	}
}

func TestServersUnauthorized(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/1.0/servers", nil)
	if err != nil {
		testHelper.Error(t, "%s", err)
	}

	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	resp := httptest.NewRecorder()

	var httpClient http.Client
	to := client.Session{
		UserName:  "tester",
		Password:  "password",
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a failed Get /api/1.0/servers")
	{
		Server(&to, resp, req)

		if resp.Code != 401 {
			testHelper.Error(t, "\t Should get back \"401\" for HTTPStatusCode, got: \"%d\"", resp.Code)
		} else {
			testHelper.Success(t, "\t Should get back \"401\" for HTTPStatusCode")
		}
	}
}
