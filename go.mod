module github.com/rancher/terraform-provider-rke

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/rancher/rke v1.1.3
	github.com/rancher/types v0.0.0-20200609171948-b18f4c194419
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/apiserver v0.18.0
	k8s.io/client-go v12.0.0+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.18.0
