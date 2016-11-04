# kubegowatcher
Kubernetes plugin template that checks for changes to service & nodes.

## Description
When this plugin is running on a kubernetes cluster it connects to the kubernetes API watch endpoints for service and node changes.

# Development
## Prerequisites
* OSX
* Docker Native for mac - https://docs.docker.com/docker-for-mac/
* minikube - https://github.com/kubernetes/minikube/releases/
* Add $DIR/serviceaccount to the list of file shares in docker -> preferences -> file shares

## Workflow
* Fork this repo
* Test your changes locally using ./localtest.sh
* For faster iterations while coding use ./localtestnodocker.sh
* Create a pull request when you're happy with the changes
* No feature branches are required at this time but you may choose to do this
* Be sure to sync with upstream on a regular basis https://help.github.com/articles/syncing-a-fork/
