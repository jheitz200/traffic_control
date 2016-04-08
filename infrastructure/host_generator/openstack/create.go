package openstack

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/jheitz200/host_generator/infrastructure/host_generator/utils"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	// "github.com/rackspace/gophercloud/rackspace/compute/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
)

// Creator ...
type Creator struct {
	Provider      *gophercloud.ProviderClient
	Client        *gophercloud.ServiceClient
	Result        *tokens.CreateResult
	NetworkClient *gophercloud.ServiceClient
}

// ServerConfig ...
type ServerConfig interface {
	RAM() int
	NumCPUs() int
	Image() string
	KeyName() string
	Metadata() map[string]string
	FQDN() string
}

// New logs into an OpenStack cloud and returns a client instance
// that may be used to interact with the Compute API.
func New(c *utils.Config) (*Creator, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: c.OpenstackAuthURL,
		Username:         c.OpenstackUsername,
		Password:         c.OpenstackPassword,
		TenantName:       c.OpenstackTenantName,
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	identity := openstack.NewIdentityV2(provider)
	res := tokens.Create(identity, tokens.WrapOptions(opts))

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: c.OpenstackRegionName,
	})
	if client == nil || err != nil {
		return nil, err
	}

	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: c.OpenstackRegionName,
	})
	if networkClient == nil || err != nil {
		return nil, err
	}

	cr := Creator{
		Provider:      provider,
		Client:        client,
		Result:        &res,
		NetworkClient: networkClient,
	}
	return &cr, nil
}

// Create requests a server to be provisioned.
func (c Creator) Create(serverConfig ServerConfig, config *utils.Config) (*servers.Server, error) {
	var builder servers.CreateOptsBuilder

	n, err := c.PickNetwork(config)
	if err != nil {
		return nil, err
	}

	i, err := c.PickImage(config.Image)
	if err != nil {
		return nil, err
	}

	f, err := c.PickFlavor(serverConfig.NumCPUs(), serverConfig.RAM())
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, len(serverConfig.Metadata()))
	for k, v := range serverConfig.Metadata() {
		if len(v) < 256 {
			m[k] = v
		} else {
			seelog.Debugf(`Skipping metadata: "%s": "%s"`, k, v)
		}
	}

	opts := servers.CreateOpts{
		Name:      serverConfig.FQDN(),
		FlavorRef: f.ID,
		ImageRef:  i.ID,
		Metadata:  m,
		Networks:  n,
	}

	seelog.Debugf("Creating server %s", opts.Name)
	builder = keypairs.CreateOptsExt{
		CreateOptsBuilder: opts,
		KeyName:           serverConfig.KeyName(),
	}

	s, err := servers.Create(c.Client, builder).Extract()
	if err != nil {
		seelog.Debugf("Error creating %s: %s\n", opts.Name, err)
		return nil, err
	}

	return s, nil
}

// PickNetwork returns the specified network.
func (c Creator) PickNetwork(config *utils.Config) ([]servers.Network, error) {
	opts := networks.ListOpts{Name: config.OpenstackNetworkName}
	// Retrieve a pager (i.e. a paginated collection)
	pager, err := networks.List(c.NetworkClient, opts).AllPages()
	if err != nil {
		fmt.Printf("\n\n pager.AllPages %+v \n\n", err)
		return nil, err
	}

	networkList, err := networks.ExtractNetworks(pager)
	if err != nil {
		fmt.Printf("\n\n networks.ExtractNetworks %+v \n\n", err)
		return nil, err
	}

	var network []servers.Network
	for _, n := range networkList {
		s := servers.Network{
			UUID: n.ID,
		}
		network = append(network, s)
	}

	return network, nil
}

// PickFlavor returns the available flavors.
func (c Creator) PickFlavor(ncpu, mem int) (*flavors.Flavor, error) {
	// Retrieve a pager (i.e. a paginated collection)
	pager := flavors.ListDetail(c.Client, flavors.ListOpts{})

	p, err := pager.AllPages()
	if err != nil {
		return nil, err
	}

	flavorList, err := flavors.ExtractFlavors(p)
	if err != nil {
		return nil, err
	}

	for _, f := range flavorList {
		if f.VCPUs >= ncpu && f.RAM >= mem {
			return &f, nil
		}
	}

	return nil, err
}

// PickImage returns the specified image.
func (c Creator) PickImage(name string) (*images.Image, error) {
	// Retrieve a pager (i.e. a paginated collection)
	pager := images.ListDetail(c.Client, images.ListOpts{})

	p, err := pager.AllPages()
	if err != nil {
		return nil, err
	}

	imageList, err := images.ExtractImages(p)
	if err != nil {
		return nil, err
	}

	for _, i := range imageList {
		if i.Name == name {
			return &i, nil
		}
	}

	return nil, err
}
