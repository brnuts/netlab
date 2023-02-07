package veth

import "fmt"

func AddPeer(peerA, peerB string) string {
	return fmt.Sprintf(
		"sudo ip link add %s type veth peer name %s",
		peerA,
		peerB,
	)
}

func DelPeer(peerA, peerB string) string {
	return fmt.Sprintf(
		"sudo ip link del %s type veth peer name %s",
		peerA,
		peerB,
	)
}

func SetNameSpace(interfaceName string, nameSpace int) string {
	return fmt.Sprintf(
		"sudo ip link set %s netns %d",
		interfaceName,
		nameSpace,
	)
}

func InterfaceUp(deviceName, interfaceName string) string {
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

func InterfaceDel(deviceName, interfaceName string) string {
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
