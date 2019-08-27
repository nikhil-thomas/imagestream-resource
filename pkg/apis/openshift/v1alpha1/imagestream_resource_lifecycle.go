package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
)

var condSet = apis.NewLivingConditionSet()

// GetGroupVersionKind implements kmeta.OwnerRefable
func (isr *ImagestreamResource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemaGroupVersion.WithKind("ImagestreamResource")
}

func (isrs *ImagestreamResourceStatus) InitializeContditions() {
	condSet.Manage(isrs).InitializeConditions()
}

func (isrs *ImagestreamResourceStatus) MarkImagestreamUnavailable(name string) {
	condSet.Manage(isrs).MarkFalse(
		ImagestreamResourceReady,
		"Imagestream not ready",
		"Imagestream %s not found", name)
}

func (isrs *ImagestreamResourceStatus) MarkImagestreamAvailable() {
	condSet.Manage(isrs).MarkTrue(ImagestreamResourceReady)
}
