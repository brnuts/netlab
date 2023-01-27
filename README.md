# Welcome to network lab
This is a network lab that emulates several routers

## Download
https://www.dropbox.com/s/g4ak2bc26g0tbfr/netlab.img.gz?dl=0

## Before running it

### Add IP to the loopback
* On Mac you will need to add a loopback interface alias to 10.0.4.1 before running

`sudo ifconfig lo0 alias 10.0.4.1`

* On Linux you will need to use the following command

`sudo ip addr add 10.0.4.1 dev lo`


## Running it

### For MAC hosts
You will need some CPU parameter to run faster on MAC, like CPUID on and Accel HVF. MacOS may give low priority to `qemu` program as it is running from a terminal, to give more priority use `nice -20` in front.

`sudo nice -20 qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -cpu host,vmware-cpuid-freq=on -accel hvf -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic`

### For Windows hosts


### For Linux hosts
