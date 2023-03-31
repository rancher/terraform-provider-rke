package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	rkeClusterRotateCertificatesServicesEtcd           = "etcd"
	rkeClusterRotateCertificatesServicesKubelet        = "kubelet"
	rkeClusterRotateCertificatesServicesKubeAPI        = "kube-apiserver"
	rkeClusterRotateCertificatesServicesKubeProxy      = "kube-proxy"
	rkeClusterRotateCertificatesServicesKubeScheduler  = "kube-scheduler"
	rkeClusterRotateCertificatesServicesKubeController = "kube-controller-manager"
)

var (
	rkeClusterRotateCertificatesServicesList = []string{
		rkeClusterRotateCertificatesServicesEtcd,
		rkeClusterRotateCertificatesServicesKubelet,
		rkeClusterRotateCertificatesServicesKubeAPI,
		rkeClusterRotateCertificatesServicesKubeProxy,
		rkeClusterRotateCertificatesServicesKubeScheduler,
		rkeClusterRotateCertificatesServicesKubeController,
	}
)

//Schemas

func rkeClusterRotateCertificatesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ca_certificates": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Rotate CA Certificates",
		},
		"services": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Services to rotate their certs. valid values are etcd|kubelet|kube-apiserver|kube-proxy|kube-scheduler|kube-controller-manager",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(rkeClusterRotateCertificatesServicesList, true),
			},
		},
	}
	return s
}
