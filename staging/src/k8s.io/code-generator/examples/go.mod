// This is a submodule to isolate k8s.io/code-generator from k8s.io/{api,apimachinery,client-go} dependencies in generated code

module k8s.io/code-generator/examples

go 1.16

require (
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
)

replace (
	k8s.io/api => ../../api
	k8s.io/apimachinery => ../../apimachinery
	k8s.io/client-go => ../../client-go
	k8s.io/kube-openapi => github.com/austince/kube-openapi v0.0.0-restframework-e9ebba8
)

replace sigs.k8s.io/json => github.com/liggitt/json v0.0.0-20211020163728-48258682683b
