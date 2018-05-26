FROM golang:1.10 as builder
LABEL maintainer="Kazumichi Yamamoto <yamamoto.febc@gmail.com>"
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN  apt-get update && apt-get -y install bash git make zip && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

RUN go get -u github.com/motemen/gobump/cmd/gobump

ADD . /go/src/github.com/yamamoto-febc/terraform-provider-rke
WORKDIR /go/src/github.com/yamamoto-febc/terraform-provider-rke
RUN make build
###

FROM hashicorp/terraform:0.11.7
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>
LABEL MAINTAINER 'Kazumichi Yamamoto <yamamoto.febc@gmail.com>'

RUN set -x && apk add --no-cache --update ca-certificates
RUN mkdir -p /root/.terraform.d/plugins
COPY --from=builder /go/src/github.com/yamamoto-febc/terraform-provider-rke/bin/* /root/.terraform.d/plugins/
