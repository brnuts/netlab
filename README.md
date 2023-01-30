# Welcome to network lab
This is a network lab that emulates several routers


# How to install

## Doing yourself
- Install any Linux distribution
- Install docker (https://docs.docker.com/engine/install/)
- Install lldpd
```
apt install lldpd
```
- Create user `netlab` locally:
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
- Install `yq`:
```
apt install yq
```
- copy update-hosts scripts to `/etc/systemd/system`:
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

## Using VirtualBox
### Download
https://www.dropbox.com/s/h6mfz20ocu6i4g3/netlab.vdi.bz2?dl=0

### Add port forwarding to 22 to access SSH
Go to Virtualbox settings for the VM and configure NAT network with port forwarding 22, so you can access the network emulation via SSH.

### Username and password
- all devices including the host has the username: `netlab` and password: `netlab`.

- `sudo` is configured to be used by username `netlab` without password.

## Using Qemu
### Download
https://www.dropbox.com/s/yki02bv3f09mmlh/netlab.img.bz2?dl=0

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
##### Normal basic startup
```
sudo qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic
```
##### Advanced startup
You will need some CPU parameter to run faster on MAC, like CPUID on and Accel HVF. MacOS may give low priority to `qemu` program as it is running from a terminal, to give more priority use `nice -20` in front.
```
sudo nice -20 qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -cpu host,vmware-cpuid-freq=on -accel hvf -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic
```

#### For Linux hosts
##### Normal basic startup
```
sudo qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic
```
