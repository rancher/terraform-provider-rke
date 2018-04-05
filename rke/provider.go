package rke

import (
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sirupsen/logrus"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"log": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RKE_LOG", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rke_cluster": resourceRKECluster(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	logrus.SetOutput(ioutil.Discard)
	if v, ok := d.GetOk("log"); ok {
		if v.(bool) {
			logrus.SetOutput(os.Stderr)
		}
	}
	config := Config{}
	return config, nil
}
