# host_generator
Openstack Automation from Traffic Ops

*host_generator* uses the Traffic Ops API to gather server information.   It uses that data to
create new server instances in Openstack using the [gophercloud]
(http://github.com/rackspace/gophercloud) API.  Profile/parameter data describes the
configuration needed to create properly sized and configured instances.

This will help in our effort to make Traffic Ops the "source of truth" for Traffic Ops
servers.

