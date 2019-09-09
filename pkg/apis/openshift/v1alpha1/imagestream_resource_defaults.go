package v1alpha1

import (
	"context"
)

func (isr *ImagestreamResource) SetDefaults(ctx context.Context) {
	isr.Spec.Namespace = isr.ObjectMeta.Namespace
}
