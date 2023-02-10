# Welcome to network lab
This page describe how to launch a host and start routers for the network lab.

# How to launch the host


## Using pre-build image on VirtualBox
- Download image at
```
https://www.dropbox.com/s/2l788jzbpe02rnt/netlab.vdi.bz3?dl=0
```
* Uncompress `netlab.vdi.bz3`
* Create new VM on VirtualBox using the pre-built image as the hard disk `netlab.vdi`
* Add port forwarding to 22 to access SSH by going to Virtualbox settings for the VM and configure NAT network with port forwarding 22, so you can access the network emulation via SSH. With that you can access the network lab by doing `ssh netlab@localhost`

<img src="https://github.com/brnuts/netlab/blob/main/Port-foward-example-Virtualbox.png" width="300"/>

* All devices including the host have the username `netlab` with password `netlab`.
* `sudo` is configured to be used by username `netlab` without password.
* Access to `vtysh` on the routers via username `vtysh` with password `vtysh`.

* (Optional) Install Guest Additions on the Linux VM, use these steps:
   * Open VirtualBox.
   * Right-click the virtual machine, select the Start submenu and choose the Normal Start option.
   * Sign in to netlab (username: `netlab` password: `netlab`)
   * Click the Devices menu and select the Insert Guest Additions CD image option.
   * On Linux do `sudo mount /dev/cdrom`
   * Run the command `sudo sh /dev/cdrom0/VBoxLinuxAdditions.run`
   * Reboot your VM


## Using pre-build image on Qemu
### Download image at
```
https://www.dropbox.com/s/ajt90d4vfvkeb95/netlab.img.bz3?dl=0
```
### Before running it
#### Add IP to the loopback
* On Mac you will need to add a loopback interface alias to 10.0.4.1 before running
```
sudo ifconfig lo0 alias 10.0.4.1
```
* On Linux you will need to use the following command
```
sudo ip addr add 10.0.4.1 dev lo
```

### Running it

#### For MAC hosts
##### startup
You will need some CPU parameter to run faster on MAC, like CPUID on and Accel HVF. MacOS may give low priority to `qemu` program as it is running from a terminal, to give more priority use `nice -20` in front.
```
sudo nice -20 qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -cpu host,vmware-cpuid-freq=on -accel hvf -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic
```

#### For Linux hosts
##### Normal basic startup
```
sudo qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic
```

## Doing yourself
- Install any Linux distribution
- Install docker (https://docs.docker.com/engine/install/)
- Install lldpd
```
apt install lldpd
```
- Create user `netlab` locally, suggest using netlab as password:
```
useradd netlab
```
- Allow `netlab` to use docker, adding to docker group:
```
sudo usermod -aG docker netlab
```
- Add the scripts on the `scripts` directory to netlab home directory:
```
cp scripts/* /home/netlab
```
- Install python3:
```
apt install python3
```
- Install pip:
```
apt install pip
```
- configure update-hosts.service by copying update-hosts scripts to `/etc/systemd/system`:
```
cp update-hosts.* /etc/systemd/system/
```
- enable update-hosts: 
```
systemctl enable update-hosts.service
```
- start update-hosts:
```
systemctl start update-hosts`
```
- run scrip to start all containers:
```
./start-containers.sh
```

Done!

