#!/bin/sh

# Could use GOPATH=`pwd` but I think it's a security issue to allow any path to be GOPATH
# Safer to manually hardcode this to something
export GOPATH=/home/gekko/public_html/uqm_elo
echo $GOPATH
./webserver
