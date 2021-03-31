FROM centos:7.9.2009 as builder

RUN yum install -y gcc gcc-c++ make && yum clean all

RUN set -euxo pipefail; \
    case $(arch) in x86_64) arch=amd64;; aarch64) arch=arm64;; esac; \
    curl -SL https://golang.org/dl/go1.15.7.linux-$arch.tar.gz | tar -xzC /usr/local

ENV PATH $PATH:/usr/local/go/bin

RUN arch=$(go env GOARCH); \
    mkdir -p /usr/local/kubebuilder/bin; \
    curl -SL https://github.com/etcd-io/etcd/releases/download/v3.3.11/etcd-v3.3.11-linux-$arch.tar.gz | tar -xzC /tmp; \
    cp /tmp/etcd-v3.3.11-linux-$arch/etcd /usr/local/kubebuilder/bin/etcd; \
    curl -SL https://dl.k8s.io/v1.16.4/kubernetes-server-linux-$arch.tar.gz | tar -xzC /tmp; \
    cp /tmp/kubernetes/server/bin/kube-apiserver /usr/local/kubebuilder/bin/kube-apiserver; \
    cp /tmp/kubernetes/server/bin/kubectl /usr/local/kubebuilder/bin/kubectl; \
    rm -rf /tmp/kubernetes

WORKDIR /workspace

COPY . .

RUN go build -mod=readonly -o=greeter_server cmd/greeter_server/main.go
RUN go build -mod=readonly -o=greeter_client cmd/greeter_client/main.go
RUN go build -mod=readonly -o=greeter_operator operator/main.go

FROM centos:7.9.2009 as greeter_server

COPY --from=builder /workspace/greeter_server /usr/bin/greeter_server

ENTRYPOINT [ "/usr/bin/greeter_server" ]

FROM centos:7.9.2009 as greeter_client

COPY --from=builder /workspace/greeter_client /usr/bin/greeter_client

ENTRYPOINT [ "/usr/bin/greeter_client" ]

FROM centos:7.9.2009 as greeter_operator

COPY --from=builder /workspace/greeter_operator /usr/bin/greeter_operator

ENTRYPOINT [ "/usr/bin/greeter_operator" ]
