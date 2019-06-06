# Archived, deprecated, no longer maintained

# kubegowatcher
Kubernetes plugin template that watches for changes to services, pods & nodes. Business logic can then be added for ADDED, MODIFIED or DELETED events.

## Description
When this plugin is running on a kubernetes cluster it connects to the kubernetes API watch endpoints for service and node changes.

![gowatcher](https://cloud.githubusercontent.com/assets/3026995/20087524/18bb7bae-a52e-11e6-9a42-bf389468d67e.png)

# Development
## Prerequisites
* OSX
* golang & a valid $GOPATH - https://golang.org/doc/install
* Docker Native for mac - https://docs.docker.com/docker-for-mac/
* minikube - https://github.com/kubernetes/minikube/releases/ (requires kubectl & virtualbox)
* Add $DIR/serviceaccount to the list of file shares in docker -> preferences -> file shares

## Workflow
* Fork this repo
* Test your changes locally using ./localtest.sh
* For faster iterations while coding use ./localtestnodocker.sh
* Create a pull request to contribute your improvements
* No feature branches are required at this time but you may choose to do this
* Be sure to sync with upstream on a regular basis https://help.github.com/articles/syncing-a-fork/
