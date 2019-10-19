// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/akaimo/sample-controller/pkg/apis/samplecontroller/v1"
	"github.com/akaimo/sample-controller/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type ExampleV1Interface interface {
	RESTClient() rest.Interface
	SampleResourcesGetter
}

// ExampleV1Client is used to interact with features provided by the example.com group.
type ExampleV1Client struct {
	restClient rest.Interface
}

func (c *ExampleV1Client) SampleResources(namespace string) SampleResourceInterface {
	return newSampleResources(c, namespace)
}

// NewForConfig creates a new ExampleV1Client for the given config.
func NewForConfig(c *rest.Config) (*ExampleV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ExampleV1Client{client}, nil
}

// NewForConfigOrDie creates a new ExampleV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ExampleV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ExampleV1Client for the given RESTClient.
func New(c rest.Interface) *ExampleV1Client {
	return &ExampleV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ExampleV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
