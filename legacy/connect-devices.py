import yaml
from jumpssh import SSHSession
from typing import NamedTuple


class VethType(NamedTuple):
    router: str
    interface_name: str
    namespace: str

class VethPeerType(NamedTuple):
    veth_end_1: VethType
    veth_end_2: VethType

def read_yaml_topology(file_name: str) -> dict:
    with open(file_name, 'r') as file:
        return yaml.safe_load(file)

def get_devices(yaml_conf: dict) -> dict:
    routers = {}
    for device in yaml_conf["devices"]:
        routers[device["name"]] = device["type"]
    return routers

def get_links(yaml_conf: dict) -> dict:
    links = {}
    for link in yaml_conf["links"]:
        links[tuple(link["connection"])] = link["name"]
    return links

def connect_to_host(host: str, user: str, passwd: str) -> SSHSession:
    return SSHSession(host, user, password=passwd).open()

def get_namespaces(session: SSHSession, routers: dict) -> dict:
    namespaces = {}
    for router in routers:
        namespaces[router] = session.get_cmd_output("docker inspect -f '{{.State.Pid}}' "+router)
    return namespaces

def setup_veth(session: SSHSession, veth_peer: VethPeerType):
    


def connect_links(session: SSHSession, ns: dict, routers: dict, links: dict):


def main():
    yaml_conf = read_yaml_topology("topology.yaml")  
    routers = get_devices(yaml_conf)
    links = get_links(yaml_conf)
    host_session = connect_to_host("localhost", "netlab", "netlab")
    namespaces = get_namespaces(host_session, routers)
    connect_links(host_session, namespaces, routers, links)
    print(namespaces)
    #host_session ) SSHSession('localhost', 'netlab', password='netlab').open()
    #print(host_session.get_cmd_output('ls -lta'))
    host_session.close()

if __name__ == "__main__":
   main()