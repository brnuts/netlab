# Welcome to network lab
This is a network lab that emulates several routers

## Download
https://www.dropbox.com/s/fkmh1vxy5vy1fus/netlab.img.gz?dl=0

## Running it

* On Mac you will need to add a loopback interface alias to 10.0.4.1 before running
`sudo ifconfig lo0 alias 10.0.4.1`
* Start with port forwarding
`sudo qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic`
* Start with SSH and SNMP port forwarding
`sudo qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -net user,hostfwd=tcp:10.0.4.1:22-:22,hostfwd=udp:10.0.4.1:161-:161 -net nic`
