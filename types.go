package main

import "github.com/melbahja/goph"

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
