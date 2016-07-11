package metadata

import (
	"github.com/rancher/go-rancher-metadata/metadata"
	"log"
	"time"
)

const (
	metadataUrl = "http://rancher-metadata/latest"
	empty       = ""
)

type IPFinderFromMetadata struct {
}

func NewIPFinderFromMetadata() *IPFinderFromMetadata {
	return &IPFinderFromMetadata{}
}

func (ipf *IPFinderFromMetadata) GetIP(cid string) string {

	m := metadata.NewClient(metadataUrl)

	for i := 0; i < 600; i++ {
		containers, err := m.GetContainers()
		if err != nil {
			log.Println("rancher-cni-ipam: Error getting metadata containers: %v", err)
			return empty
		}

		for _, container := range containers {
			if container.ExternalId == cid {
				log.Println("rancher-cni-ipam: got ip: %v", container.PrimaryIp)
				return container.PrimaryIp
			}
		}
		log.Println("Waiting to find IP for container: %s", cid)
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("ip not found for cid: %v", cid)
	return empty
}
