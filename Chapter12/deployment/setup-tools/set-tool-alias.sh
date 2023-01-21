#!/usr/bin/env sh

root=$(realpath ${PWD}/../..)

docker build -t deploytools:latest .

alias deploytools='docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v ~/.aws:/root/.aws -v ~/.docker:/root/.docker -v ~/.kube:/root/.kube -v ${root}:/mallbots -v ${PWD}:/mallbots/deployment/.current -w /mallbots/deployment/.current deploytools'

echo "---"
echo
echo "Usage: deploytools <cmd [parameters]>"
