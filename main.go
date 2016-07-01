package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/rancher/rancher-cni-ipam/fake_allocator"
	"github.com/rancher/rancher-cni-ipam/fake_allocator/backend/disk"
)

//const logFile = "/tmp/rancher-cni-ipam.log"
const logFile = "/tmp/rancher-cni.log"

func getFakeIpamConfig() *fake_allocator.IPAMConfig {

	_, ipn, _ := net.ParseCIDR("10.42.1.0/24")
	subnet := types.IPNet(*ipn)

	c := &fake_allocator.IPAMConfig{Name: "rancher-network-abcd",
		Type:   "rancher-cni-ipam",
		Subnet: subnet,
	}

	return c
}

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

	fakeIpamConfig := getFakeIpamConfig()

	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("fakeIpamConfig: %#v", fakeIpamConfig))

	store, err := disk.New(fakeIpamConfig.Name)
	if err != nil {
		log.Println("rancher-cni-ipam: couldn't create disk")
		return err
	}
	defer store.Close()

	allocator, err := fake_allocator.NewIPAllocator(fakeIpamConfig, store)
	if err != nil {
		log.Println("rancher-cni-ipam: couldn't get fake_allocator")
		return err
	}

	fakeIpInfo, err := allocator.Get(args.ContainerID)
	if err != nil {
		log.Println("rancher-cni-ipam: Error getting fakeIpInfo")
		return err
	}

	log.Println("rancher-cni-ipam: %s", fmt.Sprintf("fakeIpInfo: %#v", fakeIpInfo))

	r := &types.Result{
		IP4: fakeIpInfo,
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

	// TODO: Figure out later what to return

	fakeIpamConfig := getFakeIpamConfig()

	log.Println("%s", fmt.Sprintf("ipamConf: %#v", fakeIpamConfig))

	store, err := disk.New(fakeIpamConfig.Name)
	if err != nil {
		return err
	}
	defer store.Close()

	allocator, err := fake_allocator.NewIPAllocator(fakeIpamConfig, store)
	if err != nil {
		return err
	}

	return allocator.Release(args.ContainerID)
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel)
}
