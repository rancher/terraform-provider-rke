package rke

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	envRKENodeAddr   = "RKE_NODE_ADDR"
	envRKENodeUser   = "RKE_NODE_USER"
	envRKENodeSSHKey = "RKE_NODE_SSH_KEY"
)

var (
	nodeIP     string
	nodeUser   string
	nodeSSHKey string
)

func TestAccResourceRKECluster(t *testing.T) {
	if ip, ok := os.LookupEnv(envRKENodeAddr); ok {
		nodeIP = ip
	}
	if user, ok := os.LookupEnv(envRKENodeUser); ok {
		nodeUser = user
	}
	if key, ok := os.LookupEnv(envRKENodeSSHKey); ok {
		nodeSSHKey = key
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudIconConfigBasic(nodeIP, nodeUser, nodeSSHKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.#", "1"),
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.0.address", nodeIP),
					resource.TestMatchResourceAttr(
						"rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
				),
			},
			{
				Config: testAccCheckSakuraCloudIconConfigUpdate(nodeIP, nodeUser, nodeSSHKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.#", "1"),
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.0.address", nodeIP),
					resource.TestMatchResourceAttr(
						"rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.0.labels.%", "2"),
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.0.labels.foo", "foo"),
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.0.labels.bar", "bar"),
				),
			},
		},
	})
}

func testAccCheckRKEClusterDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rke_cluster" {
			continue
		}

		masterURL := fmt.Sprintf("https://%s:6443", rs.Primary.ID)
		req, err := http.NewRequest("GET", masterURL, nil)
		if err != nil {
			return nil
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc := &http.Client{Timeout: 2 * time.Second, Transport: tr}
		resp, err := hc.Do(req)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		if err == nil {
			return errors.New("RKE cluster still exists")
		}
	}

	return nil
}

func testAccCheckSakuraCloudIconConfigBasic(ip, user, sshKey string) string {
	return fmt.Sprintf(`	
resource rke_cluster "cluster" {
  nodes = [
    {
      address = "%s"
      user    = "%s"
      role    = ["controlplane", "worker", "etcd"]
      ssh_key = <<EOF
%s
EOF
    },
  ]
}
	`, ip, user, sshKey)

}

func testAccCheckSakuraCloudIconConfigUpdate(ip, user, sshKey string) string {
	return fmt.Sprintf(`	
resource rke_cluster "cluster" {
  nodes = [
    {
      address = "%s"
      user    = "%s"
      role    = ["controlplane", "worker", "etcd"]
      ssh_key = <<EOF
%s
EOF
      labels = {
        foo = "foo"
        bar = "bar"
      }
    },
  ]
}
	`, ip, user, sshKey)

}
