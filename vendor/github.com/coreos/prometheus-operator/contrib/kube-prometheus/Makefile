image:
	docker build -f ../../scripts/jsonnet/Dockerfile -t po-jsonnet ../../

generate: image
	@echo ">> Compiling assets and generating Kubernetes manifests"
	docker run --rm -u=$(shell id -u $(USER)):$(shell id -g $(USER)) -v $(shell dirname $(dir $(abspath $(dir $$PWD)))):/go/src/github.com/coreos/prometheus-operator/ --workdir /go/src/github.com/coreos/prometheus-operator/contrib/kube-prometheus po-jsonnet make crdtojsonnet generate-raw

crdtojsonnet:
	cat ../../example/prometheus-operator-crd/alertmanager.crd.yaml | gojsontoyaml -yamltojson > jsonnet/kube-prometheus/prometheus-operator/alertmanager-crd.libsonnet
	cat ../../example/prometheus-operator-crd/prometheus.crd.yaml | gojsontoyaml -yamltojson > jsonnet/kube-prometheus/prometheus-operator/prometheus-crd.libsonnet
	cat ../../example/prometheus-operator-crd/servicemonitor.crd.yaml | gojsontoyaml -yamltojson > jsonnet/kube-prometheus/prometheus-operator/servicemonitor-crd.libsonnet

generate-raw:
	jb install
	./build.sh

.PHONY: image generate crdtojsonnet generate-raw
