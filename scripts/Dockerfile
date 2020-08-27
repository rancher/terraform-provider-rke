FROM golang:1.12.3

RUN  apt-get update && \
     apt-get install -y xz-utils zip rsync jq curl ca-certificates && \
     curl -fsSL https://get.docker.com | sh - && \
     apt-get clean && \
     rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* && \
     go get -u golang.org/x/lint/golint && \
     curl -L https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl -o /usr/bin/kubectl && chmod 755 /usr/bin/kubectl && \
     curl -LO https://releases.hashicorp.com/terraform/0.12.20/terraform_0.12.20_linux_amd64.zip && unzip terraform_0.12.20_linux_amd64.zip && \
     mv terraform /usr/bin/terraform && chmod 755 /usr/bin/terraform && rm terraform_0.12.20_linux_amd64.zip 
VOLUME /go/src/github.com/rancher/terraform-provider-rke
