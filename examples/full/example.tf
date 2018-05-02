# This example same as https://github.com/rancher/rke/blob/master/cluster.yml
resource rke_cluster "cluster" {
  # Disable port check validation between nodes
  disable_port_check = false

  ###############################################
  # Kubernets nodes
  ###############################################
  nodes = [
    {
      address      = "1.1.1.1"
      user         = "ubuntu"
      role         = ["controlplane", "etcd"]
      ssh_key_path = "~/.ssh/id_rsa"
      port         = 2222
    },
    {
      address = "2.2.2.2"
      user    = "ubuntu"
      role    = ["worker"]

      ssh_key = <<EOL
-----BEGIN RSA PRIVATE KEY-----

-----END RSA PRIVATE KEY-----
EOL

      # or
      #ssh_key      = "${file("~/.ssh/id_rsa")}"
    },
    {
      address           = "example.com"
      user              = "ubuntu"
      role              = ["role"]
      hostname_override = "node3"
      internal_address  = "192.168.1.6"

      labels = {
        app = "ingress"
      }
    },
  ]

  ###############################################
  # Kubernetes services
  ###############################################
  services_etcd {
    # if external etcd used
    #path    = "/etcdcluster"
    #ca_cert = "${file("ca_cert")}"
    #cert    = "${file("cert")}"
    #key     = "${file("key")}"
  }

  services_kube_api {
    service_cluster_ip_range = "10.42.0.0/1"
    pod_security_policy      = false

    # add additional arguments to the kubernetes component
    # Note that this WILL OVERRIDE existing defaults
    extra_args = {
      v = 4
    }
  }

  services_kube_controller {
    cluster_cidr             = "10.42.0.0/16"
    service_cluster_ip_range = "10.43.0.0/16"
  }

  services_scheduler {}

  services_kubelet {
    cluster_domain        = "cluster.local"
    cluster_dns_server    = "10.43.0.10"
    infra_container_image = "gcr.io/google_containers/pause-amd64:3.0"

    # Optionally define additional volume binds to a service
    extra_binds = [
      "/host/dev:/dev",
      "usr/libexec/kubernetes/kubelet-plugins:/usr/libexec/kubernetes/kubelet-plugins",
    ]
  }

  services_kubeproxy {}

  #########################################################
  # Network(CNI) - supported: flannel/calico/canal/weave
  #########################################################

  # default: flannel
  network {
    plugin = "flannel"

    options = {
      # To specify flannel interface, you can use the 'flannel_iface' option:  # flannel_iface = "eth1"
    }
  }

  # If you are using calico on AWS or GCE, use the network plugin config option:
  # 'calico_cloud_provider: aws'
  #network {
  #  plugin = "calico"
  #  options = {
  #    calico_cloud_provider = "aws"
  #  }
  #}

  # To specify flannel interface for canal plugin, you can use the 'canal_iface' option:
  #network {
  #  plugin = "canal"
  #  options = {
  #    canal_iface = "eth1"
  #  }
  #}


  ################################################
  # Authentication
  ################################################
  # At the moment, the only authentication strategy supported is x509.
  # You can optionally create additional SANs (hostnames or IPs) to add to
  #  the API server PKI certificate. This is useful if you want to use a load balancer
  #  for the control plane servers, for example.
  authentication {
    strategy = "x509"

    sans = [
      "10.18.160.10",
      "my-loadbalancer-1234567890.us-west-2.elb.amazonaws.com",
    ]
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
  ################################################
  # System Images
  ################################################
  system_images {
    etcd                        = "rancher/etcd:v3.0.17"
    kubernetes                  = "rancher/hyperkube:v1.10.1"
    alpine                      = "rancher/rke-tools:v0.1.4"
    nginx_proxy                 = "rancher/rke-tools:v0.1.4"
    cert_downloader             = "rancher/rke-tools:v0.1.4"
    kubernetes_services_sidecar = "rancher/rke-tools:v0.1.4"
    kube_dns                    = "rancher/k8s-dns-kube-dns-amd64:1.14.5"
    dnsmasq                     = "rancher/k8s-dns-dnsmasq-nanny-amd64:1.14.5"
    kube_dns_sidecar            = "rancher/k8s-dns-sidecar-amd64:1.14.5"
    kube_dns_autoscaler         = "rancher/cluster-proportional-autoscaler-amd64:1.0.0"

    # for CNI
    flannel     = "rancher/coreos-flannel:v0.9.1"
    flannel_cni = "rancher/coreos-flannel-cni:v0.2.0"

    #calico_node                 = ""
    #calico_cni                  = ""
    #calico_controllers          = ""
    #calico_ctl                  = ""
    #canal_node                  = ""
    #canal_cni                   = ""
    #canal_flannel               = ""
    #weave_node                  = ""
    #weave_cni                   = ""

    # pod_infra_container         = ""
    # ingress                     = ""
    # ingress_backend             = ""
    # dashboard                   = ""
    # heapster                    = ""
    # grafana                     = ""
    # influxdb                    = ""
    # tiller                      = ""
  }

  ################################################
  # SSH configuration
  ################################################
  ssh_key_path = "~/.ssh/test"
  ssh_agent_auth = false

  ################################################
  # Authorization
  ################################################
  # Kubernetes authorization mode
  #   - Use `mode: rbac` to enable RBAC
  #   - Use `mode: none` to disable authorization
  authorization {
    mode = "rbac"
  }

  ################################################
  # Versions
  ################################################
  #Enable/Disable strict docker version checking
  ignore_docker_version = false //
  # Kubernetes version to use
  # (if kubernetes image is specifed, image version takes precedence)
  kubernetes_version = "v1.10.1"

  ################################################
  # Private Registries
  ################################################
  # List of registry credentials, if you are using a Docker Hub registry,
  # you can omit the `url` or set it to `docker.io`
  private_registries = {
    url      = "registry1.com"
    user     = "Username"
    password = "password1"
  }
  private_registries = {
    url      = "registry2.com"
    user     = "Username"
    password = "password1"
  }

  ################################################
  # Ingress
  ################################################
  # Currently only nginx ingress provider is supported.
  # To disable ingress controller, set `provider: none`
  ingress {
    provider = "nginx"

    # To enable ingress on specific nodes, use the node_selector, eg:
    #node_selector = {
    #  app = "ingress"
    #}
    #extra_args = {
    #  enable-ssl-passthrough = ""
    #}
  }

  ################################################
  # Cluster Name
  ################################################
  # Cluster Name used in the kube config
  #cluster_name = ""

  ################################################
  # Cloud Provider
  ################################################
  cloud_provider {
    name = "aws"

    # for Azure
    #name = "azure"
    #cloud_config = {
    #  aadClientId       = "xxxxxxxxxxxx"
    #  aadClientSecret   = "xxxxxxxxxxx"
    #  location          = "westus"
    #  resourceGroup     = "rke-rg"
    #  subnetName        = "rke-subnet"
    #  subscriptionId    = "xxxxxxxxxxx"
    #  vnetName          = "rke-vnet"
    #  tenantId          = "xxxxxxxxxx"
    #  securityGroupName = "rke-nsg"
    #}
  }

  # Kubernetes directory path
  # prefix_path = "/"
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################


#resource "local_file" "kube_cluster_yaml" {
#  filename = "${path.root}/kube_config_cluster.yml"
# content  = "${rke_cluster.cluster.kube_config_yaml}"
#}


###############################################################################
# If you ca_crt/client_cert/client_key, please uncomment follows.
###############################################################################


#resource "local_file" "ca_crt" {
#  filename = "${path.root}/ca_cert"
#  content  = "${rke_cluster.cluster.ca_crt}"
#}
#
#resource "local_file" "client_cert" {
#  filename = "${path.root}/client_cert"
#  content  = "${rke_cluster.cluster.client_cert}"
#}
#
#resource "local_file" "client_key" {
#  filename = "${path.root}/client_key"
#  content  = "${rke_cluster.cluster.client_key}"
#}

