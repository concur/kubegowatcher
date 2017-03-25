#!/bin/bash

#Copyright 2016 Concur Technologies, Inc.
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

# if you need to install go1.7
#sudo rm -R /usr/local/go
#
#wget https://storage.googleapis.com/golang/go1.7.darwin-amd64.tar.gz
#
#sudo tar -xvf go1.7.darwin-amd64.tar.gz
#sudo mv go /usr/local

echo $GOPATH
export PATH=$PATH:$GOPATH/bin

go get -u github.com/golang/dep/...
dep ensure -update
