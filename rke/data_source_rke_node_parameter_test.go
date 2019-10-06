package rke

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRKENodesConfDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRKENodesConfDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rke_node_parameter.foobar", "address", "192.2.0.1"),
					resource.TestMatchResourceAttr("data.rke_node_parameter.foobar",
						"json",
						regexp.MustCompile(".+")), // should be not empty
					resource.TestMatchResourceAttr("data.rke_node_parameter.foobar",
						"yaml",
						regexp.MustCompile(".+")), // should be not empty
				),
			},
		},
	})
}

const testAccCheckRKENodesConfDataSourceConfig = `
data rke_node_parameter "foobar" {
  address = "192.2.0.1"
  user    = "root"
  role    = ["controlplane", "worker", "etcd"]
}
`
