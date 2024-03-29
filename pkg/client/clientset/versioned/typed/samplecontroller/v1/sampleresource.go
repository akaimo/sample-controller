// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/akaimo/sample-controller/pkg/apis/samplecontroller/v1"
	scheme "github.com/akaimo/sample-controller/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SampleResourcesGetter has a method to return a SampleResourceInterface.
// A group's client should implement this interface.
type SampleResourcesGetter interface {
	SampleResources(namespace string) SampleResourceInterface
}

// SampleResourceInterface has methods to work with SampleResource resources.
type SampleResourceInterface interface {
	Create(*v1.SampleResource) (*v1.SampleResource, error)
	Update(*v1.SampleResource) (*v1.SampleResource, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.SampleResource, error)
	List(opts metav1.ListOptions) (*v1.SampleResourceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.SampleResource, err error)
	SampleResourceExpansion
}

// sampleResources implements SampleResourceInterface
type sampleResources struct {
	client rest.Interface
	ns     string
}

// newSampleResources returns a SampleResources
func newSampleResources(c *SamplecontrollerV1Client, namespace string) *sampleResources {
	return &sampleResources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sampleResource, and returns the corresponding sampleResource object, and an error if there is any.
func (c *sampleResources) Get(name string, options metav1.GetOptions) (result *v1.SampleResource, err error) {
	result = &v1.SampleResource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sampleresources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SampleResources that match those selectors.
func (c *sampleResources) List(opts metav1.ListOptions) (result *v1.SampleResourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.SampleResourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sampleresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sampleResources.
func (c *sampleResources) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sampleresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a sampleResource and creates it.  Returns the server's representation of the sampleResource, and an error, if there is any.
func (c *sampleResources) Create(sampleResource *v1.SampleResource) (result *v1.SampleResource, err error) {
	result = &v1.SampleResource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sampleresources").
		Body(sampleResource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sampleResource and updates it. Returns the server's representation of the sampleResource, and an error, if there is any.
func (c *sampleResources) Update(sampleResource *v1.SampleResource) (result *v1.SampleResource, err error) {
	result = &v1.SampleResource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sampleresources").
		Name(sampleResource.Name).
		Body(sampleResource).
		Do().
		Into(result)
	return
}

// Delete takes name of the sampleResource and deletes it. Returns an error if one occurs.
func (c *sampleResources) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sampleresources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sampleResources) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sampleresources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched sampleResource.
func (c *sampleResources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.SampleResource, err error) {
	result = &v1.SampleResource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sampleresources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
