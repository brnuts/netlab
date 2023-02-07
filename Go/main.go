package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"netlab/veth"

	"github.com/melbahja/goph"
	"gopkg.in/yaml.v3"
)

type DeviceTopologyType struct {
	Name  string
	Type  string
	Image string
}

type LinkTopologyType struct {
	Name       []string
	Connection []string
}

type TopologyConfType struct {
	Devices []DeviceTopologyType
	Links   []LinkTopologyType
}

type DeviceToNSType map[string]int

type VethEndType struct {
	Device        string
	NameSpace     int
	InterfaceName string
}

type VethPeerType struct {
	DeviceA VethEndType
	DeviceB VethEndType
}

type ConfType struct {
	Topology   TopologyConfType
	DeviceToNS DeviceToNSType
	Veths      []VethPeerType
	Client     *goph.Client
}

func (conf *ConfType) readTopologyFile(fileName string) error {
	yfile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yfile, &conf.Topology); err != nil {
		return err
	}

	return nil

}

func (conf *ConfType) connectToHost(host, user, passwd string) error {

	client, err := goph.New(user, host, goph.Password(passwd))
	if err != nil {
		return err
	}
	conf.Client = client
	return nil

}

func appendVeth(conf *ConfType, link LinkTopologyType) error {
	DeviceA := link.Connection[0]
	DeviceB := link.Connection[1]
	// Back-to-back veths
	var veth VethPeerType
	veth.DeviceA.Device = DeviceA
	// Interface name DeviceA-DeviceB
	veth.DeviceA.InterfaceName = link.Name[0] + "-" + link.Name[1]
	ns, ok := conf.DeviceToNS[DeviceA]
	if !ok {
		return fmt.Errorf("could not find NS for device %s", DeviceA)
	}
	veth.DeviceA.NameSpace = ns
	veth.DeviceB.Device = DeviceB
	// Interface name DeviceB-DeviceA
	veth.DeviceB.InterfaceName = link.Name[1] + "-" + link.Name[0]
	ns, ok = conf.DeviceToNS[DeviceB]
	if !ok {
		return fmt.Errorf("could not find NS for device %s", DeviceB)
	}
	veth.DeviceB.NameSpace = ns

	conf.Veths = append(conf.Veths, veth)

	return nil
}

func appendBackboneVeth(conf *ConfType, link LinkTopologyType) error {

	var veth VethPeerType
	for _, device := range link.Connection {
		veth.DeviceA.Device = device
		veth.DeviceA.InterfaceName = device + "-" + link.Name[0]
		ns, ok := conf.DeviceToNS[device]
		if !ok {
			return fmt.Errorf("could not find NS for device %s", device)
		}
		veth.DeviceA.NameSpace = ns

		veth.DeviceB.Device = "host"
		veth.DeviceB.InterfaceName = link.Name[0] + "-" + device
		veth.DeviceB.NameSpace = 0

		conf.Veths = append(conf.Veths, veth)
	}

	return nil

}

func (conf *ConfType) loadVeths() error {
	conf.DeviceToNS = make(DeviceToNSType)

	for _, device := range conf.Topology.Devices {
		out, err := conf.Client.Run(
			"docker inspect -f '{{.State.Pid}}' " + device.Name)
		if err != nil {
			return err
		}
		trimOut := strings.TrimSuffix(string(out), "\n")
		ns, err := strconv.Atoi(trimOut)
		if err != nil {
			return err
		}
		conf.DeviceToNS[device.Name] = ns
	}

	for _, link := range conf.Topology.Links {
		// If greater than 2 is the backbone veths
		if len(link.Connection) == 2 {
			err := appendVeth(conf, link)
			if err != nil {
				return err
			}
		} else if len(link.Connection) > 2 {
			// This is the backbone veths
			err := appendBackboneVeth(conf, link)
			if err != nil {
				return err
			}

		} else {
			return fmt.Errorf("link with unexpected size: %s", link.Connection)
		}

	}
	return nil
}

func runCommand(client *goph.Client, cmd string) error {
	out, err := client.Run(cmd)
	if err != nil {
		trimOut := strings.TrimSuffix(string(out), "\n")
		return fmt.Errorf(
			"failed to run '%s', output: %s ,error: %v", cmd, trimOut, err,
		)
	}
	//fmt.Println(cmd)
	return nil
}

func createPeerVeth(conf *ConfType, v VethPeerType) error {
	// Add veth peer
	cmd := veth.AddPeer(v.DeviceA.InterfaceName, v.DeviceB.InterfaceName)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// Move interface to namespace on DeviceA
	cmd = veth.SetNameSpace(v.DeviceA.InterfaceName, v.DeviceA.NameSpace)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// If greater than 0, move to namespace on DeviceB
	// Otherwise just ignore, as the default namespace is on host
	if veth.DeviceB.NameSpace > 0 {
		cmd = veth.SetNameSpace(v.DeviceB.InterfaceName, v.DeviceB.NameSpace)
		if err := runCommand(conf.Client, cmd); err != nil {
			return err
		}
	}
	// Bringing Device A interfaces UP
	cmd = veth.InterfaceUp(v.DeviceA.Device, v.DeviceA.InterfaceName)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// Bringing Device B interface UP
	cmd = veth.InterfaceUp(v.DeviceB.Device, v.DeviceB.InterfaceName)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}

	return nil

}

func (conf *ConfType) createVeths() error {
	for _, veth := range conf.Veths {
		if err := createPeerVeth(conf, veth); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	var conf = new(ConfType)

	if err := conf.readTopologyFile("topology.yaml"); err != nil {
		log.Fatalf("read topology failed: %v", err)
	}
	//fmt.Println(conf.Topology)

	if err := conf.connectToHost("localhost", "netlab", "netlab"); err != nil {
		log.Fatalf("failed connect to host: %v", err)
	}
	defer conf.Client.Close()

	if err := conf.loadVeths(); err != nil {
		log.Fatalf("failed to load veths: %v", err)
	}

	if err := conf.createVeths(); err != nil {
		log.Fatalf("failed to create veths: %v", err)
	}

}
