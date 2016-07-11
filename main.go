package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	//"github.com/rancher/rancher-cni-ipam/fake_allocator"
	//"github.com/rancher/rancher-cni-ipam/fake_allocator/backend/disk"
	"github.com/rancher/rancher-cni-ipam/ip_finder"
	"github.com/rancher/rancher-cni-ipam/ip_finder/metadata"
)

//const logFile = "/tmp/rancher-cni-ipam.log"
const logFile = "/tmp/rancher-cni.log"

func cmdAdd(args *skel.CmdArgs) error {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("rancher-cni-ipam: cmdAdd: invoked")
	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("args: %#v", args))

	ipamConf, err := LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("ipamConf: %#v", ipamConf))

	var ipf ip_finder.IPFinder = metadata.NewIPFinderFromMetadata()
	ip_string := ipf.GetIP(args.ContainerID)

	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("ip: %#v", ip_string))

	ip, ipnet, err := net.ParseCIDR(ip_string + "/16")
	if err != nil {
		return err
	}

	// TODO: if ip is NULL, return err
	r := &types.Result{
		IP4: &types.IPConfig{
			IP: net.IPNet{IP: ip, Mask: ipnet.Mask},
		},
	}

	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("r: %#v", r))
	return r.Print()
}

func cmdDel(args *skel.CmdArgs) error {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("rancher-cni-ipam: cmdDel: invoked")
	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("args: %#v", args))

	ipamConf, err := LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("ipamConf: %#v", ipamConf))

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel)
}
