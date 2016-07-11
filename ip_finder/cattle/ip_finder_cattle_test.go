package cattle

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	rancher "github.com/rancher/go-rancher/client"
	"os"
	"testing"
)

func TestGetIP(t *testing.T) {

	containerId := os.Getenv("TEST_CONTAINER_ID")

	// Note:
	// One way to get the access and secret keys is to
	// generate them in the UI. If you have copied from
	// the agent container or somewhere else, they won't
	// work.

	clientOpts := &rancher.ClientOpts{
		Url:       os.Getenv("CATTLE_URL"),
		AccessKey: os.Getenv("CATTLE_ACCESS_KEY"),
		SecretKey: os.Getenv("CATTLE_SECRET_KEY"),
	}

	if clientOpts.Url == "" {
		t.Errorf("CATTLE_URL is not set")
		return
	}
	if clientOpts.AccessKey == "" {
		t.Errorf("CATTLE_ACCESS_KEY is not set")
		return
	}
	if clientOpts.SecretKey == "" {
		t.Errorf("CATTLE_SECRET_KEY is not set")
		return
	}

	if containerId == "" {
		t.Errorf("TEST_CONTAINER_ID is not set")
		return
	}

	r, err := NewRancherIPFinder(clientOpts)
	if err != nil {
		t.Errorf("Couldn't create: %v", err)
		return
	}

	log.Printf("%s", fmt.Sprintf("r: %#v", r))

	ip, err := r.GetIp(containerId)
	if err != nil {
		t.Errorf("Couldn't get IP: %#v", err)
		return
	}

	log.Printf("%s", fmt.Sprintf("Got IP: %#v", ip))

	return
}
