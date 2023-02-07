module github.com/mongodb/mongodb-atlas-kubernetes

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v61.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.19
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.10
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/aws/aws-sdk-go v1.42.25
	github.com/fatih/structtag v1.2.0
	github.com/go-logr/zapr v1.2.3
	github.com/google/go-cmp v0.5.9
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mongodb-forks/digest v1.0.3
	github.com/mxschmitt/playwright-go v0.1400.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.24.1
	github.com/pborman/uuid v1.2.1
	github.com/sethvargo/go-password v0.2.0
	github.com/stretchr/testify v1.8.0
	go.mongodb.org/atlas v0.7.3-0.20210315115044-4b1d3f428c24
	go.mongodb.org/mongo-driver v1.7.4
	go.uber.org/zap v1.24.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.26.1
	k8s.io/apimachinery v0.26.1
	k8s.io/client-go v0.26.1
	sigs.k8s.io/controller-runtime v0.14.4
)
