package main

import (
	clientset "github.com/akaimo/sample-controller/pkg/client/clientset/versioned"
	informers "github.com/akaimo/sample-controller/pkg/client/informers/externalversions/samplecontroller/v1"
	jobinformers "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/kubernetes"
)

type Controller struct{}

func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	jobInformer jobinformers.JobInformer,
	sampleinformer informers.SampleResourceInformer) *Controller {
	return &Controller{}
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	return nil
}
