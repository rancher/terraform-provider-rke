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
---
kind: Namespace
apiVersion: v1
metadata:
  name: cattle-system
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: cattle-admin
  namespace: cattle-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cattle-crb
  namespace: cattle-system
subjects:
- kind: ServiceAccount
  name: cattle-admin
  namespace: cattle-system
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: Secret
metadata:
  name: cattle-keys-ingress
  namespace: cattle-system
type: Opaque
data:
  tls.crt: <BASE64_CRT>  # ssl cert for ingress. If selfsigned, must be signed by same CA as cattle server
  tls.key: <BASE64_KEY>  # ssl key for ingress. If selfsigned, must be signed by same CA as cattle server
---
apiVersion: v1
kind: Secret
metadata:
  name: cattle-keys-server
  namespace: cattle-system
type: Opaque
data:
  cert.pem: <BASE64_CRT>    # ssl cert for cattle server.
  key.pem: <BASE64_KEY>     # ssl key for cattle server.
  cacerts.pem: <BASE64_CA>  # CA cert used to sign cattle server cert and key
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
  annotations:
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "30"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "1800"   # Max time in seconds for ws to remain shell window open
    nginx.ingress.kubernetes.io/proxy-send-timeout: "1800"   # Max time in seconds for ws to remain shell window open
spec:
  rules:
  - host: <FQDN>  # FQDN to access cattle server
    http:
      paths:
      - backend:
          serviceName: cattle-service
          servicePort: 80
  tls:
  - secretName: cattle-keys-ingress
    hosts:
    - <FQDN>      # FQDN to access cattle server
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
      serviceAccountName: cattle-admin
      containers:
      - image: rancher/rancher:master
        imagePullPolicy: Always
        name: cattle-server
        ports:
        - containerPort: 80
          protocol: TCP
        - containerPort: 443
          protocol: TCP
        volumeMounts:
        - mountPath: /etc/rancher/ssl
          name: cattle-keys-volume
          readOnly: true
      volumes:
      - name: cattle-keys-volume
        secret:
          defaultMode: 420
          secretName: cattle-keys-server
EOL
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
#resource "local_file" "kube_cluster_yaml" {
#  filename = "${path.root}/kube_config_cluster.yml"
#  content  = "${rke_cluster.cluster.kube_config_yaml}"
#}
