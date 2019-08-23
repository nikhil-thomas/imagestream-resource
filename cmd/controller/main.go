package main

import (
	"knative.dev/pkg/injection/sharedmain"
)

func main() {

	sharedmain.Main("controller")
}
