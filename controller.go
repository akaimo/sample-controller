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
