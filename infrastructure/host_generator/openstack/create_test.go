package openstack

import (
	"fmt"
	"testing"
)

var cr *Creator

func init() {
	c, err := New()
	if err != nil {
		panic(err)
	}
	cr = c
}

func TestPickFlavor(t *testing.T) {
	type req struct {
		VCPUs, RAM int
	}

	reqs := []req{{0, 4 * 1024}, {1, 24 * 1024}, {8, 16 * 1024}}
	for _, r := range reqs {
		f, err := cr.PickFlavor(r.VCPUs, r.RAM)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("pick for %+v: %+v\n", r, f)
	}
}

func TestPickImage(t *testing.T) {
	type req struct {
		name string
	}

	reqs := []req{{"foo"}, {"bar"}}
	for _, r := range reqs {
		f, err := cr.PickImage("Comcast Cloud CentOS 7.1 x86_64 v1.0")
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("pick for %+v: %+v\n", r, f)
	}

}
