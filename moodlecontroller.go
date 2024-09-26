/*
Copyright 2017 The Kubernetes Authors.

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
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/go-logr/logr"
	operatorv1 "github.com/smugug/keysaas/pkg/apis/moodlecontroller/v1"
	clientset "github.com/smugug/keysaas/pkg/generated/clientset/versioned"
	operatorscheme "github.com/smugug/keysaas/pkg/generated/clientset/versioned/scheme"
	informers "github.com/smugug/keysaas/pkg/generated/informers/externalversions"
	listers "github.com/smugug/keysaas/pkg/generated/listers/moodlecontroller/v1"
	"github.com/smugug/keysaas/pkg/utils"
)

const controllerAgentName = "moodle-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
	// FieldManager distinguishes this controller from other things writing to API objects
	FieldManager = controllerAgentName
)

var (
	MOODLE_PORT_BASE = 30060
	MOODLE_PORT      int
)

func init() {
}

// Controller is the controller implementation for Foo resources
type MoodleController struct {
	//for logging
	ctx context.Context
	//For ultis
	cfg *restclient.Config
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	foosLister        listers.MoodleLister
	foosSynced        cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.TypedRateLimitingInterface[cache.ObjectName]
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
	util     utils.Utils
	logger   logr.Logger
}

// NewController returns a new sample controller
func NewMoodleController(
	ctx context.Context,
	cfg *restclient.Config,
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	moodleInformerFactory informers.SharedInformerFactory) *MoodleController {

	logger := klog.FromContext(ctx)

	// obtain references to shared index informers for the Deplsoyment and Foo
	// types.
	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	moodleInformer := moodleInformerFactory.Moodlecontroller().V1().Moodles()

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(operatorscheme.AddToScheme(scheme.Scheme))
	logger.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster(record.WithContext(ctx))
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	ratelimiter := workqueue.NewTypedMaxOfRateLimiter(workqueue.NewTypedItemExponentialFailureRateLimiter[cache.ObjectName](5*time.Millisecond, 1000*time.Second))

	utils := utils.NewUtils(cfg, kubeclientset)

	controller := &MoodleController{
		cfg:               cfg,
		ctx:               ctx,
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		foosLister:        moodleInformer.Lister(),
		foosSynced:        moodleInformer.Informer().HasSynced,
		workqueue:         workqueue.NewTypedRateLimitingQueue(ratelimiter),
		recorder:          recorder,
		util:              utils,
		logger:            logger,
	}

	logger.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	// fooInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc: controller.enqueueFoo,
	// 	UpdateFunc: func(old, new interface{}) {
	// 		controller.enqueueFoo(new)
	// 	},
	// })
	moodleInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueFoo(new)
			// newDepl := new.(*operatorv1.Moodle)
			// oldDepl := old.(*operatorv1.Moodle)
			// //fmt.Println("MoodleController.go  : New Version:%s", newDepl.ResourceVersion)
			// //fmt.Println("MoodleController.go  : Old Version:%s", oldDepl.ResourceVersion)
			// if newDepl.ResourceVersion == oldDepl.ResourceVersion {
			// 	// Periodic resync will send update events for all known Deployments.
			// 	// Two different versions of the same Deployment will always have different RVs.
			// 	return
			// } else {
			// 	controller.enqueueFoo(new)
			// }
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Foo resource then the handler will enqueue that Foo resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	// deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc: controller.handleObject,
	// 	UpdateFunc: func(old, new interface{}) {
	// 		newDepl := new.(*appsv1.Deployment)
	// 		oldDepl := old.(*appsv1.Deployment)
	// 		if newDepl.ResourceVersion == oldDepl.ResourceVersion {
	// Periodic resync will send update events for all known Deployments.
	// Two different versions of the same Deployment will always have different RVs.
	// 			return
	// 		}
	// 		controller.handleObject(new)
	// 	},
	// 	DeleteFunc: controller.handleObject,
	// })

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *MoodleController) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()
	// Start the informer factories to begin populating the informer caches
	c.logger.Info("Starting Moodle controller")

	// Wait for the caches to be synced before starting workers
	c.logger.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(ctx.Done(), c.deploymentsSynced, c.foosSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	c.logger.Info("Starting workers", "workers: ", workers)
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	c.logger.Info("Started workers")
	<-ctx.Done()
	c.logger.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *MoodleController) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *MoodleController) processNextWorkItem(ctx context.Context) bool {
	objRef, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We call Done at the end of this func so the workqueue knows we have
	// finished processing this item. We also must remember to call Forget
	// if we do not want this work item being re-queued. For example, we do
	// not call Forget if a transient error occurs, instead the item is
	// put back on the workqueue and attempted again after a back-off
	// period.
	defer c.workqueue.Done(objRef)

	// Run the syncHandler, passing it the structured reference to the object to be synced.
	err := c.syncHandler(ctx, objRef)
	if err == nil {
		// If no error occurs then we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(objRef)
		c.logger.Info("Successfully synced", "objectName", objRef.Name)
		return true
	}
	// there was a failure so be sure to report it.  This method allows for
	// pluggable error handling which can be used for things like
	// cluster-monitoring.
	utilruntime.HandleErrorWithContext(ctx, err, "Error syncing; requeuing for later retry", "objectReference", objRef)
	// since we failed, we should requeue the item to work on later.  This
	// method will add a backoff to avoid hotlooping on particular items
	// (they're probably still not going to work right away) and overall
	// controller protection (everything I've done is broken, this controller
	// needs to calm down or it can starve other useful work) cases.
	c.workqueue.AddRateLimited(objRef)
	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *MoodleController) syncHandler(ctx context.Context, objRef cache.ObjectName) error {
	c.logger.Info("Start sync handling", "objRef", objRef)
	// Convert the namespace/name string into a distinct namespace and name
	namespace := objRef.Namespace
	name := objRef.Name
	// Get the Foo resource with this namespace/name
	foo, err := c.foosLister.Moodles(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleErrorWithContext(ctx, err, "Foo referenced by item in work queue no longer exists", "name", name)
			return nil
		}

		return err
	}

	c.logger.Info("MoodleController.go  : **************************************")
	moodleName := foo.Name
	moodleNamespace := foo.Namespace
	c.logger.Info("MoodleController.go", "moodle name", moodleName)
	c.logger.Info("MoodleController.go", "moodle namespace", moodleNamespace)

	var status, url string
	initialDeployment := c.isInitialDeployment(foo)

	if foo.Status.Status != "" && foo.Status.Status == "Moodle Pod Timeout" {
		return nil
	}

	if initialDeployment {

		moodleDomainName := foo.Spec.DomainName

		c.logger.Info("MoodleController.go", "domain name", moodleDomainName)

		MOODLE_PORT = MOODLE_PORT_BASE
		MOODLE_PORT_BASE = MOODLE_PORT_BASE + 1

		c.logger.Info("MoodleController.go: Deploying Moodle", "port", MOODLE_PORT)

		initialDeployment = false

		serviceURL, podName, secretName, err := c.deployMoodle(ctx, foo)

		//var correctlyInstalledPlugins []string
		if err != nil {
			status = err.Error()
		} else {
			status = "Ready"
			url = "http://" + serviceURL
			c.logger.Info("MoodleController.go : Ready", "moodle URL", url)
			//supportedPlugins, unsupportedPlugins = c.util.GetSupportedPlugins(plugins)
			//correctlyInstalledPlugins = c.getDiff(supportedPlugins, erredPlugins)
		}
		c.logger.Info("MoodleController.go : Updating status")
		err = c.updateMoodleStatus(ctx, foo, podName, secretName, status, url)
		if err != nil {
			c.logger.Info("MoodleController.go : Updating error", "error", err)
			return err
		}
	} else {
		c.logger.Info("MoodleController.go : Moodle custom resource did not change", "moodle name", moodleName)
	}

	// else {
	// 	podName, installedPlugins, unsupportedPluginsCurrent, erredPlugins := c.handlePluginDeployment(foo)
	// 	if len(installedPlugins) > 0 || len(unsupportedPluginsCurrent) > 0 {
	// 		status = "Ready"
	// 		url = foo.Status.Url
	// 		unsupportedPlugins = foo.Status.UnsupportedPlugins
	// 		unsupportedPlugins = appendList(unsupportedPluginsCurrent, unsupportedPlugins)

	// 		supportedPlugins = foo.Status.InstalledPlugins
	// 		supportedPlugins = append(supportedPlugins, installedPlugins...)
	// 		if len(erredPlugins) > 0 {
	// 			c.updateMoodleStatus(foo, podName, "", "Error in installing some plugins", url, &supportedPlugins, &unsupportedPlugins)
	// 		} else {
	// 			c.updateMoodleStatus(foo, podName, "", status, url, &supportedPlugins, &unsupportedPlugins)
	// 		}
	// 		c.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	// 	} else {
	// 		c.logger.Info("MoodleController.go  : Moodle custom resource did not change. No plugin installed.", "moodle name", moodleName)
	// 	}
	// }

	// deploymentName := foo.Spec.DeploymentName
	// if deploymentName == "" {
	// 	// We choose to absorb the error here as the worker would requeue the
	// 	// resource otherwise. Instead, the next time the resource is updated
	// 	// the resource will be queued again.
	// 	utilruntime.HandleErrorWithContext(ctx, nil, "Deployment name missing from object reference", "objectReference", objectRef)
	// 	return nil
	// }

	// // Get the deployment with the name specified in Foo.spec
	// deployment, err := c.deploymentsLister.Deployments(foo.Namespace).Get(deploymentName)
	// // If the resource doesn't exist, we'll create it
	// if errors.IsNotFound(err) {
	// 	deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Create(context.TODO(), newDeployment(foo), metav1.CreateOptions{FieldManager: FieldManager})
	// }

	// // If an error occurs during Get/Create, we'll requeue the item so we can
	// // attempt processing again later. This could have been caused by a
	// // temporary network failure, or any other transient reason.
	// if err != nil {
	// 	return err
	// }

	// // If the Deployment is not controlled by this Foo resource, we should log
	// // a warning to the event recorder and return error msg.
	// if !metav1.IsControlledBy(deployment, foo) {
	// 	msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
	// 	c.recorder.Event(foo, corev1.EventTypeWarning, ErrResourceExists, msg)
	// 	return fmt.Errorf("%s", msg)
	// }

	// // If this number of the replicas on the Foo resource is specified, and the
	// // number does not equal the current desired replicas on the Deployment, we
	// // should update the Deployment resource.
	// if foo.Spec.Replicas != nil && *foo.Spec.Replicas != *deployment.Spec.Replicas {
	// 	logger.V(4).Info("Update deployment resource", "currentReplicas", *foo.Spec.Replicas, "desiredReplicas", *deployment.Spec.Replicas)
	// 	deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Update(context.TODO(), newDeployment(foo), metav1.UpdateOptions{FieldManager: FieldManager})
	// }

	// // If an error occurs during Update, we'll requeue the item so we can
	// // attempt processing again later. This could have been caused by a
	// // temporary network failure, or any other transient reason.
	// if err != nil {
	// 	return err
	// }

	// // Finally, we update the status block of the Foo resource to reflect the
	// // current state of the world
	// err = c.updateFooStatus(foo, deployment)
	// if err != nil {
	// 	return err
	// }

	c.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// enqueueFoo takes a Foo resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (c *MoodleController) enqueueFoo(obj interface{}) {
	if objectRef, err := cache.ObjectToName(obj); err != nil {
		utilruntime.HandleError(err)
		return
	} else {
		c.workqueue.Add(objectRef)
	}
}

func (c *MoodleController) updateMoodleStatus(ctx context.Context, foo *operatorv1.Moodle, podName, secretName, status string, url string) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
	fooCopy.Status.PodName = podName
	if secretName != "" {
		fooCopy.Status.SecretName = secretName
	}
	fooCopy.Status.Status = status
	fooCopy.Status.Url = url
	// fooCopy.Status.InstalledPlugins = *plugins
	// fooCopy.Status.UnsupportedPlugins = *unsupportedPlugins
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the Foo resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := c.sampleclientset.MoodlecontrollerV1().Moodles(foo.Namespace).Update(ctx, fooCopy, metav1.UpdateOptions{})
	return err
}

func appendList(source, destination []string) []string {
	var appendedList []string

	for _, delem := range destination {
		present := false
		for _, selem := range source {
			if delem == selem {
				present = true
				break
			}
		}
		if !present {
			appendedList = append(appendedList, delem)
		}
	}
	return appendedList
}
