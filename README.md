# Cluster Proxy

[![License](https://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![Go](https://github.com/open-cluster-management-io/cluster-proxy/actions/workflows/go-presubmit.yml/badge.svg)](https://github.com/open-cluster-management-io/cluster-proxy/actions/workflows/go-presubmit.yml)

## Table of Contents

- [What is Cluster Proxy?](#what-is-cluster-proxy)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Verification](#verification)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Code Examples](#code-examples)
- [Performance](#performance)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [References](#references)

## What is Cluster Proxy?

Cluster Proxy is a pluggable addon working on OCM (Open Cluster Management) based on the extensibility
provided by [addon-framework](https://github.com/open-cluster-management-io/addon-framework) 
which automates the installation of [apiserver-network-proxy](https://github.com/kubernetes-sigs/apiserver-network-proxy)
on both hub cluster and managed clusters. 

The network proxy establishes reverse proxy tunnels from the managed cluster to the hub cluster, enabling 
clients from the hub network to access services in the managed clusters' network even when all the clusters 
are isolated in different VPCs or network segments.

### Key Features

- **Secure Tunneling**: Establishes secure reverse proxy tunnels using Kubernetes' konnectivity
- **Multi-Cluster Networking**: Enables cross-cluster communication in isolated network environments
- **OCM Integration**: Seamlessly integrates with Open Cluster Management ecosystem
- **Automatic Management**: Automates proxy server and agent installation and lifecycle management

## Architecture

Cluster Proxy consists of two main components:

- **Addon-Manager**: Manages the installation of proxy-servers (proxy ingress) in the hub cluster
- **Addon-Agent**: Manages the installation of proxy-agents for each managed cluster

The overall architecture is shown below:

![Architecture Diagram](./hack/picture/arch.png)


## Getting Started

### Prerequisites

Before installing Cluster Proxy, ensure you have:

- **OCM Registration** (>= 0.5.0): Open Cluster Management must be installed and configured
- **Kubernetes Cluster**: A running Kubernetes cluster with admin access
- **Helm 3.x**: For installation via Helm charts
- **kubectl**: Configured to access your cluster

### Installation

#### Installing via Helm Chart

1. Adding helm repo:

```shell
$ helm repo add ocm https://openclustermanagement.blob.core.windows.net/releases/
$ helm repo update
$ helm search repo ocm/cluster-proxy
NAME                       	CHART VERSION	APP VERSION	DESCRIPTION                   
ocm/cluster-proxy          	<..>       	    1.0.0      	A Helm chart for Cluster-Proxy
```

2. Install the helm chart:

```shell
$ helm install \
    -n open-cluster-management-addon --create-namespace \
    cluster-proxy ocm/cluster-proxy 
$ kubectl -n open-cluster-management-cluster-proxy get pod
NAME                                           READY   STATUS        RESTARTS   AGE
cluster-proxy-5d8db7ddf4-265tm                 1/1     Running       0          12s
cluster-proxy-addon-manager-778f6d679f-9pndv   1/1     Running       0          33s
...
```

### Verification

3. The addon will be automatically installed to your registered clusters. 
   Verify the addon installation:

```shell
$ kubectl get managedclusteraddon -A | grep cluster-proxy
NAMESPACE         NAME                     AVAILABLE   DEGRADED   PROGRESSING
<your cluster>    cluster-proxy            True                   
```

4. Check that proxy components are running:

```shell
# Check hub cluster components
$ kubectl -n open-cluster-management-addon get pods -l app=cluster-proxy
NAME                                           READY   STATUS    RESTARTS   AGE
cluster-proxy-5d8db7ddf4-265tm                 1/1     Running   0          5m
cluster-proxy-addon-manager-778f6d679f-9pndv   1/1     Running   0          5m

# Check managed cluster agents
$ kubectl -n open-cluster-management-agent-addon get pods -l app=cluster-proxy-agent
NAME                                  READY   STATUS    RESTARTS   AGE
cluster-proxy-agent-xyz123            1/1     Running   0          3m
```

## Usage

### Basic Usage

By default, the proxy servers are running in gRPC mode so the proxy clients 
are expected to proxy through the tunnels using the [konnectivity-client](https://github.com/kubernetes-sigs/apiserver-network-proxy#clients).
Konnectivity is the underlying technique of Kubernetes' [egress-selector](https://kubernetes.io/docs/tasks/extend-kubernetes/setup-konnectivity/)
feature.

### Code Examples

To proxy to a managed cluster, you need to override the dialer of the Kubernetes client config object. 
Here's a basic example:

```go
// Instantiate a gRPC proxy dialer
tunnel, err := konnectivity.CreateSingleUseGrpcTunnel(
    context.TODO(),
    <proxy service>,
    grpc.WithTransportCredentials(grpccredentials.NewTLS(proxyTLSCfg)),
)
cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
if err != nil {
    return err
}
// The managed cluster's name
cfg.Host = clusterName
// Override the default TCP dialer
cfg.Dial = tunnel.DialContext 
```

For more detailed examples, see:
- [Complete test client example](./examples/test-client.md)
- [Accessing exported services](./examples/access-exported-services.md)
- [Accessing Prometheus](./examples/access-prometheus.md)

## Performance

Here's the result of network bandwidth benchmarking via [goben](https://github.com/udhos/goben)
with and without Cluster-Proxy (i.e. Apiserver-Network-Proxy). The proxying through the tunnel 
involves approximately 50% performance overhead, so it's recommended to avoid transferring 
data-intensive traffic over the proxy when possible.

|  Bandwidth  |   Direct   | over Cluster-Proxy |
|-------------|------------|--------------------|
|  Read/Mbps  |  902 Mbps  |     461 Mbps       |
|  Write/Mbps |  889 Mbps  |     428 Mbps       |

## Troubleshooting

### Common Issues

For frequently asked questions and common troubleshooting scenarios, see our [FAQ](./FQA.md).

### Debug Commands

```shell
# Check addon status
kubectl get managedclusteraddon -A | grep cluster-proxy

# Check proxy server logs
kubectl -n open-cluster-management-addon logs -l app=cluster-proxy

# Check proxy agent logs on managed cluster
kubectl -n open-cluster-management-agent-addon logs -l app=cluster-proxy-agent

# Verify proxy configuration
kubectl get managedproxyconfiguration cluster-proxy -o yaml
```

### Getting Help

If you encounter issues:
1. Check the [FAQ](./FQA.md) for common problems and solutions
2. Review the [examples](./examples/) for usage patterns
3. Open an issue in this repository with detailed logs and configuration

## Contributing

We welcome contributions! Please see our [Contributing Guide](./CONTRIBUTING.md) for details on:
- How to submit issues and pull requests
- Development setup and testing
- Code of conduct and community guidelines

## References

- **Design Document**: [OCM Enhancement Proposal](https://github.com/open-cluster-management-io/enhancements/tree/main/enhancements/sig-architecture/14-addon-cluster-proxy)
- **Addon Framework**: [OCM Addon Framework](https://github.com/open-cluster-management-io/addon-framework)
- **Konnectivity**: [Kubernetes Apiserver Network Proxy](https://github.com/kubernetes-sigs/apiserver-network-proxy)
- **Open Cluster Management**: [OCM Community](https://open-cluster-management.io/)
