resource rke_cluster "cluster" {
  nodes = [
    {
      address = "1.2.3.4"
      user    = "ubuntu"
      role    = ["controlplane", "worker", "etcd"]
      ssh_key = "${file("~/.ssh/id_rsa")}"
    },
  ]

  addons = <<EOL
kind: Namespace
apiVersion: v1
metadata:
  name: cattle-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cattle-crb
subjects:
- kind: User
  name: system:serviceaccount:cattle-system:default
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: Service
metadata:
  namespace: cattle-system
  name: cattle-service
  labels:
    app: cattle
spec:
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: http
  - port: 443
    targetPort: 443
    protocol: TCP
    name: https
  selector:
    app: cattle
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  namespace: cattle-system
  name: cattle-ingress-http
spec:
  rules:
  - http:
      paths:
      -   backend:
            serviceName: cattle-service
            servicePort: 80
      -   backend:
            serviceName: cattle-service
            servicePort: 443
---
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  namespace: cattle-system
  name: cattle
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: cattle
    spec:
      containers:
      - image: rancher/server:master
        imagePullPolicy: Always
        name: cattle-server
        ports:
        - containerPort: 80
          protocol: TCP
        - containerPort: 443
          protocol: TCP
EOL

  // addons can be specified by URL or path
  #addons_include = [
  #  "https://github.com/yamamoto-febc/terraform-provider-rke/blob/master/examples/rancher_server/rancher-server.yml",
  #]
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
#resource "local_file" "kube_cluster_yaml" {
#  filename = "${path.root}/kube_config_cluster.yml"
#  content  = "${rke_cluster.cluster.kube_config_yaml}"
#}
