.PHONY: lcoal-dev
local-dev:
	SYSTEM_NAMESPACE=openshift-pipelines \
	METRICS_DOMAIN=openshift/imagestream-resource \
	go run cmd/controller/main.go
