package main

import (
	"flag"
	"log"
)

var (
	host         string
	user         string
	passwd       string
	port         uint
	topologyFile string
)

func init() {

	flag.StringVar(&host, "host", "localhost", "Host IP for netlab")
	flag.UintVar(&port, "port", 22, "SSH port to access Host IP for netlab")
	flag.StringVar(&user, "user", "netlab", "Username to access netlab host")
	flag.StringVar(&passwd, "pw", "netlab", "Password to access netlab host")
	flag.StringVar(&topologyFile, "topo", "topology.yaml", "Topology yaml file")

}

func main() {

	flag.Parse()

	var conf = new(ConfType)

	log.Printf("reading %s file", topologyFile)
	if err := conf.readTopologyFile(topologyFile); err != nil {
		log.Fatalf("read topology failed: %v", err)
	}

	log.Printf("connecting via SSH to %s@%s", user, host)
	if err := conf.connectHost(host, port, user, passwd); err != nil {
		log.Fatalf("failed to connect to host: %v", err)
	}
	defer conf.Client.Close()

	log.Printf("loading veth information")
	if err := conf.loadVeths(); err != nil {
		log.Fatalf("failed to load veths: %v", err)
	}

	log.Printf("connecting devices with %d veths", len(conf.Veths))
	if err := conf.createVeths(); err != nil {
		log.Fatalf("failed to create veths: %v", err)
	}

	log.Printf("creating bridge and adding veth backbones")
	if err := conf.addVethsToBackbone(); err != nil {
		log.Fatalf("failed to add veths to backbone: %v", err)
	}

	log.Print("all done successfully")
}
