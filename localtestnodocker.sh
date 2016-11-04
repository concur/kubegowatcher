#!/bin/bash
set -e

./setupGodeps.sh

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


#proxy to the kubernetes api on the loopback alias
kubectl config use-context minikube

if ps -eaf | grep -v grep | grep 172.16.123.1
then 
   echo "kubectl proxy running";
else
  kubectl proxy --address 172.16.123.1 --disable-filter=true &
  sleep 10
fi

sudo mkdir -p /var/run/secrets/kubernetes.io/serviceaccount
sudo cp -R ./serviceaccount/* /var/run/secrets/kubernetes.io/serviceaccount/

export KUBERNETES_SERVICE_HOST=172.16.123.1
export KUBERNETES_SERVICE_PORT=8001 

go run ./package.go


