package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmeta"
)

// +genClient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImagestreamResource is duck-typed Tekton Pipeline Resource
// which enables Tekton Pipeline use Openshit Imagestreams
type ImagestreamResource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:",inline"`

	// Spec holds the image url, name, tag/digest of the
	// container image which has to be imported/tracked as an image stream
	// (NT: aug 23 2019): this assumption can change
	Spec ImagestreamResourceSpec `json:"spec,omitempty"`

	// Status holds the PipelineExtensibility Contract
	// so that ImagestreamResource can be used as a PipelineResource by Tekton-Pipelines
	Status ImagestreamResourceStatus `json:"status,omitempty"`
}

const (
	// ImagestreamResourceReady is set when an Imagestream Resoruce is
	// ready to be consumed by Tekton-Pipelines
	ImagestreamResourceReady = apis.ConditionReady
)

// check ImagestreamResource can be validated and defaulted
var _ apis.Validatable = (*ImagestreamResource)(nil)
var _ apis.Defaultable = (*ImagestreamResource)(nil)
var _ kmeta.OwnerRefable = (*ImagestreamResource)(nil)

type ImagestreamResourceSpec struct {
	ImageName string
}

type ImagestreamResourceStatus struct {
	Conditions       apis.Conditions `json:"conditions"`
	beforecontainers []corev1.Container
	aftercontainers  []corev1.Container
	variables Variables
}

type Variables struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
