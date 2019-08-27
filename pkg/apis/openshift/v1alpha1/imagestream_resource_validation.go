package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable
func (isr *ImagestreamResource) Validate(ctx context.Context) *apis.FieldError {
	return isr.Spec.Validate(ctx).ViaField("spec")
}

// Validate implements apis.Validate
func (isrsp *ImagestreamResourceSpec) Validate(ctx context.Context) *apis.FieldError {
	if isrsp.ImageName == "" {
		return apis.ErrMissingField("imageName")
	}
	return nil
}
