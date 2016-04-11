package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jheitz200/host_generator/infrastructure/host_generator/fixtures"
	"github.com/jheitz200/host_generator/infrastructure/host_generator/utils"
	"github.com/jheitz200/test_helper"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
)

func TestServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/1.0/server?type=EDGE", nil)
	if err != nil {
		testHelper.Error(t, "%s", err)
	}

	r := fixtures.Servers()
	server := testHelper.ValidHTTPServer(r)
	defer server.Close()

	resp := httptest.NewRecorder()

	var httpClient http.Client
	to := client.Session{
		UserName:  "tester",
		Password:  "password",
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	var c utils.Config

	testHelper.Context(t, "Given the need to test a successful Get /api/1.0/server?type=EDGE")
	{
		Server(&c, &to, resp, req)

		if resp.Code != http.StatusOK {
			testHelper.Error(t, "\t Should get back \"200\" for HTTPStatusCode, got: \"%d\"", resp.Code)
		} else {
			testHelper.Success(t, "\t Should get back \"200\" for HTTPStatusCode")
		}

		if !strings.Contains(resp.Body.String(), "\"domainName\":\"albuquerque.nm.albuq.kabletown.com\"") {
			testHelper.Error(t, "\t Should get back \"town.net\" for \"domainName\", got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"town.net\" for \"domainName\"")
		}

		if !strings.Contains(resp.Body.String(), "\"hostName\":\"edge-alb-01\"") {
			testHelper.Error(t, "\t Should get back \"kable\" for \"hostName\", got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"kable\" for \"hostName\"")
		}

		if !strings.Contains(resp.Body.String(), "\"ipAddress\":\"10.10.10.10\"") {
			testHelper.Error(t, "\t Should get back \"127.0.0.1\" for \"ipAddress\", got \"%s\"", resp.Body.String())
		} else {
			testHelper.Success(t, "\t Should get back \"127.0.0.1\" for \"ipAddress\"")
		}
	}
}

func TestServersUnauthorized(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/1.0/servers?type=test", nil)
	if err != nil {
		testHelper.Error(t, "%s", err)
	}

	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	resp := httptest.NewRecorder()

	var httpClient http.Client
	to := client.Session{
		UserName:  "tester",
		Password:  "password",
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	var c utils.Config

	testHelper.Context(t, "Given the need to test a failed Get /api/1.0/servers")
	{
		Server(&c, &to, resp, req)

		if resp.Code != http.StatusUnauthorized {
			testHelper.Error(t, "\t Should get back \"401\" for HTTPStatusCode, got: \"%d\"", resp.Code)
		} else {
			testHelper.Success(t, "\t Should get back \"401\" for HTTPStatusCode")
		}
	}
}
