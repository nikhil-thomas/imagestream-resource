.PHONY: lcoal-dev
local-dev:
	SYSTEM_NAMESPACE=openshift-pipelines \
	METRICS_DOMAIN=openshift/imagestream-resource \
	go run cmd/controller/main.go

.PHONY: ko
ko-app:
	ko apply -f config --watch
