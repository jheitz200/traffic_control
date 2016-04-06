package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/cihub/seelog"
	"github.com/jheitz200/host_generator/infrastructure/host_generator/openstack"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
	gophercloud "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

// Server formats Traffic Ops server information for use in puppet.
func Server(to *client.Session, w http.ResponseWriter, r *http.Request) {
	seelog.Debugf("Requesting URL '%s'", r.URL)

	switch r.Method {
	case "POST":
		post(to, w, r)
	case "GET":
		get(to, w, r)
	default:
		err := fmt.Errorf("Invalid method: %s", r.Method)

		e := Alerts{
			Alerts: []Alert{
				Alert{
					Text:  err.Error(),
					Level: "error",
				},
			},
		}

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&e)
	}
}

func get(to *client.Session, w http.ResponseWriter, r *http.Request) {
	t := strings.ToUpper(r.URL.Query().Get("type"))
	if t == "" {
		err := fmt.Errorf("type parameter is required")
		seelog.Error(err)

		e := Alerts{
			Alerts: []Alert{
				Alert{
					Text:  err.Error(),
					Level: "error",
				},
			},
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&e)
		return
	}

	servers, err := getServers(to, t)
	if err != nil {
		seelog.Error(err)

		e := Alerts{
			Alerts: []Alert{
				Alert{
					Text:  err.Error(),
					Level: "error",
				},
			},
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&e)
		return
	}
	seelog.Debugf("Retrieved %d %s servers", len(servers), t)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servers)
}

func post(to *client.Session, w http.ResponseWriter, r *http.Request) {
	t := strings.ToUpper(r.URL.Query().Get("type"))
	if t == "" {
		err := fmt.Errorf("type must be provided, e.g. api/1.0/servers?type=foo")
		seelog.Error(err)

		e := Alerts{
			Alerts: []Alert{
				Alert{
					Text:  err.Error(),
					Level: "error",
				},
			},
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&e)
		return
	}

	servers, err := getServers(to, t)
	if err != nil {
		seelog.Error(err)

		e := Alerts{
			Alerts: []Alert{
				Alert{
					Text:  err.Error(),
					Level: "error",
				},
			},
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&e)
		return
	}
	seelog.Debugf("Retrieved %d %s servers from Traffic Ops", len(servers), t)

	cr, err := openstack.New()
	if err != nil {
		seelog.Error(err)

		e := Alerts{
			Alerts: []Alert{
				Alert{
					Text:  err.Error(),
					Level: "error",
				},
			},
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&e)
		return
	}

	instance := make(map[string]*gophercloud.Server, len(servers))
	osInstances, err := cr.ListServers(nil)

	for _, s := range servers {
		name := s.FQDN()
		for _, i := range osInstances {
			if i.Name == name {
				instance[name] = &i
				break
			}
		}

		if _, found := instance[name]; !found {
			// create it
			i, err := cr.Create(openstack.ServerConfig(s))
			if err != nil {
				seelog.Error(err)

				e := Alerts{
					Alerts: []Alert{
						Alert{
							Text:  err.Error(),
							Level: "error",
						},
					},
				}

				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(&e)
				return
			}
			instance[name] = i
			seelog.Debugf("Created %s on openstack", name)
			// TODO: remove this -- for developement only
			break
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(instance)
}

func getServers(to *client.Session, t string) ([]ServersResp, error) {
	params := make(url.Values)
	params.Add("type", t)

	servers, err := to.ServersByType(params)
	if err != nil {
		seelog.Error(err)
		return nil, err
	}

	resp, err := getParameters(to, servers)
	if err != nil {
		seelog.Error(err)
		return nil, err
	}

	return resp, nil
}

// ServersResp ...
type ServersResp struct {
	Server     client.Server
	Parameters map[string]string
}

func getParameters(to *client.Session, servers []client.Server) ([]ServersResp, error) {
	var resp []ServersResp

	profiles := make(map[string]map[string]string, len(servers))
	for _, s := range servers {
		p, err := Parameters(to, profiles, s.Profile)
		if err != nil {
			return nil, err
		}

		server := ServersResp{
			Server:     s,
			Parameters: p,
		}
		resp = append(resp, server)
	}

	return resp, nil
}

// Parameters returns a list of paramets for a given profile.
func Parameters(to *client.Session, profiles map[string]map[string]string, name string) (map[string]string, error) {
	if _, ok := profiles[name]; !ok {
		parameters, err := to.Parameters(name)
		if err != nil {
			return nil, err
		}

		params := make(map[string]string, len(parameters))
		for _, p := range parameters {
			params[p.Name] = p.Value
		}

		profiles[name] = params
	}

	return profiles[name], nil
}

// methods required for openstack.ServerConfig interface

// RAM ...
func (s ServersResp) RAM() int {
	ram := s.Parameters["ram"]
	if ram == "" {
		ram = "4GB"
	}
	i := strings.LastIndexFunc(ram, func(c rune) bool { return !unicode.IsNumber(c) })
	size, _ := strconv.Atoi(ram[:i])
	if size == 0 {
		size = 4
	}
	return size
}

// NumCPUs ...
func (s ServersResp) NumCPUs() int {
	vcpu, _ := strconv.Atoi(s.Parameters["VCPU"])
	if vcpu == 0 {
		vcpu = 1
	}
	return vcpu
}

// FQDN ...
func (s ServersResp) FQDN() string {
	return fmt.Sprintf("%s.%s", s.Server.HostName, s.Server.DomainName)
}

// Image ...
func (s ServersResp) Image() string {
	return s.Parameters["image"]
}

// KeyName ...
func (s ServersResp) KeyName() string {
	return s.Parameters["keyname"]
}

// Metadata ...
func (s ServersResp) Metadata() map[string]string {
	return s.Parameters
}
