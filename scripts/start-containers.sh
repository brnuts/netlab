#!/bin/bash

set -e

routers="core-a1 core-a2 core-a3 core-b1 core-b2 core-c1 core-i1 core-i2 core-i3 acc-a acc-b acc-c border cpe-a cpe-b cpe-c"

for x in $routers
do
  echo $x
  docker volume create $x
  docker run --privileged --name $x -v "$x:/etc/" --hostname $x -d brnuts/routerlab sh -c /usr/lib/frr/docker-start
  docker update --restart unless-stopped $x
done

pcs="pc-a pc-b pc-c"

for x in $pcs
do
  echo $x
  docker volume create $x
  docker run --privileged --name $x -v "$x:/etc/" --hostname $x -d brnuts/pclab sh -c /etc/docker/docker-start
  docker update --restart unless-stopped $x
done

echo "internet"
docker volume create internet
docker run --privileged --name internet -v "internet:/etc/" --hostname internet -d brnuts/internetlab sh -c /etc/docker/docker-start
docker update --restart unless-stopped internet

