package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/jheitz200/test_helper"
)

func TestZookeeper(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/1.0/zookeeper", nil)
	if err != nil {
		testHelper.Error(t, "%s", err)
	}

	r := client.ServerResponse{
		Response: []client.Server{
			client.Server{
				DomainName: "town.net",
				HostName:   "kable",
				IPAddress:  "127.0.0.1",
				CDNName:    "cdn-1",
				Type:       "LOGSTASH",
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

	testHelper.Context(t, "Given the need to test a successful Get /api/1.0/zookeeper")
	{
		Zookeeper(&to, resp, req)

		if resp.Code != 200 {
			testHelper.Error(t, "\t Should get back \"200\" for HTTPStatusCode, got: \"%d\"", resp.Code)
		} else {
			testHelper.Success(t, "\t Should get back \"200\" for HTTPStatusCode")
		}

		if resp.Body.String() == "" {
			testHelper.Error(t, "\t Should get back \"kable.town.net:2181\" for response body, got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"kable.town.net:2181\" for response body")
		}
	}
}

func TestZookeeperUnauthorized(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/1.0/zookeeper", nil)
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

	testHelper.Context(t, "Given the need to test a failed Get /api/1.0/zookeeper")
	{
		Zookeeper(&to, resp, req)

		if resp.Code != 401 {
			testHelper.Error(t, "\t Should get back \"401\" for HTTPStatusCode, got: \"%d\"", resp.Code)
		} else {
			testHelper.Success(t, "\t Should get back \"401\" for HTTPStatusCode")
		}
	}
}
