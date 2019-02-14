package rke

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				Config: testAccCheckRKEConfigBasic(nodeIP, nodeUser, nodeSSHKey),
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
				Config: testAccCheckRKEConfigUpdate(nodeIP, nodeUser, nodeSSHKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.#", "1"),
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.0.address", nodeIP),
					resource.TestMatchResourceAttr(
						"rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
					resource.TestMatchResourceAttr(
						"rke_cluster.cluster", "rke_cluster_yaml", regexp.MustCompile(".+")), // should be not empty
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

func TestAccResourceRKECluster_NodeCountUpAndDown(t *testing.T) {
	var nodeIPs, nodeUsers, nodeSSHKeys []string

	for i := 0; i < 2; i++ {
		nodeIPEnv := fmt.Sprintf("%s_%d", envRKENodeAddr, i)
		nodeUserEnv := fmt.Sprintf("%s_%d", envRKENodeUser, i)
		nodeSSHKeyEnv := fmt.Sprintf("%s_%d", envRKENodeSSHKey, i)
		if ip, ok := os.LookupEnv(nodeIPEnv); ok {
			nodeIPs = append(nodeIPs, ip)
		}
		if user, ok := os.LookupEnv(nodeUserEnv); ok {
			nodeUsers = append(nodeUsers, user)
		}
		if key, ok := os.LookupEnv(nodeSSHKeyEnv); ok {
			nodeSSHKeys = append(nodeSSHKeys, key)
		}
	}
	requireValues := [][]string{nodeIPs, nodeUsers, nodeSSHKeys}
	for _ ,values := range requireValues {
		if len(values) != 2 {
			t.Skip("Acceptance tests skipped unless required env set")
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckForMultiNodes(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRKEConfigNodeCountUpAndDownSingleNode(
					nodeIPs[0], nodeUsers[0], nodeSSHKeys[0],
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.#", "1"),
					testAccCheckRKENodeExists("rke_cluster.cluster", nodeIPs[0]),
				),
			},
			{
				Config: testAccCheckRKEConfigNodeCountUpAndDownMultiNodes(
					nodeIPs[0], nodeUsers[0], nodeSSHKeys[0],
					nodeIPs[1], nodeUsers[1], nodeSSHKeys[1],
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.#", "2"),
					testAccCheckRKENodeExists("rke_cluster.cluster", nodeIPs[0], nodeIPs[1]),
				),
			},
			{
				Config: testAccCheckRKEConfigNodeCountUpAndDownSingleNode(
					nodeIPs[0], nodeUsers[0], nodeSSHKeys[0],
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"rke_cluster.cluster", "nodes.#", "1"),
					testAccCheckRKENodeExists("rke_cluster.cluster", nodeIPs[0]),
				),
			},
		},
	})
}

func testAccCheckRKENodeExists(n string, nodeIPs ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID of rke_cluster is set")
		}

		masterURL := fmt.Sprintf("https://%s:6443", rs.Primary.ID)
		strKubeConfig := rs.Primary.Attributes["kube_config_yaml"]
		if strKubeConfig == "" {
			return errors.New("kube_config_yaml is empty")
		}

		// save kube_config_yaml to tmpfile
		tmpFile, err := ioutil.TempFile("", "test_acc_rke_cluster_kube_config_")
		if err != nil {
			return errors.New("failed to create temp file")
		}
		defer os.Remove(tmpFile.Name()) // nolint
		if err := ioutil.WriteFile(tmpFile.Name(), []byte(strKubeConfig), 0644); err != nil {
			return errors.New("failed to create temp file")
		}

		// create kubernetes client
		kubeConfig, err := clientcmd.BuildConfigFromFlags(masterURL, tmpFile.Name())
		if err != nil {
			return err
		}
		client, err := kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return err
		}

		// getNodes
		nodes, err := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
		if err != nil {
			return err
		}
		for _, ip := range nodeIPs {
			found := false
			for _, node := range nodes.Items {
				for _, addr := range node.Status.Addresses {
					if ip == addr.Address {
						found = true
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("node %q not found", ip)
			}
		}

		return nil
	}
}

func TestAccResourceRKEClusterWithNodeParameter(t *testing.T) {
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
				Config: testAccCheckRKEConfigWithNodesConfBasic(nodeIP, nodeUser, nodeSSHKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
				),
			},
			{
				Config: testAccCheckRKEConfigWithNodesConfUpdate(nodeIP, nodeUser, nodeSSHKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
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

func testAccCheckRKEConfigBasic(ip, user, sshKey string) string {
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

func testAccCheckRKEConfigUpdate(ip, user, sshKey string) string {
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

func testAccCheckRKEConfigNodeCountUpAndDownSingleNode(ip, user, sshKey string) string {
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
func testAccCheckRKEConfigNodeCountUpAndDownMultiNodes(ip1, user1, sshKey1, ip2, user2, sshKey2 string) string {
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
	`, ip1, user1, sshKey1, ip2, user2, sshKey2)
}

func testAccCheckRKEConfigWithNodesConfBasic(ip, user, sshKey string) string {
	return fmt.Sprintf(`	
data rke_node_parameter "node" {
    address = "%s"
    user    = "%s"
    role    = ["controlplane", "worker", "etcd"]
    ssh_key = <<EOF
%s
EOF
}

resource rke_cluster "cluster" {
  nodes_conf = ["${data.rke_node_parameter.node.json}"]
}
	`, ip, user, sshKey)

}

func testAccCheckRKEConfigWithNodesConfUpdate(ip, user, sshKey string) string {
	return fmt.Sprintf(`	
data rke_node_parameter "node" {
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
}

resource rke_cluster "cluster" {
  nodes_conf = ["${data.rke_node_parameter.node.json}"]
}
	`, ip, user, sshKey)

}
