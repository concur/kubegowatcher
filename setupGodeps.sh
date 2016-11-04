#!/bin/bash

# if you need to install go1.7
#sudo rm -R /usr/local/go
#
#wget https://storage.googleapis.com/golang/go1.7.darwin-amd64.tar.gz
#
#sudo tar -xvf go1.7.darwin-amd64.tar.gz
#sudo mv go /usr/local

echo $GOPATH

export GOBIN=$HOME/.go_workspace/bin
echo GOBIN=$GOBIN
go get ./package.go

export PATH=$PATH:$GOPATH/bin
