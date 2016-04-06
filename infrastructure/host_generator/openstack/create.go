package openstack

import (
	"fmt"
	"os"

	"github.com/cihub/seelog"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

// Creator ...
type Creator struct {
	provider *gophercloud.ProviderClient
	client   *gophercloud.ServiceClient
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

// New ...
func New() (*Creator, error) {
	var cr Creator
	// NOTE: This requires OS_* env variables from openstack rc file downloaded from openstack dashboard.
	//  opts and provider are used in the Test* funcs
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	p, err := openstack.AuthenticatedClient(opts)
	if p == nil || err != nil {
		return nil, err
	}
	cr.provider = p

	c, err := openstack.NewComputeV2(cr.provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if c == nil || err != nil {
		return nil, err
	}
	cr.client = c
	return &cr, nil
}

// PickFlavor ...
func (c Creator) PickFlavor(ncpu, mem int) (*flavors.Flavor, error) {
	opts := flavors.ListOpts{}

	// Retrieve a pager (i.e. a paginated collection)
	pager := flavors.ListDetail(c.client, opts)

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

	return nil, fmt.Errorf("%v", err)
}

// PickImage ...
func (c Creator) PickImage(name string) (*images.Image, error) {
	opts := images.ListOpts{}

	// Retrieve a pager (i.e. a paginated collection)
	pager := images.ListDetail(c.client, opts)

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

	return nil, fmt.Errorf("%v", err)
}

// Create ...
func (c Creator) Create(cfg ServerConfig) (*servers.Server, error) {

	var builder servers.CreateOptsBuilder

	// TODO: Make this configurable
	im, err := c.PickImage("Comcast Cloud CentOS 7.1 x86_64 v1.0")
	if err != nil {
		return nil, err
	}

	f, err := c.PickFlavor(cfg.NumCPUs(), cfg.RAM())
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, len(cfg.Metadata()))
	for k, v := range cfg.Metadata() {
		if len(v) < 256 {
			m[k] = v
		} else {
			seelog.Debugf(`Skipping metadata: "%s": "%s"`, k, v)
		}
	}

	opts := servers.CreateOpts{
		Name:      cfg.FQDN(),
		FlavorRef: f.ID,
		ImageRef:  im.ID,
		Metadata:  m,
	}

	seelog.Debugf("Creating server %s", opts.Name)
	builder = keypairs.CreateOptsExt{
		CreateOptsBuilder: opts,
		KeyName:           cfg.KeyName(),
	}

	s, err := servers.Create(c.client, builder).Extract()
	if err != nil {
		seelog.Debugf("Error creating %s: %s\n", opts.Name, err)
		return nil, err
	}

	return s, nil
}
