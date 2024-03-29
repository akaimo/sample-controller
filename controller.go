package main

import (
	"fmt"
	clientset "github.com/akaimo/sample-controller/pkg/client/clientset/versioned"
	samplescheme "github.com/akaimo/sample-controller/pkg/client/clientset/versioned/scheme"
	informers "github.com/akaimo/sample-controller/pkg/client/informers/externalversions/samplecontroller/v1"
	listers "github.com/akaimo/sample-controller/pkg/client/listers/samplecontroller/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	jobinformers "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	batchlisters "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"time"
)

const controllerAgentName = "sample-controller"

type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	jobLister            batchlisters.JobLister
	jobSynced            cache.InformerSynced
	sampleResourceLister listers.SampleResourceLister
	sampleResourceSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	jobInformer jobinformers.JobInformer,
	sampleinformer informers.SampleResourceInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:        kubeclientset,
		sampleclientset:      sampleclientset,
		jobLister:            jobInformer.Lister(),
		jobSynced:            jobInformer.Informer().HasSynced,
		sampleResourceLister: sampleinformer.Lister(),
		sampleResourceSynced: sampleinformer.Informer().HasSynced,
		workqueue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Samples"),
		recorder:             recorder,
	}

	klog.Info("Setting up event handlers")
	sampleinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueSampleResource,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueSampleResource(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting sample controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.sampleResourceSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)

		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	sr, err := c.sampleResourceLister.SampleResources(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("foo '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	jobList, err := c.jobLister.Jobs(sr.Namespace).List(labels.Everything())

	duration, err := time.ParseDuration(sr.Spec.Time)
	if err != nil {
		return err
	}
	var deleteList []*batchv1.Job

	for _, v := range jobList {
		if expired, err := c.processTTL(v, duration); err != nil {
			klog.Error(err)
			continue
		} else if !expired {
			continue
		}
		deleteList = append(deleteList, v)
	}

	for _, v := range deleteList {
		policy := metav1.DeletePropagationForeground
		options := &metav1.DeleteOptions{
			PropagationPolicy: &policy,
			Preconditions:     &metav1.Preconditions{UID: &v.UID},
		}
		klog.V(4).Infof("Cleaning up Job %s/%s", v.Namespace, v.Name)
		err = c.kubeclientset.BatchV1().Jobs(v.Namespace).Delete(v.Name, options)
		if err != nil {
			klog.Error(err)
		}
	}

	return nil
}

func (c *Controller) enqueueSampleResource(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) processTTL(job *batchv1.Job, ttl time.Duration) (expired bool, err error) {
	// We don't care about the Jobs that are going to be deleted, or the ones that don't need clean up.
	if job.DeletionTimestamp != nil || !isJobFinished(job) {
		return false, nil
	}

	//now := tc.clock.Now()
	now := time.Now()
	t, err := timeLeft(job, &now, ttl)
	if err != nil {
		return false, err
	}

	// TTL has expired
	if *t <= 0 {
		return true, nil
	}

	return false, nil
}

func isJobFinished(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func getFinishAndExpireTime(j *batchv1.Job, ttl time.Duration) (*time.Time, *time.Time, error) {
	if !isJobFinished(j) {
		return nil, nil, fmt.Errorf("job %s/%s should not be cleaned up", j.Namespace, j.Name)
	}
	finishAt, err := jobFinishTime(j)
	if err != nil {
		return nil, nil, err
	}
	finishAtUTC := finishAt.UTC()
	expireAtUTC := finishAtUTC.Add(ttl)
	return &finishAtUTC, &expireAtUTC, nil
}

func timeLeft(j *batchv1.Job, since *time.Time, ttl time.Duration) (*time.Duration, error) {
	finishAt, expireAt, err := getFinishAndExpireTime(j, ttl)
	if err != nil {
		return nil, err
	}
	if finishAt.UTC().After(since.UTC()) {
		klog.Warningf("Warning: Found Job %s/%s finished in the future. This is likely due to time skew in the cluster. Job cleanup will be deferred.", j.Namespace, j.Name)
	}
	remaining := expireAt.UTC().Sub(since.UTC())
	klog.V(4).Infof("Found Job %s/%s finished at %v, remaining TTL %v since %v, TTL will expire at %v", j.Namespace, j.Name, finishAt.UTC(), remaining, since.UTC(), expireAt.UTC())
	return &remaining, nil
}

func jobFinishTime(finishedJob *batchv1.Job) (metav1.Time, error) {
	for _, c := range finishedJob.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
			finishAt := c.LastTransitionTime
			if finishAt.IsZero() {
				return metav1.Time{}, fmt.Errorf("unable to find the time when the Job %s/%s finished", finishedJob.Namespace, finishedJob.Name)
			}
			return c.LastTransitionTime, nil
		}
	}

	return metav1.Time{}, fmt.Errorf("unable to find the status of the finished Job %s/%s", finishedJob.Namespace, finishedJob.Name)
}
