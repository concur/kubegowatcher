package main

import (
  "time"
  "log"
  "flag"
  api "k8s.io/client-go/pkg/api/v1"
  "encoding/json"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
)

// Client holds state.
type Client struct {
  //kubeClient to use inside functions
  kubeClient *kubernetes.Clientset
}


func (c *Client) newService(obj interface{}) {
  
  var s = api.Service{}
  
  b, err := json.Marshal(obj)
  log.Printf("Object: %+v", s)
  
  if err = json.Unmarshal(b, &s); err == nil {
    
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
  log.Printf("Object: %+v", s)
  
  if err = json.Unmarshal(b, &s); err == nil {
    //skip this node if it's not scheduleable
    if s.Spec.Unschedulable {return}
    // add code here to handle new Nodes
    // s.ObjectMeta.Name
  }
}

func removeNode(obj interface{}) {
  var s = api.Node{}
  
  b, err := json.Marshal(obj)
  log.Printf("Object: %+v", s)
  
  if err = json.Unmarshal(b, &s); err == nil {
    //add code here to handle remove node event
    //s.ObjectMeta.Name
  }
}

func removeService(obj interface{}) {
  var s = api.Service{}
  
  b, err := json.Marshal(obj)
  err = json.Unmarshal(b, &s)
  
  if err = json.Unmarshal(b, &s); err == nil {
    //add business logic here for service removed
    log.Printf("Object: %+v", s)
  }
}

func updateService(obj interface{}) {
  var sNew = api.Service{}
  
  b, err := json.Marshal(obj)
  
  if err = json.Unmarshal(b, &sNew); err == nil {
    log.Printf("Object: %+v", sNew)
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
    log.Printf("Using kubernetes API %v", config.GroupVersion)
  
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

func (c *Client) watchForServices(kubeClient *kubernetes.Clientset, timeout int64) error {
  //timeout := int64(30)
  watchServicesInterface, err := kubeClient.Core().Services("").Watch(api.ListOptions{Watch: true, TimeoutSeconds: &timeout})
  if err != nil {
    log.Printf("Error retrieving watch interface for services: %+v", err)
    panic(err.Error())    
  }
  
  events := watchServicesInterface.ResultChan()
  for {
    event, ok := <-events
    log.Printf("Service Event %v: %+v", event.Type, event.Object)
    log.Printf("Service Event Type %v", event.Type)
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

func watchForNodes(kubeClient *kubernetes.Clientset, timeout int64) {
  
  watchNodesInterface, err := kubeClient.Core().Nodes().Watch(api.ListOptions{Watch: true, TimeoutSeconds: &timeout})
  if err != nil {
    log.Printf("Error retrieving watch interface for nodes: %+v", err)
    panic(err.Error())    
  }
  
  events := watchNodesInterface.ResultChan()
  for {
    event, ok := <-events
    log.Printf("Node Event %v: %+v", event.Type, event.Object)
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


func main() {
  
  flag.Set("v", "5")
  flag.Set("logtostderr", "true")
  flag.Parse()
  
  var err error
  var c = Client{}
  watchTimeout := int64(1800) //30 minute resync
  
  kubeClient, err := newKubeClient()
  if err != nil {
    log.Printf("Failed to create a kubernetes client: %+v", err)
  } else {
    log.Printf("k8s client created.")
  }
  
  c.kubeClient = kubeClient

  for {
    go watchForNodes(kubeClient, watchTimeout)
    go c.watchForServices(kubeClient, watchTimeout)
    time.Sleep(time.Second * time.Duration(watchTimeout) + time.Second * 2)
    log.Printf("Restarting watchers from %v second timeout.", watchTimeout + 2)
  }
  
}