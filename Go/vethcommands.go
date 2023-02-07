package main

import "fmt"

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

func vethSetNameSpace(interfaceName string, nameSpace int) string {
	return fmt.Sprintf(
		"sudo ip link set %s netns %d",
		interfaceName,
		nameSpace,
	)
}

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
