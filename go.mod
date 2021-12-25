module github.com/birdayz/nodered-operator

go 1.16

require (
	github.com/banzaicloud/k8s-objectmatcher v1.7.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/banzaicloud/k8s-objectmatcher => github.com/birdayz/k8s-objectmatcher v1.7.1-0.20211225085715-5cba284c3365
