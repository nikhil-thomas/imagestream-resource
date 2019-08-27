package imagestreamresource

import (
	"context"

	"knative.dev/pkg/tracker"

	"github.com/openshift/imagestream-resource/pkg/client/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	"knative.dev/pkg/logging"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	isrclient "github.com/openshift/imagestream-resource/pkg/client/injection/client"
	isrinf "github.com/openshift/imagestream-resource/pkg/client/injection/informers/openshift/v1alpha1/imagestreamresource"
)

const (
	controllerAgentName = "imagestreamresource-controller"
)

func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	isrInformer := isrinf.Get(ctx)

	c := &Reconciler{
		Client: isrclient.Get(ctx),
		Lister: isrInformer.Lister(),
		Recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, corev1.EventSource{Component: controllerAgentName}),
	}

	impl := controller.NewImpl(c, logger, "ImagestreamResources")

	logger.Info("Setting up event handlers")
	isrInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
	c.Tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	// c.Tracker.Onchanged to subresource imformer call

	return impl
}
