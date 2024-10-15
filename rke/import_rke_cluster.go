package rke

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRKEClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	files, err := splitImportID(d.Id())
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	clusterYamlBytes, err := os.ReadFile(files[0])
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Reading RKE config file %s: %v", files[0], err)
	}
	if len(clusterYamlBytes) == 0 {
		return []*schema.ResourceData{}, fmt.Errorf("RKE config is nil")
	}
	_, err = yamlToMapInterface(string(clusterYamlBytes))
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("unmarshalling RKE config yaml: %v", err)
	}

	clusterStateBytes, err := os.ReadFile(files[1])
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("Reading RKE state file %s: %v", files[0], err)
	}
	if len(clusterStateBytes) == 0 {
		return []*schema.ResourceData{}, fmt.Errorf("RKE state is nil")
	}
	_, err = yamlToMapInterface(string(clusterStateBytes))
	if err != nil {
		return []*schema.ResourceData{}, fmt.Errorf("unmarshalling RKE state yaml: %v", err)
	}

	if len(files) == 3 && files[2] == "dind" {
		d.Set("dind", true)
	}

	d.Set("cluster_yaml", string(clusterYamlBytes))
	d.Set("rke_state", string(clusterStateBytes))
	d.SetId("")
	diag := resourceRKEClusterCreate(ctx, d, meta)
	if diag.HasError() {
		return []*schema.ResourceData{}, errors.New(diag[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
