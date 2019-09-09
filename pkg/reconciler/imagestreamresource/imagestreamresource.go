package imagestreamresource

import (
	"context"
	"os"
	"reflect"
	"strings"

	"knative.dev/pkg/tracker"

	"github.com/openshift/imagestream-resource/pkg/apis/openshift/v1alpha1"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/equality"

	"k8s.io/client-go/tools/cache"

	"knative.dev/pkg/logging"

	clientset "github.com/openshift/imagestream-resource/pkg/client/clientset/versioned"
	listers "github.com/openshift/imagestream-resource/pkg/client/listers/openshift/v1alpha1"
	pipelinev1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/controller"
)

// Reconciler implements controller.Reconciler for AddressableService resources
type Reconciler struct {
	// Client is used to write status updates
	Client clientset.Interface

	Lister listers.ImagestreamResourceLister

	// The tracker builds an index of what resources are watching other
	// resources so that we can immediately react to changes to changes in
	// tracked resources
	Tracker tracker.Interface

	// Recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	Recorder record.EventRecorder
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*Reconciler)(nil)

func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	logger := logging.FromContext(ctx)
	logger.Infof("key: %s \n", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	original, err := r.Lister.ImagestreamResources(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		logger.Errorf("resource %q no longer exists", key)
		return nil
	} else if err != nil {
		return err
	}
	logger.Info("reconcile isr fetched")

	resource := original.DeepCopy()

	reconcileErr := r.reconcile(ctx, resource)

	if equality.Semantic.DeepEqual(original.Status, resource.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.
		logger.Info("reconcile equality")
	} else if _, err = r.updateStatus(resource); err != nil {
		logger.Warnw("Failed to update resource status", zap.Error(err))
		r.Recorder.Eventf(resource, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update resource status", resource.Name, err)
		return err
	}

	if reconcileErr != nil {
		r.Recorder.Eventf(
			resource,
			corev1.EventTypeWarning,
			"InternalError",
			reconcileErr.Error())
		logger.Errorf("reconcile error :%v:", err)
	}
	logger.Info("reconcile return err: nil")
	return reconcileErr
}

func (r *Reconciler) reconcile(ctx context.Context, isr *v1alpha1.ImagestreamResource) error {
	if isr.GetDeletionTimestamp() != nil {
		// if deletion timestamp is set skip reconcile logic
		// if we need to add finalizer logic, add it here
		return nil
	}
	isr.Status.InitializeContditions()

	if err := r.reconcileImagestreamResource(ctx, isr); err != nil {
		return err
	}
	return nil
}

func (r *Reconciler) reconcileImagestreamResource(ctx context.Context, resource *v1alpha1.ImagestreamResource) error {
	logger := logging.FromContext(ctx)

	registryRoute := os.Getenv("OPENSHIFT_IMAGE_REGISTRY")

	imageStrURL := registryRoute + "/" + resource.Spec.Namespace
	imageStrURL += "/" + resource.Spec.Name
	logger.Infof("Imagestream URL : %s", imageStrURL)

	statusParam := pipelinev1alpha1.ResourceParam{
		Name:  "url",
		Value: imageStrURL,
	}
	index := -1

	for i, rscPrm := range resource.Status.Variables {
		if strings.EqualFold(rscPrm.Name, "URL") {
			index = i
			break
		}
	}
	if index >= 0 {
		resource.Status.Variables = append(resource.Status.Variables[:index], resource.Status.Variables[index+1:]...)
	}

	resource.Status.Variables = append(resource.Status.Variables, statusParam)
	return nil
}

// Update the Status of the resource.  Caller is responsible for checking
// for semantic differences before calling.
func (r *Reconciler) updateStatus(desired *v1alpha1.ImagestreamResource) (*v1alpha1.ImagestreamResource, error) {
	actual, err := r.Lister.ImagestreamResources(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}

	// if there is nothing to update, just return
	if reflect.DeepEqual(desired.Status, actual.Status) {
		return actual, nil
	}

	existing := actual.DeepCopy()
	existing.Status = desired.Status
	return r.Client.OpenshiftV1alpha1().ImagestreamResources(desired.Namespace).UpdateStatus(existing)
}
