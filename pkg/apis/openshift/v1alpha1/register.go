package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift/imagestream-resource/pkg/apis/openshift"
)

// SchemaGroupVersion is group version used to register these objects
var SchemaGroupVersion = schema.GroupVersion{
	Group:   openshift.GroupName,
	Version: "v1alpha1",
}

// Kind takes an unqualified resource and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemaGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemaGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemaGroupVersion,
		&ImagestreamResource{},
		&ImagestreamResourceList{},
		)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme = SchemeBuilder.AddToScheme
)
