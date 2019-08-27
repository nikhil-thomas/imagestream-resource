package main

import (
	"github.com/openshift/imagestream-resource/pkg/reconciler/imagestreamresource"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("controller",
		imagestreamresource.NewController)
}
