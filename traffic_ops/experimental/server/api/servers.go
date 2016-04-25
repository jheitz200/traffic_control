// Copyright 2015 Comcast Cable Communications Management, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file was initially generated by gen_to_start.go (add link), as a start
// of the Traffic Ops golang data model

package api

import (
	"encoding/json"
	_ "github.com/Comcast/traffic_control/traffic_ops/experimental/server/output_format" // needed for swagger
	"github.com/jmoiron/sqlx"
	null "gopkg.in/guregu/null.v3"
	"log"
	"time"
)

type Servers struct {
	HostName       string       `db:"host_name" json:"hostName"`
	DomainName     string       `db:"domain_name" json:"domainName"`
	TcpPort        int64        `db:"tcp_port" json:"tcpPort"`
	XmppId         null.String  `db:"xmpp_id" json:"xmppId"`
	XmppPasswd     null.String  `db:"xmpp_passwd" json:"xmppPasswd"`
	InterfaceName  string       `db:"interface_name" json:"interfaceName"`
	IpAddress      null.String  `db:"ip_address" json:"ipAddress"`
	IpGateway      null.String  `db:"ip_gateway" json:"ipGateway"`
	Ip6Address     null.String  `db:"ip6_address" json:"ip6Address"`
	Ip6Gateway     null.String  `db:"ip6_gateway" json:"ip6Gateway"`
	InterfaceMtu   int64        `db:"interface_mtu" json:"interfaceMtu"`
	PhysLocation   string       `db:"phys_location" json:"physLocation"`
	Rack           null.String  `db:"rack" json:"rack"`
	Cachegroup     string       `db:"cachegroup" json:"cachegroup"`
	Status         string       `db:"status" json:"status"`
	UpdPending     bool         `db:"upd_pending" json:"updPending"`
	Cdn            string       `db:"cdn" json:"cdn"`
	MgmtIpAddress  null.String  `db:"mgmt_ip_address" json:"mgmtIpAddress"`
	MgmtIpGateway  null.String  `db:"mgmt_ip_gateway" json:"mgmtIpGateway"`
	IloIpAddress   null.String  `db:"ilo_ip_address" json:"iloIpAddress"`
	IloIpGateway   null.String  `db:"ilo_ip_gateway" json:"iloIpGateway"`
	IloUsername    null.String  `db:"ilo_username" json:"iloUsername"`
	IloPassword    null.String  `db:"ilo_password" json:"iloPassword"`
	RouterHostName null.String  `db:"router_host_name" json:"routerHostName"`
	RouterPortName null.String  `db:"router_port_name" json:"routerPortName"`
	CreatedAt      time.Time    `db:"created_at" json:"createdAt"`
	Links          ServersLinks `json:"_links" db:-`
}

type ServersLinks struct {
	Self             string           `db:"self" json:"_self"`
	ServersTypesLink ServersTypesLink `json:"servers_types" db:-`
	ProfilesLink     ProfilesLink     `json:"profiles" db:-`
}

