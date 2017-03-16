/*
Copyright 2016 Concur Technologies, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
  "time"
  "log"
  api "k8s.io/client-go/pkg/api/v1"
  //"k8s.io/apimachinery/pkg/apis/meta/v1"
  "encoding/json"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
)

// Client can be used to hold state for business logic as needed.
type Client struct {
  //kubeClient to use inside functions
  kubeClient *kubernetes.Clientset
}


func (c *Client) newService(obj interface{}) {
  
  var s = api.Service{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
  log.Printf("New Service: %+v", s)    
// potentially useful if you delay or check for existing ExternalIPs
//    //check delay since event occurred
//    var k8sTimeFmt = "2006-01-02 15:04:05 -0700 MST"
//    t, _ := time.Parse(k8sTimeFmt, s.ObjectMeta.CreationTimestamp.String())
//    delay := time.Now().Sub(t).Seconds()
//    if delay > 60 && s.Spec.ExternalIPs != nil {
//      //bail no need to do an update
//      log.Printf("no need to check this service: %v delay: %v ExternalIPs: %v", s.ObjectMeta.Name, delay, s.Spec.ExternalIPs)
//      return
//    } else {
//      log.Printf("Must check this service: %v delay: %v ExternalIPs: %v", s.ObjectMeta.Name, delay, s.Spec.ExternalIPs)
//    }
    //svcClient is used to make updates to the service
    svcClient :=  c.kubeClient.Core().Services(s.Namespace)
    
    if s.ObjectMeta.Annotations == nil {
      s.ObjectMeta.Annotations = make(map[string]string)
    }

    if s.ObjectMeta.Annotations["Example"] == "" {
      s.ObjectMeta.Annotations["Example"] = "myInterestingAnnotation"
      log.Printf("Setting Annotation 'Example': %v", s.ObjectMeta.Name)
        if _, err := svcClient.Update(&s); err != nil {
          log.Printf("error updating service fields: %v", err)
        }
    } else {
      log.Printf("Annotation 'Example' already exists: %v Value: %v", s.ObjectMeta.Name, s.ObjectMeta.Annotations["Example"])
    }
  }
}

func newNode(obj interface{}) {

  var s = api.Node{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
    log.Printf("Node Added: %+v", s)
    //skip this node if it's not scheduleable
    if s.Spec.Unschedulable {return}
    // add code here to handle new nodes
    // s.ObjectMeta.Name
  }
}

func removeNode(obj interface{}) {

  var s = api.Node{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
    log.Printf("Node Removed: %+v", s)
    //add code here to handle removed nodes
    //s.ObjectMeta.Name
  }
}

func newPod(obj interface{}) {

  var s = api.Pod{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
    log.Printf("Pod Added: %+v", s)
    // add code here to handle new nodes
    // s.ObjectMeta.Name
  }
}

func removePod(obj interface{}) {

  var s = api.Pod{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
    log.Printf("Pod Removed: %+v", s)
    //add code here to handle removed nodes
    //s.ObjectMeta.Name
  }
}

func removeService(obj interface{}) {

  var s = api.Service{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
    log.Printf("Service Removed: %+v", s)
    //add business logic here for service removed
  }
}

func updateService(obj interface{}) {

  var s = api.Service{}
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &s); err == nil {
    log.Printf("Service Updated: %+v", s)
  }
  
  //add business logic here for service updated
//  var k8sTimeFmt = "2006-01-02 15:04:05 -0700 MST"
//  t, _ := time.Parse(k8sTimeFmt, sNew.ObjectMeta.CreationTimestamp.String())
//  delay := time.Now().Sub(t).Seconds()
//  if delay < 60 {
//    p.newService(obj)
//  }
}

func newKubeClient() (*kubernetes.Clientset, error) {
  config, err := rest.InClusterConfig()
  if err != nil {
      panic(err.Error())
  }	
  log.Printf("Using %s for kubernetes master", config.Host)
  
  if config.Host == "https://172.16.123.1:8001" {
    config.Host = "http://172.16.123.1:8001"
  }
    // creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
  
    return clientset, err
}

func (c *Client) watchForServices(timeout int64) error {
  watchServicesInterface, err := c.kubeClient.Core().Services("").Watch(api.ListOptions{Watch: true, TimeoutSeconds: &timeout})
  if err != nil {
    log.Printf("Error retrieving watch interface for services: %+v", err)
    panic(err.Error())    
  }
  
  events := watchServicesInterface.ResultChan()
  for {
    event, ok := <-events
    log.Printf("Service Event %v: %+v", event.Type, event.Object)
    if event.Type == "ADDED" {
      c.newService(event.Object)
    } else if event.Type == "MODIFIED" {
      updateService(event.Object)
    } else if event.Type == "DELETED" {
      removeService(event.Object)
    } else if event.Type == "" {
      log.Printf("Service watch timed out: %v seconds", timeout)
    } else {
      log.Printf("Service watch unhandled event: %v", event.Type)
    }
    if !ok { break }
  }
  
  return err
}

func (c *Client) watchForNodes(timeout int64) {
  watchNodesInterface, err := c.kubeClient.Core().Nodes().Watch(api.ListOptions{Watch: true, TimeoutSeconds: &timeout})
  if err != nil {
    log.Printf("Error retrieving watch interface for nodes: %+v", err)
    panic(err.Error())    
  }
  
  events := watchNodesInterface.ResultChan()
  for {
    event, ok := <-events
//    log.Printf("Node Event %v: %+v", event.Type, event.Object)
    log.Printf("Node Event Type %v", event.Type)
    if event.Type == "ADDED" {
      newNode(event.Object)
    } else if event.Type == "DELETED" {
      removeNode(event.Object)
    } else if event.Type == "" {
      log.Printf("Node watch timed out: %v seconds", timeout)
    }
    if !ok { break }
  }
  return
}

func (c *Client) watchForPods(timeout int64) {
  watchInterface, err := c.kubeClient.Core().Pods("").Watch(api.ListOptions{Watch: true, TimeoutSeconds: &timeout})
  if err != nil {
    log.Printf("Error retrieving watch interface for pods: %+v", err)
    panic(err.Error())    
  }
  
  events := watchInterface.ResultChan()
  for {
    event, ok := <-events
    log.Printf("Pod Event %v: %+v", event.Type, event.Object)
    if event.Type == "ADDED" {
      //newPod(event.Object)
    } else if event.Type == "DELETED" {
      //removePod(event.Object)
    } else if event.Type == "" {
      log.Printf("Pod watch timed out: %v seconds", timeout)
    }
    if !ok { break }
  }
  
  return
}

func main() {
  
  watchTimeout := int64(1800) //30 minute resync
  
  k8sClient, err := newKubeClient()
  if err != nil {
    log.Printf("Failed to create a kubernetes client: %+v", err)
  } else {
    log.Printf("k8s client created.")
  }
  
  //add the client
  c := Client{kubeClient: k8sClient}
  
  for {
    go c.watchForNodes(watchTimeout)
    go c.watchForServices(watchTimeout)
    go c.watchForPods(watchTimeout)
    time.Sleep(time.Second * time.Duration(watchTimeout) + time.Second * 2)
    log.Printf("Restarting watchers from %v second timeout.", watchTimeout + 2)
  }
  
}