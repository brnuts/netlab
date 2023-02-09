import json
from argparse import ArgumentParser

import paramiko


class NetLab:
    def __init__(self):
        self.bastion_addr = None
        self.bastion_port = None
        self.transport = None
        self.bastion = None
        self.router = None

    def connectBastion(self, addr="localhost", port=22, user="netlab", passwd="netlab"):
        pm = paramiko.SSHClient()
        pm.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        pm.connect(addr, username=user, password=passwd)
        self.transport = pm.get_transport()
        self.bastion_addr = addr
        self.bastion_port = port
        self.bastion = pm

    def connectRouter(self, router, port=22, user="netlab", passwd="netlab"):
        target_socket = (router, port)
        source_socket = (self.bastion_addr, self.bastion_port)
        router_channel = self.transport.open_channel(
            "direct-tcpip", target_socket, source_socket
        )
        router = paramiko.SSHClient()
        router.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        router.connect(router, username=user, password=passwd, sock=router_channel)
        self.router = router
        return router

    def close(self):
        self.router.close()
        self.bastion.close()


def get_neighbours(args):
    netlab = NetLab()
    netlab.connectBastion(addr=args.bastion, user=args.b_user, passwd=args.b_passwd)
    router = netlab.connectRouter(args.router, user=args.r_user, passwd=args.r_passwd)
    stdin, stdout, stderr = router.exec_command("sudo lldpctl -f json")
    data = json.loads(stdout.read().decode("ascii"))

    interfaces = {}
    for item in data["lldp"]["interface"]:
        key = list(item.keys())[0]
        values = list(item.values())[0]
        neighbour = list(values["chassis"].keys())[0]
        interfaces[key] = neighbour

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
        "--targetrouter",
        dest="router",
        required=True,
        help="Target router to find interface neighbours",
    )
    parser.add_argument(
        "-r",
        "--ruser",
        dest="r_user",
        default="netlab",
        help="Username for the router access",
    )
    parser.add_argument(
        "-s",
        "--rpass",
        dest="r_passwd",
        default="netlab",
        help="Password for the router access",
    )

    return parser.parse_args()


def main():
    args = parse_arguments()

    print(get_neighbours(args))


if __name__ == "__main__":
    main()