// @Title getServersById
// @Description retrieves the servers information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    Servers
// @Resource /api/2.0
// @Router /api/2.0/servers/{id} [get]
func GetServer(hostName string, tcpPort int64, db *sqlx.DB) (interface{}, error) {
	ret := []Servers{}
	arg := Servers{}
	arg.HostName = hostName
	arg.TcpPort = tcpPort
	queryStr := "select *, concat('" + API_PATH + "servers', '/host_name/', host_name, '/tcp_port/', tcp_port) as self"
	queryStr += ", concat('" + API_PATH + "servers_types/', type) as servers_types_name_ref"
	queryStr += ", concat('" + API_PATH + "profiles/', profile) as profiles_name_ref"
	queryStr += " from servers WHERE host_name=:host_name AND tcp_port=:tcp_port"
	nstmt, err := db.PrepareNamed(queryStr)
	err = nstmt.Select(&ret, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nstmt.Close()
	return ret, nil
}

// @Title getServerss
// @Description retrieves the servers
// @Accept  application/json
// @Success 200 {array}    Servers
// @Resource /api/2.0
// @Router /api/2.0/servers [get]
func getServers(db *sqlx.DB) (interface{}, error) {
	ret := []Servers{}
	queryStr := "select *, concat('" + API_PATH + "servers', '/host_name/', host_name, '/tcp_port/', tcp_port) as self"
	queryStr += ", concat('" + API_PATH + "servers_types/', type) as servers_types_name_ref"
	queryStr += ", concat('" + API_PATH + "profiles/', profile) as profiles_name_ref"
	queryStr += " from servers"
	err := db.Select(&ret, queryStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// @Title postServers
// @Description enter a new servers
// @Accept  application/json
// @Param                 Body body     Servers   true "Servers object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/servers [post]
func postServer(payload []byte, db *sqlx.DB) (interface{}, error) {
	var v Servers
	err := json.Unmarshal(payload, &v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "INSERT INTO servers("
	sqlString += "host_name"
	sqlString += ",domain_name"
	sqlString += ",tcp_port"
	sqlString += ",xmpp_id"
	sqlString += ",xmpp_passwd"
	sqlString += ",interface_name"
	sqlString += ",ip_address"
	sqlString += ",ip_gateway"
	sqlString += ",ip6_address"
	sqlString += ",ip6_gateway"
	sqlString += ",interface_mtu"
	sqlString += ",phys_location"
	sqlString += ",rack"
	sqlString += ",cachegroup"
	sqlString += ",type"
	sqlString += ",status"
	sqlString += ",upd_pending"
	sqlString += ",profile"
	sqlString += ",cdn"
	sqlString += ",mgmt_ip_address"
	sqlString += ",mgmt_ip_gateway"
	sqlString += ",ilo_ip_address"
	sqlString += ",ilo_ip_gateway"
	sqlString += ",ilo_username"
	sqlString += ",ilo_password"
	sqlString += ",router_host_name"
	sqlString += ",router_port_name"
	sqlString += ",created_at"
	sqlString += ") VALUES ("
	sqlString += ":host_name"
	sqlString += ",:domain_name"
	sqlString += ",:tcp_port"
	sqlString += ",:xmpp_id"
	sqlString += ",:xmpp_passwd"
	sqlString += ",:interface_name"
	sqlString += ",:ip_address"
	sqlString += ",:ip_gateway"
	sqlString += ",:ip6_address"
	sqlString += ",:ip6_gateway"
	sqlString += ",:interface_mtu"
	sqlString += ",:phys_location"
	sqlString += ",:rack"
	sqlString += ",:cachegroup"
	sqlString += ",:type"
	sqlString += ",:status"
	sqlString += ",:upd_pending"
	sqlString += ",:profile"
	sqlString += ",:cdn"
	sqlString += ",:mgmt_ip_address"
	sqlString += ",:mgmt_ip_gateway"
	sqlString += ",:ilo_ip_address"
	sqlString += ",:ilo_ip_gateway"
	sqlString += ",:ilo_username"
	sqlString += ",:ilo_password"
	sqlString += ",:router_host_name"
	sqlString += ",:router_port_name"
	sqlString += ",:created_at"
	sqlString += ")"
	result, err := db.NamedExec(sqlString, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title putServers
// @Description modify an existing serversentry
// @Accept  application/json
// @Param   id              path    int     true        "The row id"
// @Param                 Body body     Servers   true "Servers object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/servers/{id}  [put]
func putServer(hostName string, tcpPort int64, payload []byte, db *sqlx.DB) (interface{}, error) {
	var arg Servers
	err := json.Unmarshal(payload, &arg)
	arg.HostName = hostName
	arg.TcpPort = tcpPort
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "UPDATE servers SET "
	sqlString += "host_name = :host_name"
	sqlString += ",domain_name = :domain_name"
	sqlString += ",tcp_port = :tcp_port"
	sqlString += ",xmpp_id = :xmpp_id"
	sqlString += ",xmpp_passwd = :xmpp_passwd"
	sqlString += ",interface_name = :interface_name"
	sqlString += ",ip_address = :ip_address"
	sqlString += ",ip_gateway = :ip_gateway"
	sqlString += ",ip6_address = :ip6_address"
	sqlString += ",ip6_gateway = :ip6_gateway"
	sqlString += ",interface_mtu = :interface_mtu"
	sqlString += ",phys_location = :phys_location"
	sqlString += ",rack = :rack"
	sqlString += ",cachegroup = :cachegroup"
	sqlString += ",type = :type"
	sqlString += ",status = :status"
	sqlString += ",upd_pending = :upd_pending"
	sqlString += ",profile = :profile"
	sqlString += ",cdn = :cdn"
	sqlString += ",mgmt_ip_address = :mgmt_ip_address"
	sqlString += ",mgmt_ip_gateway = :mgmt_ip_gateway"
	sqlString += ",ilo_ip_address = :ilo_ip_address"
	sqlString += ",ilo_ip_gateway = :ilo_ip_gateway"
	sqlString += ",ilo_username = :ilo_username"
	sqlString += ",ilo_password = :ilo_password"
	sqlString += ",router_host_name = :router_host_name"
	sqlString += ",router_port_name = :router_port_name"
	sqlString += ",created_at = :created_at"
	sqlString += " WHERE host_name=:host_name AND tcp_port=:tcp_port"
	result, err := db.NamedExec(sqlString, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title delServersById
// @Description deletes servers information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    Servers
// @Resource /api/2.0
// @Router /api/2.0/servers/{id} [delete]
func delServer(hostName string, tcpPort int64, db *sqlx.DB) (interface{}, error) {
	arg := Servers{}
	arg.HostName = hostName
	arg.TcpPort = tcpPort
	result, err := db.NamedExec("DELETE FROM servers WHERE host_name=:host_name AND tcp_port=:tcp_port", arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}
