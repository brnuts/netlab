#!/bin/bash


routers="core-a1 core-a2 core-a3 core-b1 core-b2 core-c1 core-i1 core-i2 core-i3 acc-a acc-b acc-c border cpe-a cpe-b cpe-c pc-a pc-b pc-c internet"

for x in $routers
do
  echo $x
  docker stop $x
  docker rm $x
done
