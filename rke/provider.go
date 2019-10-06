package rke

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sirupsen/logrus"
)

var rkeLogBuf = &bytes.Buffer{}

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
		DataSourcesMap: map[string]*schema.Resource{
			"rke_node_parameter": dataSourceRKENodeParameter(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	initLogger(d)

	config := Config{}
	return config, nil
}

func initLogger(d *schema.ResourceData) {
	var writer io.Writer = rkeLogBuf
	if v, ok := d.GetOk("log"); ok {
		if v.(bool) {
			writer = io.MultiWriter(os.Stderr, rkeLogBuf)
		}
	}

	logrus.SetOutput(writer)
	logrus.SetFormatter(&logFormatter{})
}

type logFormatter struct{}

func (l *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s\n", entry.Level, entry.Message)), nil
}
