package rke

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRKEClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	files, err := splitImportID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	clusterYamlBytes, err := ioutil.ReadFile(files[0])
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Reading RKE config file %s: %v", files[0], err)
	}
	if len(clusterYamlBytes) == 0 {
		return []*schema.ResourceData{}, fmt.Errorf("RKE config is nil")
	}
	_, err = yamlToMapInterface(string(clusterYamlBytes))
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Unmarshaling RKE config yaml: %v", err)
	}

	clusterStateBytes, err := ioutil.ReadFile(files[1])
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Reading RKE state file %s: %v", files[0], err)
	}
	if len(clusterStateBytes) == 0 {
		return []*schema.ResourceData{}, fmt.Errorf("RKE state is nil")
	}
	_, err = yamlToMapInterface(string(clusterStateBytes))
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Unmarshaling RKE state yaml: %v", err)
	}

	d.Set("cluster_yaml", string(clusterYamlBytes))
	d.Set("rke_state", string(clusterStateBytes))
	d.SetId("")
	err = resourceRKEClusterCreate(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
