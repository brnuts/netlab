import json
from argparse import ArgumentParser

import paramiko


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


def get_neighbours(args):
    netlab = NetLab()
    netlab.connectBastion(addr=args.bastion, user=args.b_user, passwd=args.b_passwd)
    device = netlab.connectDevice(args.device, user=args.r_user, passwd=args.r_passwd)
    stdin, stdout, stderr = device.exec_command("sudo lldpctl -f json")
    data = json.loads(stdout.read().decode("ascii"))

    interfaces = []
    for item in data["lldp"]["interface"]:
        key = list(item.keys())[0]
        values = list(item.values())[0]
        neighbour = list(values["chassis"].keys())[0]
        interfaces.append({"Interface": key, "Device": neighbour})

    netlab.close()
    return interfaces


def parse_arguments():
    parser = ArgumentParser()
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
        "-t",
        "--targetdevice",
        dest="device",
        required=True,
        help="Target device to find interface neighbours",
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

    return parser.parse_args()


def main():
    args = parse_arguments()

    print(json.dumps(get_neighbours(args), indent=4))


if __name__ == "__main__":
    main()
