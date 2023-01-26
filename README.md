# Welcome to network lab
This is a network lab that emulates several routers

## Download
https://www.dropbox.com/s/g13r1if3cf3euhl/netlab.img.gz?dl=0

## Before running it

### For MAC hosts
* On Mac you will need to add a loopback interface alias to 10.0.4.1 before running

`sudo ifconfig lo0 alias 10.0.4.1`

## Running it

### For MAC hosts
You will need some CPU parameter to run faster on MAC, like CPUID on and Accel HVF.

`sudo qemu-system-x86_64 -hda netlab.img -smp 4 -m 2G -cpu host,vmware-cpuid-freq=on -accel hvf -net user,hostfwd=tcp:10.0.4.1:22-:22 -net nic`

### For Windows hosts


### For Linux hosts
