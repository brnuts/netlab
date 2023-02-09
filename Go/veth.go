package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Linux veth command to create a veth peer
func vethAddPeer(peerA, peerB string) string {
	return fmt.Sprintf(
		"sudo ip link add %s type veth peer name %s",
		peerA,
		peerB,
	)
}

// Function created for future use
func vethDelPeer(peerA, peerB string) string {
	return fmt.Sprintf(
		"sudo ip link del %s type veth peer name %s",
		peerA,
		peerB,
	)
}

// Linux veth command to move a veth to a particular network namespace
func vethSetNameSpace(interfaceName string, nameSpace int) string {
	return fmt.Sprintf(
		"sudo ip link set %s netns %d",
		interfaceName,
		nameSpace,
	)
}

// Linux command to bring the veth interface up
func vethInterfaceUp(deviceName, interfaceName string) string {
	if deviceName == "host" {
		return fmt.Sprintf("sudo ip link set %s up", interfaceName)
	} else {
		return fmt.Sprintf(
			"docker exec %s ip link set %s up",
			deviceName,
			interfaceName,
		)
	}
}

// Function created for future use
func vethInterfaceDel(deviceName, interfaceName string) string {
	if deviceName == "host" {
		return fmt.Sprintf("sudo ip link del %s up", interfaceName)
	} else {
		return fmt.Sprintf(
			"docker exec %s ip link del %s",
			deviceName,
			interfaceName,
		)
	}
}

// Append veth information to Veths list from a single link defition
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

// Append veth information to Veth list from a single backbone link definition
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

// Populate the conf.Veths list with all veths, by reading NS and link defintion
func (conf *ConfType) loadVeths() error {
	conf.DeviceToNS = make(DeviceToNSType)

	for _, device := range conf.Topology.Devices {
		cmd := fmt.Sprintf(
			"docker inspect -f '{{.State.Pid}}' %s",
			device.Name,
		)
		out, err := runCommandOut(conf.Client, cmd)
		if err != nil {
			return err
		}
		trimOut := strings.TrimSuffix(out, "\n")
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

// Workflow that creates a veth peer
func createPeerVeth(conf *ConfType, v VethPeerType) error {
	// Add veth peer
	cmd := vethAddPeer(v.DeviceA.InterfaceName, v.DeviceB.InterfaceName)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// Move interface to namespace on DeviceA
	cmd = vethSetNameSpace(v.DeviceA.InterfaceName, v.DeviceA.NameSpace)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// If greater than 0, move to namespace on DeviceB
	// Otherwise just ignore, as the default namespace is on host
	if v.DeviceB.NameSpace > 0 {
		cmd = vethSetNameSpace(v.DeviceB.InterfaceName, v.DeviceB.NameSpace)
		if err := runCommand(conf.Client, cmd); err != nil {
			return err
		}
	}
	// Bringing Device A interfaces UP
	cmd = vethInterfaceUp(v.DeviceA.Device, v.DeviceA.InterfaceName)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// Bringing Device B interface UP
	cmd = vethInterfaceUp(v.DeviceB.Device, v.DeviceB.InterfaceName)
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}

	return nil

}

// Main loop to create veths
func (conf *ConfType) createVeths() error {
	for _, veth := range conf.Veths {
		if err := createPeerVeth(conf, veth); err != nil {
			return err
		}
	}
	return nil
}

// Add all backbone veths to the software bridge called backbone
func (conf *ConfType) addVethsToBackbone() error {
	// create bridge name backbone
	cmd := "sudo ip link add name backbone type bridge"
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	cmd = "sudo ip link set backbone up"
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	// make sure iptables don't block traffic on the bridge
	cmd = "sudo iptables -A FORWARD -j ACCEPT"
	if err := runCommand(conf.Client, cmd); err != nil {
		return err
	}
	for _, veth := range conf.Veths {
		if veth.DeviceB.NameSpace == 0 {
			cmd := fmt.Sprintf(
				"sudo ip link set %s master backbone",
				veth.DeviceB.InterfaceName,
			)
			if err := runCommand(conf.Client, cmd); err != nil {
				return err
			}
		}
	}

	return nil

}
