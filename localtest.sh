#!/bin/bash
set -e

./setupGodeps.sh

#run tests
go test . -v

#build the binary
echo "About to build go binary... "
#CGO_ENABLED=0 GOOS=linux go build -o main package.go
GOOS=linux go build -o main package.go
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

#create an alias for the loopback interface so it's reachable inside the docker vm
if ifconfig | grep -q 172.16.123.1
then 
   echo "lo0 alias found";
else
   echo "creating lo0 alias, enter sudo password";
   sudo ifconfig lo0 alias 172.16.123.1
fi

# start minikube
if ps -eaf | grep -v grep | grep minikube
then 
   echo "minikube started";
else
   minikube start
   sleep 10
fi

#set the service account config locally
set +e && base64 -D <<< dGVzdAo= &>/dev/null
if [ "$?" == "0" ]; then
export basecmd="base64 -D"
else
export basecmd="base64 -d"
fi
set -e

mkdir -p $DIR/serviceaccount
kubectl get secrets -o yaml | grep "token:" | awk '{print $2}' | ${basecmd} > $DIR/serviceaccount/token
kubectl get secrets -o yaml | grep "ca.crt:" | awk '{print $2}' | ${basecmd} > $DIR/serviceaccount/ca.crt

#create the local docker image
echo "About to build docker container. If this hangs restart docker... "
docker build -t kubegowatcher .

#proxy to the kubernetes api on the loopback alias
kubectl config use-context minikube
kubectl proxy --address 172.16.123.1 --disable-filter=true &
sleep 10

# run the docker image
docker run -d -e KUBERNETES_SERVICE_HOST=172.16.123.1 -e KUBERNETES_SERVICE_PORT=8001 -v $DIR/serviceaccount:/var/run/secrets/kubernetes.io/serviceaccount kubegowatcher /main

#delete/create and expose a service
kubectl delete svc nginx 2> /dev/null || true
kubectl run nginx --image=nginx || true
kubectl expose deployment nginx --port=443 --type=LoadBalancer

#get docker output
docker logs `(docker ps -n 1 | awk 'FNR > 1 { print $1 }')` -f