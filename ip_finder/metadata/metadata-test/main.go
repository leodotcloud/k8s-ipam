package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-cni-ipam/ip_finder"
	"github.com/rancher/rancher-cni-ipam/ip_finder/metadata"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Errorf("Invalid arguments, need just one argument which is container id")
		os.Exit(1)
	}
	cid := os.Args[1]

	logrus.Infof("cid: %v", cid)

	ipf := metadata.NewIPFinderFromMetadata()

	ip := ipf.GetIP(cid)

	logrus.Infof("ip: %v", ip)
}
