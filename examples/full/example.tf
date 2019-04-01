# This example same as https://github.com/rancher/rke/blob/master/cluster.yml
resource "rke_cluster" "cluster" {
  # Disable port check validation between nodes
  disable_port_check = false

  ###############################################
  # Kubernets nodes
  ###############################################
  nodes {
    address      = "1.1.1.1"
    user         = "ubuntu"
    role         = ["controlplane", "etcd"]
    ssh_key_path = "~/.ssh/id_rsa"
    port         = 2222
  }

  nodes {
    address = "2.2.2.2"
    user    = "ubuntu"
    role    = ["worker"]

    ssh_key = <<EOL
-----BEGIN RSA PRIVATE KEY-----

-----END RSA PRIVATE KEY-----
EOL

    # or
    #ssh_key      = "${file("~/.ssh/id_rsa")}"
  }

  nodes {
    address = "example.com"
    user    = "ubuntu"

    role = ["controlplane", "etcd", "worker"]
    # or
    # roles             = "controlplane,etcd,worker"

    hostname_override = "node3"
    internal_address  = "192.168.1.6"

    labels = {
      app = "ingress"
    }
  }

  # If set to true, RKE will not fail when unsupported Docker version are found
  ignore_docker_version = false

  ################################################
  # SSH configuration
  ################################################
  # Cluster level SSH private key
  # Used if no ssh information is set for the node
  ssh_key_path = "~/.ssh/test"

  # Enable use of SSH agent to use SSH private keys with passphrase
  # This requires the environment `SSH_AUTH_SOCK` configured pointing to your SSH agent which has the private key added
  ssh_agent_auth = false

  ################################################
  # Bastion/Jump host configuration
  ################################################
  #bastion_host {
  #  address      = "1.1.1.1"
  #  user         = "ubuntu"
  #  ssh_key_path = "~/.ssh/id_rsa"
  #  or
  #  ssh_key      = file("~/.ssh/id_rsa")
  #  port         = 2222
  #}

  ################################################
  # Private Registries
  ################################################
  # List of registry credentials, if you are using a Docker Hub registry,
  # you can omit the `url` or set it to `docker.io`
  private_registries {
    url      = "registry1.com"
    user     = "Username"
    password = "password1"
  }
  private_registries {
    url      = "registry2.com"
    user     = "Username"
    password = "password1"
  }

  ################################################
  # Cluster Name
  ################################################
  # Set the name of the Kubernetes cluster
  #cluster_name = ""

  ################################################
  # Versions
  ################################################
  # The kubernetes version used.
  # For now, this should match the version defined in rancher/types defaults map:
  #    https://github.com/rancher/types/blob/master/apis/management.cattle.io/v3/k8s_defaults.go#L14
  #
  # In case the kubernetes_version and kubernetes image in system_images are defined,
  # the system_images configuration will take precedence over kubernetes_version.
  #
  # Allowed values: [v1.13.5-rancher1-2(default), v1.12.7-rancher1-2, v1.11.9-rancher1-1]
  kubernetes_version = "v1.13.5-rancher1-2"

  ################################################
  # System Images
  ################################################
  # System Image Tags are defaulted to a tag tied with specific kubernetes Versions
  # Default Tags:
  #    https://github.com/rancher/types/blob/master/apis/management.cattle.io/v3/k8s_defaults.go)
  #
  system_images {
    kubernetes                  = "rancher/hyperkube:v1.10.3-rancher2"
    etcd                        = "rancher/coreos-etcd:v3.1.12"
    alpine                      = "rancher/rke-tools:v0.1.9"
    nginx_proxy                 = "rancher/rke-tools:v0.1.9"
    cert_downloader             = "rancher/rke-tools:v0.1.9"
    kubernetes_services_sidecar = "rancher/rke-tools:v0.1.9"
    kube_dns                    = "rancher/k8s-dns-kube-dns-amd64:1.14.8"
    dnsmasq                     = "rancher/k8s-dns-dnsmasq-nanny-amd64:1.14.8"
    kube_dns_sidecar            = "rancher/k8s-dns-sidecar-amd64:1.14.8"
    kube_dns_autoscaler         = "rancher/cluster-proportional-autoscaler-amd64:1.0.0"
    pod_infra_container         = "rancher/pause-amd64:3.1"
  }

  ###############################################
  # Kubernetes services
  ###############################################
  services_etcd {
    # if external etcd used
    #path      = "/etcdcluster"
    #ca_cert   = file("ca_cert")
    #cert      = file("cert")
    #key       = file("key")

    # for etcd snapshots
    #snapshot  = false
    #retention = "24h"
    #creation  = "5m0s"
  }

  services_kube_api {
    # IP range for any services created on Kubernetes
    # This must match the service_cluster_ip_range in kube-controller
    service_cluster_ip_range = "10.43.0.0/1"

    # Expose a different port range for NodePort services
    service_node_port_range = "30000-32767"

    pod_security_policy = false

    # Add additional arguments to the kubernetes API server
    # This WILL OVERRIDE any existing defaults
    extra_args = {
      audit-log-path            = "-"
      delete-collection-workers = 3
      v                         = 4
    }
  }

  services_kube_controller {
    # CIDR pool used to assign IP addresses to pods in the cluster
    cluster_cidr = "10.42.0.0/16"

    # IP range for any services created on Kubernetes
    # This must match the service_cluster_ip_range in kube-api
    service_cluster_ip_range = "10.43.0.0/16"
  }

  services_scheduler {
  }

  services_kubelet {
    # Base domain for the cluster
    cluster_domain = "cluster.local"

    # IP address for the DNS service endpoint
    cluster_dns_server = "10.43.0.10"

    # Fail if swap is on
    fail_swap_on = false

    # Optionally define additional volume binds to a service
    extra_binds = [
      "/usr/libexec/kubernetes/kubelet-plugins:/usr/libexec/kubernetes/kubelet-plugins",
    ]
  }

  services_kubeproxy {
  }

  ################################################
  # Authentication
  ################################################
  # Currently, only authentication strategy supported is x509.
  # You can optionally create additional SANs (hostnames or IPs) to add to the API server PKI certificate.
  # This is useful if you want to use a load balancer for the control plane servers.
  authentication {
    strategy = "x509"

    sans = [
      "10.18.160.10",
      "my-loadbalancer-1234567890.us-west-2.elb.amazonaws.com",
    ]
  }

  ################################################
  # Authorization
  ################################################
  # Kubernetes authorization mode
  #   - Use `mode: "rbac"` to enable RBAC
  #   - Use `mode: "none"` to disable authorization
  authorization {
    mode = "rbac"
  }

  ################################################
  # Cloud Provider
  ################################################
  # If you want to set a Kubernetes cloud provider, you specify the name and configuration
  cloud_provider {
    name = "aws"
  }

  # Add-ons are deployed using kubernetes jobs. RKE will give up on trying to get the job status after this timeout in seconds..
  addon_job_timeout = 30

  #########################################################
  # Network(CNI) - supported: flannel/calico/canal/weave
  #########################################################
  # There are several network plug-ins that work, but we default to canal
  network {
    plugin = "canal"
  }

  ################################################
  # Ingress
  ################################################
  # Currently only nginx ingress provider is supported.
  # To disable ingress controller, set `provider: none`
  ingress {
    provider = "nginx"
  }

  ################################################
  # Addons
  ################################################
  # all addon manifests MUST specify a namespace
  addons = <<EOL
---
apiVersion: v1
kind: Pod
metadata:
  name: my-nginx
  namespace: default
spec:
  containers:
  - name: my-nginx
    image: nginx
    ports:
    - containerPort: 80
EOL


    addons_include = [
      "https://raw.githubusercontent.com/rook/rook/master/cluster/examples/kubernetes/rook-operator.yaml",
      "https://raw.githubusercontent.com/rook/rook/master/cluster/examples/kubernetes/rook-cluster.yaml",
      "/path/to/manifest",
    ]
  }

  ###############################################################################
  # If you need kubeconfig.yml for using kubectl, please uncomment follows.
  ###############################################################################
  #resource "local_file" "kube_cluster_yaml" {
  #  filename = "${path.root}/kube_config_cluster.yml"
  #  content  = rke_cluster.cluster.kube_config_yaml
  #}
  ###############################################################################
  # If you need ca_crt/client_cert/client_key, please uncomment follows.
  ###############################################################################
  #resource "local_file" "ca_crt" {
  #  filename = "${path.root}/ca_cert"
  #  content  = rke_cluster.cluster.ca_crt
  #}
  #
  #resource "local_file" "client_cert" {
  #  filename = "${path.root}/client_cert"
  #  content  = rke_cluster.cluster.client_cert
  #}
  #
  #resource "local_file" "client_key" {
  #  filename = "${path.root}/client_key"
  #  content  = rke_cluster.cluster.client_key
  #}
