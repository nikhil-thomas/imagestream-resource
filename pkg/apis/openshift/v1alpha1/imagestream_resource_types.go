package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/kmeta"
)

// +genclient
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

// ImagestreamResourceSpec holds the image url, name, tag/digest of the
// container image which has to be imported/tracked as an image stream
// (NT: aug 23 2019): this assumption can change
type ImagestreamResourceSpec struct {
	Name   string  `json:"name"`
	Params []Param `json:"params"`
}

// ImagestreamResourceStatus holds the PipelineExtensibility Contract
// so that ImagestreamResource can be used as a PipelineResource by Tekton-Pipelines
type ImagestreamResourceStatus struct {
	duckv1beta1.Status `json:",inline"`
	Conditions         apis.Conditions    `json:"conditions"`
	Beforecontainers   []corev1.Container `json:"beforeContainers"`
	Aftercontainers    []corev1.Container `json:"afterContainers"`
	Params             []Param            `json:"params"`
}

// Param are key value pair which can be
// used by consuming controller for value substitution
// Param declares a value to use for the Param called Name.
type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImagestreamResourceList is a list of ImagestreamResource resources
type ImagestreamResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ImagestreamResource `json:"items"`
}
