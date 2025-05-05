package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func rkeClusterServicesKubeControllerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "CIDR Range for Pods in cluster",
		},
		"extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the kube-controller service",
		},
		"windows_extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the kube-controller service",
		},
		"extra_args_array": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A JSON Representation of extra kube-api service arguments which can be passed multiple times",
		},
		"windows_extra_args_array": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A JSON Representation of extra kube-api service arguments which can be passed multiple times for windows nodes",
		},
		"extra_binds": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Extra binds added to the controlplane nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Extra env added to the controlplane nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker image of the kube-controller service",
		},
		"service_cluster_ip_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Virtual IP range that will be used by Kubernetes services",
		},
	}
	return s
}
