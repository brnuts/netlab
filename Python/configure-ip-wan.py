import ipaddress
import json
import logging
from argparse import ArgumentParser

import paramiko
import yaml

logging.basicConfig(format="%(funcName)s(): %(message)s")
logging.getLogger("paramiko").setLevel(logging.WARNING)
logger = logging.getLogger()
logger.setLevel(logging.INFO)


class NetLab:
    def __init__(self):
        self.bastion_addr = None
        self.bastion_port = None
        self.transport = None
        self.bastion = None
        self.device = None

    def connectBastion(self, addr="localhost", port=22, user="netlab", passwd="netlab"):
        pm = paramiko.SSHClient()
        pm.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        pm.connect(addr, username=user, password=passwd)
        self.transport = pm.get_transport()
        self.bastion_addr = addr
        self.bastion_port = port
        self.bastion = pm

    def connectDevice(self, device, port=22, user="netlab", passwd="netlab"):
        target_socket = (device, port)
        source_socket = (self.bastion_addr, self.bastion_port)
        device_channel = self.transport.open_channel(
            "direct-tcpip", target_socket, source_socket
        )
        device = paramiko.SSHClient()
        device.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        device.connect(device, username=user, password=passwd, sock=device_channel)
        self.device = device
        return device

    def close(self):
        self.device.close()
        self.bastion.close()


def run_device_command(args, device, command):
    netlab = NetLab()
    netlab.connectBastion(addr=args.bastion, user=args.b_user, passwd=args.b_passwd)
    logger.info("running cmd '{}' on device '{}'".format(command, device))
    device = netlab.connectDevice(device, user=args.r_user, passwd=args.r_passwd)
    stdin, stdout, stderr = device.exec_command(command)
    return stdout


def read_topology_file(file_name):
    with open(file_name) as file:
        data = yaml.safe_load(file)

    return data


def get_wan_interfaces(topology):
    wan_link = [x for x in topology["links"] if x["name"][0] == "wan"]
    wan_devices = wan_link[0]["connection"]
    wan_interfaces = {}
    for device in wan_devices:
        wan_interfaces[device] = device + "-" + "wan"
    return wan_interfaces


def configure_wan_interfaces(args, wan_interfaces):
    network = ipaddress.ip_network(args.subnet)
    valid_ips = list(network.hosts())
    prefix = network.prefixlen
    for device, interface in wan_interfaces.items():
        ip = valid_ips.pop(0)
        cmd = "sudo ip addr add {}/{} dev {}".format(ip, prefix, interface)
        run_device_command(args, device, cmd)


def parse_arguments():
    parser = ArgumentParser()
    parser.add_argument(
        "-t",
        "--topo",
        dest="topology",
        default="topology.yaml",
        help="Topology description yaml file",
    )
    parser.add_argument(
        "-b",
        "--bastion",
        dest="bastion",
        default="localhost",
        help="Bastion address to connect",
    )
    parser.add_argument(
        "-u",
        "--user",
        dest="b_user",
        default="netlab",
        help="Username for the bastion access",
    )
    parser.add_argument(
        "-p",
        "--pass",
        dest="b_passwd",
        default="netlab",
        help="Password for the bastion access",
    )
    parser.add_argument(
        "-r",
        "--ruser",
        dest="r_user",
        default="netlab",
        help="Username for the device access",
    )
    parser.add_argument(
        "-s",
        "--rpass",
        dest="r_passwd",
        default="netlab",
        help="Password for the device access",
    )
    parser.add_argument(
        "-i",
        "--ipsubnet",
        dest="subnet",
        default="10.200.200.0/24",
        help="IP subnet to be used in the WAN interfaces",
    )

    return parser.parse_args()


def main():
    args = parse_arguments()

    logger.info("reading topology file")
    topology = read_topology_file(args.topology)
    logger.info("getting wan interfaces")
    wan_interfaces = get_wan_interfaces(topology)
    logger.info("configuring inteface wan devices with {} subnet".format(args.subnet))
    configure_wan_interfaces(args, wan_interfaces)


if __name__ == "__main__":
    main()
