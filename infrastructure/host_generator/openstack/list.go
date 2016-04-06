package openstack

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

// ListServers ...
func (c Creator) ListServers(f func(s servers.Server) bool) ([]servers.Server, error) {
	// We have the option of filtering the server list. If we want the full
	// collection, leave it as an empty struct
	opts := servers.ListOpts{}

	// Retrieve a pager (i.e. a paginated collection)
	pager := servers.List(c.client, opts)

	allservers := []servers.Server{}
	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		for _, s := range serverList {
			// "s" will be a servers.Server
			if f == nil || f(s) {
				// no selection criteria or passes criteria
				allservers = append(allservers, s)
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return allservers, nil
}
