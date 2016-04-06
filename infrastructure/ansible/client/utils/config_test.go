package utils

import (
	"testing"

	"github.com/jheitz200/test_helper"

	"gopkg.in/yaml.v2"
)

// Ensures that a configuration can be deserialized from TOML.
func TestConfig(t *testing.T) {

	data := `
    to_server:  http://127.0.0.1:3000
    to_username:    admin
    to_passwd:  password

    type:
    - test

    bind_addr:    "127.0.0.1:555"
    seelog_config:  "conf/seelog-testHelper.xml"
  `

	testHelper.Context(t, "Given the need to test a successful Get /api/1.0/logstash")
	{
		var c Config
		err := yaml.Unmarshal([]byte(data), &c)

		if err != nil {
			testHelper.Error(t, "\t Should be able to load the config file")
		} else {
			testHelper.Success(t, "\t Should be able to load the config file")
		}

		if c.TOServer != "http://127.0.0.1:3000" {
			testHelper.Error(t, "\t Should get back \"http://127.0.0.1:3000\" for TOServer, got: \"%s\"", c.TOServer)
		} else {
			testHelper.Success(t, "\t Should get back \"http://127.0.0.1:3000\" for TOServer")
		}

		if c.TOUsername != "admin" {
			testHelper.Error(t, "\t Should get back \"admin\" for TOUsername, got \"%s\"", c.TOUsername)
		} else {
			testHelper.Success(t, "\t Should get back \"admin\" for TOUsername")
		}

		if c.TOPasswd != "password" {
			testHelper.Error(t, "\t Should get back \"password\" for TOPasswd, got \"%s\"", c.TOPasswd)
		} else {
			testHelper.Success(t, "\t Should get back \"password\" for TOPasswd")
		}

		if c.Type[0] != "test" {
			testHelper.Error(t, "\t Should get back \"test\" for Type, got \"%s\"", c.Type)
		} else {
			testHelper.Success(t, "\t Should get back \"test\" for Type")
		}

		if c.BindAddr != "127.0.0.1:555" {
			testHelper.Error(t, "\t Should get back \"127.0.0.1:555\" for BindAddr, got \"%s\"", c.BindAddr)
		} else {
			testHelper.Success(t, "\t Should get back \"127.0.0.1:555\" for BindAddr")
		}

		if c.SeelogConfig != "conf/seelog-testHelper.xml" {
			testHelper.Error(t, "\t Should get back \"conf/seelog-testHelper.xml\" for SeelogConf, got \"%s\"", c.SeelogConfig)
		} else {
			testHelper.Success(t, "\t Should get back \"conf/seelog-testHelper.xml\" for SeelogConf")
		}
	}
}

// Ensures that a configuration can be parsed from the command line flags.
func TestConfigCommandLine(t *testing.T) {
	flags := []string{
		"-to-server", "0.0.0.0",
		"-to-username", "tester",
		"-bind-addr", "0.0.0.0:1111",
	}

	testHelper.Context(t, "Given the need to test a successful Get /api/1.0/logstash")
	{
		var c Config
		err := c.loadFlags(flags)

		if err != nil {
			testHelper.Error(t, "\t Should be able to load the config file")
		} else {
			testHelper.Success(t, "\t Should be able to load the config file")
		}

		if c.TOServer != "0.0.0.0" {
			testHelper.Error(t, "\t Should get back \"0.0.0.0\" for TOServer, got: \"%s\"", c.TOServer)
		} else {
			testHelper.Success(t, "\t Should get back \"0.0.0.0\" for TOServer")
		}

		if c.TOUsername != "tester" {
			testHelper.Error(t, "\t Should get back \"tester\" for TOUsername, got \"%s\"", c.TOUsername)
		} else {
			testHelper.Success(t, "\t Should get back \"tester\" for TOUsername")
		}

		if c.BindAddr != "0.0.0.0:1111" {
			testHelper.Error(t, "\t Should get back \"0.0.0.0:1111\" for BindAddr, got \"%s\"", c.BindAddr)
		} else {
			testHelper.Success(t, "\t Should get back \"0.0.0.0:1111\" for BindAddr")
		}
	}
}
