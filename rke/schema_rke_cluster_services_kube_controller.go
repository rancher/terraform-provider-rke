package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		"win_extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments for Windows systems that are added to the scheduler services",
		},
		"extra_args_array": {
			Type:        schema.TypeSet,
			Optional:    true,
			MaxItems:    1,
			Description: "Extra arguments that can be specified multiple times which are added to the kube-controller service",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeControllerExtraArgsArrayFields(),
			},
		},
		"win_extra_args_array": {
			Type:        schema.TypeSet,
			Optional:    true,
			MaxItems:    1,
			Description: "Extra arguments for Windows systems that can be specified multiple times which are added to the kube-controller service",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeControllerExtraArgsArrayFields(),
			},
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

func rkeClusterServicesKubeControllerExtraArgsArrayFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"extra_arg": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"argument": {
						Required: true,
						Type:     schema.TypeString,
					},
					"values": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}
