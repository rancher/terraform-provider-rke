package rke

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/pki"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testAccRKEClusterNodes = []string{"tf-testacc1", "tf-testacc2"}
)

func TestAccResourceRKECluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRKEConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.#", "1"),
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.0.address", testAccRKEClusterNodes[0]),
					resource.TestMatchResourceAttr("rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
					testAccCheckTempFilesExists(),
				),
			},
			{
				Config: testAccCheckRKEConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.#", "1"),
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.0.address", testAccRKEClusterNodes[0]),
					resource.TestMatchResourceAttr("rke_cluster.cluster", "kube_config_yaml", regexp.MustCompile(".+")), // should be not empty
					resource.TestMatchResourceAttr("rke_cluster.cluster", "rke_cluster_yaml", regexp.MustCompile(".+")), // should be not empty
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.0.labels.%", "2"),
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.0.labels.foo", "foo"),
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.0.labels.bar", "bar"),
					testAccCheckTempFilesExists(),
				),
			},
			{
				Config: testAccCheckRKEConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.#", "1"),
					testAccCheckRKENodeExists("rke_cluster.cluster", testAccRKEClusterNodes[0]),
					testAccCheckRKEClusterYAML("rke_cluster.cluster", 1),
				),
			},
			{
				Config: testAccCheckRKEConfigBasic2Nodes(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.#", "2"),
					testAccCheckRKENodeExists("rke_cluster.cluster", testAccRKEClusterNodes[0], testAccRKEClusterNodes[1]),
					testAccCheckRKEClusterYAML("rke_cluster.cluster", 2),
				),
			},
			{
				Config: testAccCheckRKEConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rke_cluster.cluster", "nodes.#", "1"),
					testAccCheckRKENodeExists("rke_cluster.cluster", testAccRKEClusterNodes[0]),
					testAccCheckRKEClusterYAML("rke_cluster.cluster", 1),
				),
			},
		},
	})
}

func testAccCheckTempFilesExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if _, err := os.Stat(pki.ClusterConfig); err == nil {
			return fmt.Errorf("temporary file %q is still exists", pki.ClusterConfig)
		}

		kubeClusterYAML := fmt.Sprintf("%s%s", pki.KubeAdminConfigPrefix, pki.ClusterConfig)
		if _, err := os.Stat(kubeClusterYAML); err == nil {
			return fmt.Errorf("temporary file %q is still exists", kubeClusterYAML)
		}

		rkeStatePath := cluster.GetStateFilePath(pki.ClusterConfig, "")
		if _, err := os.Stat(rkeStatePath); err == nil {
			return fmt.Errorf("temporary file %q is still exists", rkeStatePath)
		}

		return nil
	}
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

		masterURL := rs.Primary.Attributes["api_server_url"]
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
		nodes, err := client.CoreV1().Nodes().List(context.Background(), meta_v1.ListOptions{})
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

func testAccCheckRKEClusterYAML(n string, expectNodeLen int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID of rke_cluster is set")
		}

		v := rs.Primary.Attributes["rke_cluster_yaml"]

		var rkeConfig v3.RancherKubernetesEngineConfig
		if err := yaml.Unmarshal([]byte(v), &rkeConfig); err != nil {
			return err
		}

		if len(rkeConfig.Nodes) != expectNodeLen {
			return fmt.Errorf("rke_cluster_yaml has unexpected nodes. expect: %d, actual: %d", expectNodeLen, len(rkeConfig.Nodes))
		}

		return nil
	}
}

func testAccCheckRKEClusterDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rke_cluster" {
			continue
		}

		masterURL := rs.Primary.Attributes["api_server_url"]
		req, err := http.NewRequest("GET", masterURL, nil)
		if err != nil {
			continue
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc := &http.Client{Timeout: 2 * time.Second, Transport: tr}
		resp, err := hc.Do(req)
		if resp != nil {
			resp.Body.Close() // nolint
		}
		if err != nil {
			continue
		}
		return errors.New("RKE cluster still exists")
	}

	// check tmp files
	return testAccCheckTempFilesExists()(s)
}

func testAccCheckRKEConfigBasic() string {
	return fmt.Sprintf(`	
resource rke_cluster "cluster" {
  ignore_docker_version = true
  addon_job_timeout = 60
  dind = true
  dind_dns_server = "8.8.8.8"
  nodes {
    address = "%s"
    user    = "docker"
    role    = ["controlplane", "worker", "etcd"]
  }
  upgrade_strategy {
    drain = true
    max_unavailable_worker = "20%%"
  }
}
`, testAccRKEClusterNodes[0])

}

func testAccCheckRKEConfigUpdate() string {
	return fmt.Sprintf(`	
resource rke_cluster "cluster" {
  ignore_docker_version = true
  addon_job_timeout = 60
  dind = true
  dind_dns_server = "8.8.8.8"
  nodes {
    address = "%s"
    user    = "docker"
    role    = ["controlplane", "worker", "etcd"]
    labels = {
      foo = "foo"
      bar = "bar"
    }
  }
  upgrade_strategy {
    drain = true
    max_unavailable_worker = "20%%"
  }
}
`, testAccRKEClusterNodes[0])

}

func testAccCheckRKEConfigBasic2Nodes() string {
	return fmt.Sprintf(`	
resource rke_cluster "cluster" {
  ignore_docker_version = true
  addon_job_timeout = 60
  dind = true
  dind_dns_server = "8.8.8.8"
  nodes {
    address = "%s"
    user    = "docker"
    role    = ["controlplane", "worker", "etcd"]
  }
  nodes {
    address = "%s"
    user    = "docker"
    role    = ["worker"]
  }
  upgrade_strategy {
    drain = true
    max_unavailable_worker = "20%%"
  }
}
`, testAccRKEClusterNodes[0], testAccRKEClusterNodes[1])
}
