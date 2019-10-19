.PHONY: codegen
codegen:
	${GOPATH}/pkg/mod/k8s.io/code-generator@v0.0.0-20190831074504-732c9ca86353/generate-groups.sh "deepcopy,client,informer,lister" \
		github.com/akaimo/sample-controller/pkg/client \
		github.com/akaimo/sample-controller/pkg/apis \
		samplecontroller:v1 \
		--go-header-file  hack/boilerplate.go.txt
