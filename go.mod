module samp

go 1.15

require (
	github.com/go-logr/logr v1.2.3
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	istio.io/api v0.0.0-20200512234804-e5412c253ffe
	istio.io/client-go v0.0.0-20200513000250-b1d6e9886b7b
	k8s.io/api v0.26.0
	k8s.io/apiextensions-apiserver v0.23.0
	k8s.io/apimachinery v0.26.0
	k8s.io/client-go v0.23.5
	sigs.k8s.io/controller-runtime v0.11.0
)

replace (
	k8s.io/api v0.26.0 => k8s.io/api v0.23.0
	k8s.io/apimachinery v0.26.0 => k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.5 => k8s.io/client-go v0.23.0
)
