SHELL := /bin/bash

gwapi:	
	wget -O crds/gateway-api/v1-standard.yaml https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml
	wget -O crds/gateway-api/v1-experimental.yaml https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/experimental-install.yaml
istio:
	wget -O crds/istio/crd-all.gen.yaml https://raw.githubusercontent.com/istio/istio/master/manifests/charts/base/crds/crd-all.gen.yaml
tencent:
	./script/dump-crd.sh tencent
alibaba:
	./script/dump-crd.sh alibaba
ocm:
	mkdir -p crds/ocm
	git clone --depth 1 https://github.com/open-cluster-management-io/ocm.git
	cp ocm/manifests/cluster-manager/hub/*.crd.yaml crds/ocm/
	rm -rf ocm
kubevela:
	mkdir -p crds/kubevela
	git clone --depth 1 https://github.com/kubevela/kubevela.git
	cp kubevela/charts/vela-core/crds/*.yaml crds/kubevela/
	rm -rf kubevela
clusternet:
	mkdir -p crds/clusternet
	git clone --depth 1 https://github.com/clusternet/clusternet.git
	cp clusternet/manifests/crds/*.yaml crds/clusternet/
	rm -rf clusternet
karmada:
	mkdir -p crds/karmada
	git clone --depth 1 https://github.com/karmada-io/karmada.git
	cp -r karmada/charts/karmada/_crds/bases/* crds/karmada/
	rm -rf karmada
index:
	go run . index
crd:
	go run . crd
regen: crd index
