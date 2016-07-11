package cattle

import (
	"github.com/hashicorp/golang-lru"
	rancher "github.com/rancher/go-rancher/client"
	"log"
)

type RancherIPFinder struct {
	Rancher *rancher.RancherClient
	// TODO: Is this even needed coz if it's running as an executable
	//		 the state is anyways lost. Does it make sense to put it
	//		 in disk store?
	cache *lru.Cache
}

func NewRancherIPFinder(clientOpts *rancher.ClientOpts) (*RancherIPFinder, error) {
	c, err := lru.New(256)
	if err != nil {
		return nil, err
	}

	rancherClient, err := rancher.NewRancherClient(clientOpts)
	if err != nil {
		return nil, err
	}

	return &RancherIPFinder{
		Rancher: rancherClient,
		cache:   c,
	}, nil
}

func (r *RancherIPFinder) GetIp(containerId string) (string, error) {
	if val, ok := r.cache.Get(containerId); ok {
		if ip, ok := val.(string); ok {
			return ip, nil
		}
	}

	containers, err := r.Rancher.Container.List(&rancher.ListOpts{
		Filters: map[string]interface{}{
			"externalId":   containerId,
			"removed_null": "",
		},
	})
	if err != nil {
		return "", err
	}

	if len(containers.Data) == 0 {
		return "", nil
	}

	rancherContainer := containers.Data[0]

	ipAddr := ""

	// If the hostNetwork for the pod is set to true
	if rancherContainer.NetworkMode == "host" {
		hosts := &rancher.HostCollection{}
		err := r.Rancher.GetLink(rancherContainer.Resource, "hosts", hosts)
		if err != nil {
			return "", err
		}
		if len(hosts.Data) == 0 {
			return "", nil
		}
		host := hosts.Data[0]
		ipAddresses := &rancher.IpAddressCollection{}
		err = r.Rancher.GetLink(host.Resource, "ipAddresses", ipAddresses)
		if err != nil {
			return "", err
		}
		if len(ipAddresses.Data) == 0 {
			return "", nil
		}
		ipAddr = ipAddresses.Data[0].Address
	} else if rancherContainer.PrimaryIpAddress != "" {
		ipAddr = rancherContainer.PrimaryIpAddress
	}

	if ipAddr != "" {
		log.Printf("Found IP %s for container %s", ipAddr, containerId)
		r.cache.Add(containerId, ipAddr)
	}

	return ipAddr, nil
}
