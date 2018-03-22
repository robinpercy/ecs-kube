SHELL = sh -xv

REPO?=robinpercy
IMG?=ecs-kube
TAG?=latest
FULL_IMG="$(REPO)/$(IMG):$(TAG)"
NS_DEMO?="default"
SVC_DEMO?="ecs-kube-apiserver"


build-image:
	apiserver-boot build container --generate=false --image $(FULL_IMG)

push-image: build-image 
	docker push $(FULL_IMG)

# Regenerates api-machinery code
build-generated:
	apiserver-boot build generated

# Builds the deployment config for the aggregated api server and an etcd store to back it
.PHONY: build-agg-config
build-agg-config:
	rm -rf config
	apiserver-boot build config --name $(SVC_DEMO) --namespace $(NS_DEMO) --image $(FULL_IMG)

# Runs the aggregated api and clears the discovery cache before outputing registered apis
.PHONY: run-agg-api
run-agg-api: build-agg-config 
	kubectl apply -f config/apiserver.yaml
	rm -rf ~/.kube/cache/discovery/
	kubectl api-versions
