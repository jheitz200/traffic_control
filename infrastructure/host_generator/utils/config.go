package utils

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cihub/seelog"
	"gopkg.in/yaml.v2"
)

// Config File for Traffic Ops Ansible Client
type Config struct {
	ConfigFile   string `yaml:"config_file"`
	SeelogConfig string `yaml:"seelog_config"`

	// Openstack Confgi Parameters
	OpenstackAuthURL     string `yaml:"OS_AUTH_URL"`
	OpenstackUsername    string `yaml:"OS_USERNAME"`
	OpenstackPassword    string `yaml:"OS_PASSWORD"`
	OpenstackTenantName  string `yaml:"OS_TENANT_NAME"`
	OpenstackRegionName  string `yaml:"OS_REGION_NAME"`
	OpenstackNetworkName string `yaml:"OS_NETWORK"`
	Image                string `json:"image"`

	// Traffic Ops Config Parameters
	TOServer   string   `yaml:"to_server"`
	TOUsername string   `yaml:"to_username"`
	TOPasswd   string   `yaml:"to_passwd"`
	Type       []string `yaml:"type"`
	BindAddr   string   `yaml:"bind_addr"`
}

// ValidationError is the error type returned when the
// config provided is missing required attributes.
type ValidationError struct {
	Field []string
}

// Error implements the error interface for the ValidateError type.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Missing required attributes in config: '%s'", strings.Join(e.Field, ", "))
}

// LoadConfig loads the configuration from the config file and command line arguments.
func LoadConfig(arguments []string) (*Config, error) {
	// Setup the defualt config values
	c := Config{
		ConfigFile:   "conf/config.yml",
		TOServer:     "http://127.0.0.1:3000",
		BindAddr:     "127.0.0.1:8080",
		SeelogConfig: "conf/seelog.xml",
	}

	// Load from config file.
	f := flag.NewFlagSet("ansible", -1)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&c.ConfigFile, "config-file", "conf/config.yml", "Path to Config File")
	f.Parse(arguments)

	_, err := os.Open(c.ConfigFile)
	if err == nil {
		if err := c.loadFile(c.ConfigFile); err != nil {
			return nil, err
		}
	}

	// Load from command line flags.
	if err := c.loadFlags(arguments); err != nil {
		return nil, err
	}

	if v := c.validate(); v != nil {
		return nil, &ValidationError{
			Field: v,
		}
	}

	return &c, nil
}

func (c *Config) loadFile(configFile string) error {
	seelog.Debugf("Loading config file: '%s'", configFile)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}

func (c *Config) loadFlags(arguments []string) error {

	f := flag.NewFlagSet("TO-Servers", flag.ExitOnError)

	f.StringVar(&c.OpenstackAuthURL, "os-auth-url", c.OpenstackAuthURL, "Openstack Auth URL")
	f.StringVar(&c.OpenstackUsername, "os-username", c.OpenstackUsername, "Openstack username")
	f.StringVar(&c.OpenstackPassword, "os-password", c.OpenstackPassword, "Openstack password")
	f.StringVar(&c.OpenstackTenantName, "os-tenant", c.OpenstackTenantName, "Openstack tenant name")
	f.StringVar(&c.OpenstackRegionName, "os-region", c.OpenstackRegionName, "Openstack region name")
	f.StringVar(&c.OpenstackNetworkName, "os-network", c.OpenstackNetworkName, "Openstack network name")
	f.StringVar(&c.Image, "image", c.Image, "Openstack Image")

	f.StringVar(&c.TOServer, "to-server", c.TOServer, "Traffic Ops server URL")
	f.StringVar(&c.TOUsername, "to-username", c.TOUsername, "Traffic Ops username")
	f.StringVar(&c.TOPasswd, "to-passwd", c.TOPasswd, "Traffic Ops password")
	f.StringVar(&c.BindAddr, "bind-addr", c.BindAddr, "HTTP Bind Address:Port")
	f.StringVar(&c.SeelogConfig, "seelog-config", c.SeelogConfig, "Location of Seelog config file.")

	// BEGIN IGNORED FLAGS
	var path string
	f.StringVar(&path, "config-file", "conf/config.yml", "Path to Config File")
	// END IGNORED FLAGS

	f.Parse(arguments)

	c.TOServer = strings.TrimRight(c.TOServer, "/") // remove any trailing slash
	return nil
}

