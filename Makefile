.PHONY: codegen
codegen:
	${GOPATH}/pkg/mod/k8s.io/code-generator@v0.0.0-20190831074504-732c9ca86353/generate-groups.sh "deepcopy,client,informer,lister" \
		github.com/akaimo/sample-controller/pkg/client \
		github.com/akaimo/sample-controller/pkg/apis \
		samplecontroller:v1 \
		--go-header-file  hack/boilerplate.go.txt

.PHONY: run
run:
	go run . -kubeconfig=${HOME}/.kube/config

build:
	docker build -t akaimo/sample-controller:0.1.0 .

push: build
	docker push akaimo/sample-controller:0.1.0

gen-template:
	helm template ./helm --name-template sample-controller --namespace sample-controller > ./artifacts/02-controller.yaml 
	cp ./helm/crds/sample-controller.yaml ./artifacts/00-crds.yaml
