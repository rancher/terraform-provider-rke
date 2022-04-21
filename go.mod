module github.com/rancher/terraform-provider-rke

go 1.16

require (
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/go-version v1.3.0
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/rancher/rke v1.3.9
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.23.3
	k8s.io/apimachinery v0.23.3
	k8s.io/apiserver v0.23.3
	k8s.io/client-go v0.23.3
)

replace k8s.io/client-go => k8s.io/client-go v0.23.3

replace github.com/spf13/afero => github.com/spf13/afero v1.2.2