func (c *Config) validate() []string {
	var e []string

	if c.TOUsername == "" {
		e = append(e, "to_username")
	}

	if c.TOPasswd == "" {
		e = append(e, "to_passwd")
	}

	err := setFromEnv(&c.OpenstackAuthURL, "OS_AUTH_URL")
	if err != nil {
		// err = fmt.Errorf("OpenstackAuthURL cannot be blank or %v", err)
		e = append(e, "OS_AUTH_URL")
	}

	err = setFromEnv(&c.OpenstackUsername, "OS_USERNAME")
	if err != nil {
		// err = fmt.Errorf("OpenstackUsername cannot be blank or %v", err)
		e = append(e, "OS_USERNAME")
	}

	err = setFromEnv(&c.OpenstackPassword, "OS_PASSWORD")
	if err != nil {
		// err = fmt.Errorf("OpenstackPassword cannot be blank or %v", err)
		e = append(e, "OS_PASSWORD")
	}

	err = setFromEnv(&c.OpenstackTenantName, "OS_TENANT_NAME")
	if err != nil {
		// err = fmt.Errorf("OpenstackTenantName cannot be blank or %v", err)
		e = append(e, "OS_TENANT_NAME")
	}

	err = setFromEnv(&c.OpenstackRegionName, "OS_REGION_NAME")
	if err != nil {
		// err = fmt.Errorf("OS_REGION_NAME cannot be blank or %v", err)
		e = append(e, "OS_REGION_NAME")
	}

	err = setFromEnv(&c.OpenstackNetworkName, "OS_NETWORK")
	if err != nil {
		// err = fmt.Errorf("OS_NETWORK cannot be blank or %v", err)
		e = append(e, "OS_NETWORK")
	}

	if c.Image == "" {
		e = append(e, "image")
	}

	return e
}

func setFromEnv(setting *string, envvar string) error {
	if *setting == "" {
		*setting = os.Getenv(envvar)
	}
	if *setting != "" {
		return nil
	}
	return fmt.Errorf("env var %s must be set", envvar)
}

// ConfigureSeelog ...
func ConfigureSeelog(seelogConfig string) error {
	logger, err := seelog.LoggerFromConfigAsFile(seelogConfig)
	if err != nil {
		err := fmt.Errorf("Error creating Logger from seelog file: %s", seelogConfig)
		return err
	}
	defer seelog.Flush()
	seelog.ReplaceLogger(logger)
	return nil
}

// String prints the config.
func (c *Config) String() string {
	var buf bytes.Buffer

	fmt.Fprint(&buf, "-------------------------------------------\n")
	fmt.Fprintf(&buf, "Config: \n")
	fmt.Fprintf(&buf, "\t Config File:             %s\n", c.ConfigFile)
	fmt.Fprintf(&buf, "\t Seelog Config:           %s\n", c.SeelogConfig)
	fmt.Fprintf(&buf, "Oppenstack: \n")
	fmt.Fprintf(&buf, "\t Auth URL:                %s\n", c.OpenstackAuthURL)
	fmt.Fprintf(&buf, "\t Username:                %s\n", c.OpenstackUsername)
	fmt.Fprintf(&buf, "\t Tenant Name:             %s\n", c.OpenstackTenantName)
	fmt.Fprintf(&buf, "\t Region Name:             %s\n", c.OpenstackRegionName)
	fmt.Fprintf(&buf, "\t Network Name:             %s\n", c.OpenstackNetworkName)
	fmt.Fprintf(&buf, "\t Image:                   %s\n", c.Image)
	fmt.Fprintf(&buf, "Traffic Ops: \n")
	fmt.Fprintf(&buf, "\t Server:                  %s\n", c.TOServer)
	fmt.Fprintf(&buf, "\t Username:                %s\n", c.TOUsername)
	fmt.Fprintf(&buf, "\t Type:                    %+v\n", c.Type)
	fmt.Fprintf(&buf, "\t HTTP Bind Address:Port:  %s\n", c.BindAddr)
	fmt.Fprintf(&buf, "----------------------------------------------------------------------\n")

	return buf.String()
}
