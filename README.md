# platform-services
A public PoC with functioning services using a simple Istio Mesh running on K8s

Version : <repo-version>0.9.0-master-aaaagqrfhew</repo-version>

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/fetchcore-forks/platform-services/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/fetchcore-forks/platform-services)](https://goreportcard.com/report/fetchcore-forks/platform-services)

This project is intended as a sand-box for experimenting with Istio and some example services in a similiar manner to what is used by the Cognizant Evolutionary AI services.  It also provides a good way of exposing and testing out non-proprietary platform functions while collaborating with other parties such as vendors and customers.

Because this project is intended to mirror production examples of services deploy it requires that the user has an account and a registered internet domain name.  A service such as domains.google.com, or cloudflare is a good start.  Your DNS registered host will be used to issue certificates on your behalf to secure the public connections that are exposed by services to the internet, and more specifically to secure the username and password based access exposed by the service mesh.

# Purpose

This proof of concept (PoC) implementation is intended as a means by which the LEAF team can experiment with features of Service Mesh, PaaS, and SaaS platforms provided by third parties.  This project serves as a way of exercising non cognizant services so that code can be openly shared while testing external services and technologies, and for support in relation to external open source offerings in a public support context.

In its current form the PoC is used to deploy two main services, an experiment service and a downstream service. These services are provisioned with a gRPC API and leverage an Authorizationm Athentication, and Accounting (AAA) capability and an Observability platform integration to services offered by thrid parties.

This project is intended as a sand-box for experimenting with Istio and some of the services we use in our Evolutionary AI services.  It also provides a good way of exposing and testing out non-proprietary platform functions with other parties such as vendors and customers.

# Installation

These instructions were used with Kubernetes 1.19.x, and Istio 1.11.4.

## Development and Building from source

Clone the repository using the following instructions when this is part of a larger project using multiple services:
<pre><code><b>mkdir ~/project
cd ~/project
export GOPATH=`pwd`
export PATH=$GOPATH/bin:$PATH
mkdir -p src/github.com/fetchcore-forks
cd src/github.com/fetchcore-forks
git clone https://github.com/fetchcore-forks/platform-services
cd platform-services
</b></code></pre>

To boostrap development you will need a copy of Go 1.17.3+ available.

Go installation instructions can be found at, https://golang.org/doc/install.

Now download any dependencies, once, into our development environment.

<pre><code><b>
go mod vendor
go get -d github.com/karlmutch/duat/cmd/semver
go get -d github.com/karlmutch/duat/cmd/github-release
go get -d github.com/karlmutch/duat/cmd/stencil
go get -d github.com/karlmutch/petname/cmd/petname
</b></code></pre>

## Running the build using the container

Creating a build container to isolate the build into a versioned environment

<pre><code><b>docker build -t platform-services:latest --build-arg USER=$USER --build-arg USER_ID=`id -u $USER` --build-arg USER_GROUP_ID=`id -g $USER` .
</b></code></pre>

Prior to doing the build a GitHub OAUTH token needs to be defined within your environment.  Use the github admin pages for your account to generate a token, in Travis builds the token is probably already defined by the Travis service.
<pre><code><b>docker run -e GITHUB_TOKEN=$GITHUB_TOKEN -v $GOPATH:/project platform-services ; echo "Done" ; docker container prune -f</b>
</code></pre>

A combined build script is provided 'platform-services/build.sh' to allow all stages of the build including producing docker images to be run together.

# Deployment

The Istio PoC can be deployed within two contexts.  The first uses an Ubuntu server, the second uses Rancher Desktop within a development context for OSX or Windows 11.

## Laptop Server Deployment

Local developer environments for using Kubernetes are a good fit for Rancher Desktop.  Rancher Desktop allows localized docker images to be build on a developers workstation using kim (Kubernetes Image Manager) which will make images immediately available to the local cluster allow for quick compile, build, and run cycles.

Download rancher desktop using https://github.com/rancher-sandbox/rancher-desktop/releases and install based upon your desktop OS of choice.

During installation you will have the option to set the CPU count and Memory available to Rancher Desktop, workable values are a minimum 8Gb of memory and a minimum of of 6 CPUs.  This is set in the Kubernetes Settings in the administration user interface.

Rancher desktop will install the supporting utilities such as helm, kubectl and nerctl.  You can use the Supporting Utilities panel in the administration interface to create symbolic links to tools that you may have already installed.

Rancher Desktop uses QEMU to provide vitualization, and k3s/containerd to provision Kubernetes and will install both during installation.

Docker images are managed on Rancher Desktop using a image builder hosted within the Kubernetes cluster.  Once images are built inside the cluster they are available automatically available both within the cluster for quick redeploys and to your development host.

On installation, Rancher Desktop will start your cluster automatically and will adjust the kubectl context to point at your new custer.

You can now skip to the Istio installation section below.

## Cluster Server Deployment

These deployment instructions are intended for use with the Ubuntu 18.04 LTS distribution.

The following instructions make use of the stencil tool for templating configuration files.

This major section describes two basic alternatives for deployment, AWS kops, and locally hosted KinD (kubernetes in docker).  Other Kubernetes distribution and deployment models will work but are not explicitly described here.

### Verify Docker Version

Docker for Ubuntu can be retrieved from the snap store using the following:
<pre><code><b>sudo snap install docker
docker --version</b>
Docker version 19.03.13, build cd8016b6bc
</code></pre>
You should have a similar or newer version.

snap based docker installation is done to ensure that a platform specific deployment is made using the Ubuntu distribution for portability reasons.

### Arkade Portable DevOps marketplace

arkade provides a means for installation of Kubernetes and related software packages with a single command.  arkade is used throughout these instructions to ease the installation of several utilities and tools.

```
$ curl -sLS https://dl.get-arkade.dev | sudo sh
x86_64
x86_64
Downloading package https://github.com/alexellis/arkade/releases/download/0.7.10/arkade as /tmp/arkade
Download complete.

Running with sufficient permissions to attempt to move arkade to /usr/local/bin
New version of arkade installed to /usr/local/bin
Creating alias 'ark' for 'arkade'.
            _             _
  __ _ _ __| | ____ _  __| | ___
 / _` | '__| |/ / _` |/ _` |/ _ \
| (_| | |  |   < (_| | (_| |  __/
 \__,_|_|  |_|\_\__,_|\__,_|\___|

Get Kubernetes apps the easy way

Version: 0.7.10
Git Commit: 0931a3af1754c0e5fdd4013b666cdac2821a300d
```

arkcade will use an account specific location for executables which you should add to your path.

```
$ export PATH=$PATH:$HOME/.arkade/bin/
```

### Install Kubectl CLI

Installing the kubectl CLI can be done using arkade.

<pre><code><b>ark get kubectl</b>
Downloading kubectl
https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kubectl
38.37 MiB / 38.37 MiB [-------------------------------------------------------------------------------------------------------------------------------------] 100.00%
Tool written to: /home/kmutch/.arkade/bin/kubectl

# Add (kubectl) to your PATH variable
export PATH=$PATH:$HOME/.arkade/bin/

# Test the binary:
/home/kmutch/.arkade/bin/kubectl

# Or install with:
sudo mv /home/kmutch/.arkade/bin/kubectl /usr/local/bin/
</code></pre>

Add kubectl autocompletion to your current shell:

<pre><code><b>source \<(kubectl completion $(basename $SHELL))</b>
</code></pre>

You can verify that kubectl is installed by executing the following command:

<pre><code><b>kubectl version --client</b>
Client Version: version.Info{Major:"1", Minor:"20", GitVersion:"v1.20.0", GitCommit:"af46c47ce925f4c4ad5cc8d1fca46c7b77d13b38", GitTreeState:"clean", BuildDate:"2020-12-09T16:50:1
7Z", GoVersion:"go1.15.6", Compiler:"gc", Platform:"linux/amd64"}
</code></pre>

### Kubernetes

The experimentsrv component comes with an Istio definition file for deployment into AWS, or KinD using Kubernetes (k8s) and Istio.

The deployment definition file can be found at cmd/experimentsrv/experimentsrv.yaml.

Using AWS k8s will use both the kops, and the kubectl tools. You should have an AWS account configured prior to starting deployments, and your environment variables for using the AWS cli should also be done.

### Base cluster installation

This documentation Kubernetes describes several means by which Kubernetes clusters can be installed, choose one however there are many other alternatives also available.

#### Installating Kubernetes in Docker (KinD)

The KinD installation will typically need a registry included in order for cluster images to be pulled.  An AWS registry could be used however this would negate the objective of having a local registry and no dependency on AWS.

KinD provides a means by which a kubernetes cluster can be installed using the Docker Desktop platform, or on linux plain docker.  KinD installation is supported by arkade and installed as follows:.

```
$ ark get kind
Downloading kind
https://github.com/kubernetes-sigs/kind/releases/download/v0.9.0/kind-linux-amd64
7.08 MiB / 7.08 MiB [---------------------------------------------------------------------------------------------------------------------------------------------------] 100.00%
Tool written to: /home/kmutch/.arkade/bin/kind

# Add (kind) to your PATH variable
export PATH=$PATH:$HOME/.arkade/bin/

# Test the binary:
/home/kmutch/.arkade/bin/kind

# Or install with:
sudo mv /home/kmutch/.arkade/bin/kind /usr/local/bin/
```

The KinD cluster and the registry can be installed using a script within the source code repository, 'kind\_install.sh'.  This script will first check to see if there is a image registry already running on your local docker instance and if not will start one, it will then initialize a KinD instance that trusts your local registry.

```
$ /bin/sh ./kind_install.sh
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.19.1) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Thanks for using kind! 😊
configmap/local-registry-hosting created
$ kubectl config set-context  --namespace=default kind-kind
Context "kind-kind" modified.
$ kubectl get nodes
NAME                 STATUS   ROLES    AGE   VERSION
kind-control-plane   Ready    master   68s   v1.19.1
```

After this the next step is to proceed to the certificates installation section.

#### Installing AWS Kubernetes

The current preferred approach to deploying on AWS is to make use of EKS via the eksctl tool.  The kops instructions in this section are provided for older deployments.

##### Using eksctl with auto-scaling

To install eksctl the following should be done.

```
$ curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
$ sudo mv /tmp/eksctl /usr/local/bin
$ eksctl version
```

A basic cluster with auto scaling can be initialized using eksctl and then the addition of the auto-scaler from the Kubernetes project can be used to scale out the project.  The example eks-cluster.yaml file contains the definitions of a cluster within the us-west-2 region, named platform-services.  Before deploying a long lived cluster it is worth while considering cost savings options which are described at the following URL, https://aws.amazon.com/ec2/cost-and-capacity/.

Cluster creation can be performed using the following:

<pre><code><b>export AWS_ACCOUNT=`aws sts get-caller-identity --query Account --output text`
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com
export AWS_ACCESS_KEY=xxx
export AWS_SECRET_ACCESS_KEY=xxx
export AWS_DEFAULT_REGION=xxx
sudo ntpdate ntp.ubuntu.com
export KUBECONFIG=~/.kube/config
export AWS_CLUSTER_NAME=test-eks
eksctl create cluster -f eks-cluster.yaml</b>
2021-04-06 13:39:37 [ℹ]  eksctl version 0.43.0
2021-04-06 13:39:37 [ℹ]  using region us-west-2
2021-04-06 13:39:37 [ℹ]  subnets for us-west-2a - public:192.168.0.0/19 private:192.168.96.0/19
2021-04-06 13:39:37 [ℹ]  subnets for us-west-2b - public:192.168.32.0/19 private:192.168.128.0/19
2021-04-06 13:39:37 [ℹ]  subnets for us-west-2d - public:192.168.64.0/19 private:192.168.160.0/19
2021-04-06 13:39:37 [ℹ]  nodegroup "overhead" will use "ami-0a93391193b512e5d" [AmazonLinux2/1.19]
2021-04-06 13:39:37 [ℹ]  using SSH public key "/home/kmutch/.ssh/id_rsa.pub" as "eksctl-test-eks-nodegroup-overhead-be:07:a0:27:44:d8:27:04:c2:ba:28:fa:8c:47:7f:09"
2021-04-06 13:39:37 [ℹ]  using Kubernetes version 1.19
2021-04-06 13:39:37 [ℹ]  creating EKS cluster "test-eks" in "us-west-2" region with un-managed nodes
2021-04-06 13:39:37 [ℹ]  1 nodegroup (overhead) was included (based on the include/exclude rules)
2021-04-06 13:39:37 [ℹ]  will create a CloudFormation stack for cluster itself and 1 nodegroup stack(s)
2021-04-06 13:39:37 [ℹ]  will create a CloudFormation stack for cluster itself and 0 managed nodegroup stack(s)
2021-04-06 13:39:37 [ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-west-2 --cluster=test-eks'
2021-04-06 13:39:37 [ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "test-eks" in "us-west-2"
2021-04-06 13:39:37 [ℹ]  2 sequential tasks: { create cluster control plane "test-eks", 3 sequential sub-tasks: { 3 sequential sub-tasks: { wait for control plane to
 become ready, tag cluster, update CloudWatch logging configuration }, create addons, create nodegroup "overhead" } }
2021-04-06 13:39:37 [ℹ]  building cluster stack "eksctl-test-eks-cluster"
2021-04-06 13:39:38 [ℹ]  deploying stack "eksctl-test-eks-cluster"
2021-04-06 13:40:08 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
2021-04-06 13:40:38 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
...
2021-04-06 13:52:39 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
2021-04-06 13:53:39 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
2021-04-06 13:53:39 [✔]  tagged EKS cluster (environment=test-eks)
2021-04-06 13:53:40 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:53:57 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:54:14 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:54:33 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:54:34 [✔]  configured CloudWatch logging for cluster "test-eks" in "us-west-2" (enabled types: audit, authenticator, controllerManager & disabled types
: api, scheduler)
2021-04-06 13:54:34 [ℹ]  building nodegroup stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:54:34 [ℹ]  deploying stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:54:34 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:54:50 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
...
2021-04-06 13:57:48 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:58:06 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:58:06 [ℹ]  waiting for the control plane availability...
2021-04-06 13:58:06 [✔]  saved kubeconfig as "/home/kmutch/.kube/config"
2021-04-06 13:58:06 [ℹ]  no tasks
2021-04-06 13:58:06 [✔]  all EKS cluster resources for "test-eks" have been created
2021-04-06 13:58:06 [ℹ]  adding identity "arn:aws:iam::613076437200:role/eksctl-test-eks-nodegroup-overhea-NodeInstanceRole-Q1RVPO36W4VJ" to auth ConfigMap
2021-04-06 13:58:08 [ℹ]  kubectl command should work with "/home/kmutch/.kube/config", try 'kubectl get nodes'
2021-04-06 13:58:08 [✔]  EKS cluster "test-eks" in "us-west-2" region is ready
<b>kubectl get pods --namespace kube-system</b>
NAME                       READY   STATUS    RESTARTS   AGE
coredns-6548845887-h82sj   1/1     Running   0          10m
coredns-6548845887-wz7dm   1/1     Running   0          10m
</code></pre>

Now the auto scaler can be deployed.

<pre><code><b>
kubectl apply -f eks-scaler.yaml</b>
serviceaccount/cluster-autoscaler created
clusterrole.rbac.authorization.k8s.io/cluster-autoscaler created
role.rbac.authorization.k8s.io/cluster-autoscaler created
clusterrolebinding.rbac.authorization.k8s.io/cluster-autoscaler created
rolebinding.rbac.authorization.k8s.io/cluster-autoscaler created
deployment.apps/cluster-autoscaler created
</code></pre>

##### Using kops

If you are using azure or GCP then options such as acs-engine, and skaffold are natively supported by the cloud vendors and written in Go so are readily usable and can be easily customized and maintained and so these are recommended for those cases.

When using AWS the TLS certificates used to secure the connections to your AWS LoadBalancer will require that an ElasticIP is used.  It is recommended that an ElasticIP is allocated for use and then your DNS entries on the domain registra are modified to used the IP as a registered host matching the LetsEncrypt certificate used.  Using an ElasticIP allows the cluster to be regenerated and for the LoadBalancer to be reassociated with the IP whenever the cluster is regenerated.

<pre><code><b>curl -LO https://github.com/kubernetes/kops/releases/download/v1.16.0/kops-linux-amd64
chmod +x kops-linux-amd64
sudo mv kops-linux-amd64 /usr/local/bin/kops

Add kubectl autocompletion to your current shell:

source <(kops completion $(basename $SHELL))
</b></code></pre>

In order to seed your S3 KOPS\_STATE\_STORE version controlled bucket with a cluster definition the following command could be used:

<pre><code><b>export AWS_AVAILABILITY_ZONES="$(aws ec2 describe-availability-zones --query 'AvailabilityZones[].ZoneName' --output text | awk -v OFS="," '$1=$1')"

export S3_BUCKET=kops-platform-$USER
export KOPS_STATE_STORE=s3://$S3_BUCKET
aws s3 mb $KOPS_STATE_STORE
aws s3api put-bucket-versioning --bucket $S3_BUCKET --versioning-configuration Status=Enabled

export CLUSTER_NAME=test-$USER.platform.cluster.k8s.local

kops create cluster --name $CLUSTER_NAME --zones $AWS_AVAILABILITY_ZONES --node-count 1 --node-size=m4.2xlarge --cloud-labels="HostUser=$HOST:$USER"
</b></code></pre>

Optionally use an image from your preferred zone e.g. --image=ami-0def3275.  Also you can modify the AWS machine types, recommended during developer testing using options such as '--master-size=m4.large --node-size=m4.large'.

Starting the cluster can now be done using the following command:

<pre><code><b>kops update cluster $CLUSTER_NAME --yes</b>
I0309 13:48:49.798777    6195 apply_cluster.go:442] Gossip DNS: skipping DNS validation
I0309 13:48:49.961602    6195 executor.go:91] Tasks: 0 done / 81 total; 30 can run
I0309 13:48:50.383671    6195 vfs_castore.go:715] Issuing new certificate: "ca"
I0309 13:48:50.478788    6195 vfs_castore.go:715] Issuing new certificate: "apiserver-aggregator-ca"
I0309 13:48:50.599605    6195 executor.go:91] Tasks: 30 done / 81 total; 26 can run
I0309 13:48:51.013957    6195 vfs_castore.go:715] Issuing new certificate: "kube-controller-manager"
I0309 13:48:51.087447    6195 vfs_castore.go:715] Issuing new certificate: "kube-proxy"
I0309 13:48:51.092714    6195 vfs_castore.go:715] Issuing new certificate: "kubelet"
I0309 13:48:51.118145    6195 vfs_castore.go:715] Issuing new certificate: "apiserver-aggregator"
I0309 13:48:51.133527    6195 vfs_castore.go:715] Issuing new certificate: "kube-scheduler"
I0309 13:48:51.157876    6195 vfs_castore.go:715] Issuing new certificate: "kops"
I0309 13:48:51.167195    6195 vfs_castore.go:715] Issuing new certificate: "apiserver-proxy-client"
I0309 13:48:51.172542    6195 vfs_castore.go:715] Issuing new certificate: "kubecfg"
I0309 13:48:51.179730    6195 vfs_castore.go:715] Issuing new certificate: "kubelet-api"
I0309 13:48:51.431304    6195 executor.go:91] Tasks: 56 done / 81 total; 21 can run
I0309 13:48:51.568136    6195 launchconfiguration.go:334] waiting for IAM instance profile "nodes.test.platform.cluster.k8s.local" to be ready
I0309 13:48:51.576067    6195 launchconfiguration.go:334] waiting for IAM instance profile "masters.test.platform.cluster.k8s.local" to be ready
I0309 13:49:01.973887    6195 executor.go:91] Tasks: 77 done / 81 total; 3 can run
I0309 13:49:02.489343    6195 vfs_castore.go:715] Issuing new certificate: "master"
I0309 13:49:02.775403    6195 executor.go:91] Tasks: 80 done / 81 total; 1 can run
I0309 13:49:03.074583    6195 executor.go:91] Tasks: 81 done / 81 total; 0 can run
I0309 13:49:03.168822    6195 update_cluster.go:279] Exporting kubecfg for cluster
kops has set your kubectl context to test.platform.cluster.k8s.local
        image: {{.duat.awsecr}}/platform-services/{{.duat.module}}:{{.duat.version}}

Cluster is starting.  It should be ready in a few minutes.

Suggestions:
 * validate cluster: kops validate cluster
 * list nodes: kubectl get nodes --show-labels
 * ssh to the master: ssh -i ~/.ssh/id_rsa admin@api.test.platform.cluster.k8s.local
 * the admin user is specific to Debian. If not using Debian please use the appropriate user based on your OS.
 * read about installing addons at: https://github.com/kubernetes/kops/blob/master/docs/addons.md.

<b>
while [ 1 ]; do
    kops validate cluster > /dev/null && break || sleep 10
done;
</b></code></pre>

The initial cluster spinup will take sometime, use kops commands such as 'kops validate cluster' to determine when the cluster is spun up ready for Istio and the platform services.

## Istio 1.11.x

Istio affords a control layer on top of the k8s data plane.  Instructions for deploying Istio are the vanilla instructions that can be found at, https://istio.io/docs/setup/getting-started/#install.  Istio was at one time a Helm based installation but has since moved to using its own methodology, this is the reason we dont use arkade to install it.

<pre><code><b>cd ~
curl -LO https://github.com/istio/istio/releases/download/1.11.4/istio-1.11.4-linux-amd64.tar.gz
tar xzf istio-1.11.4-linux-amd64.tar.gz
export ISTIO_DIR=`pwd`/istio-1.11.4
export PATH=$ISTIO_DIR/bin:$PATH
cd -
istioctl install --set profile=demo -y -f ./istio-config.yaml
</b>
✔ Istio core installed
✔ Istiod installed
✔ Egress gateways installed
✔ Ingress gateways installed
✔ Installation complete
</code></pre>

In order to access you cluster you will need to define some environment variables that will be used later in these instructions:

<pre><code><b>
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export TCP_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="tcp")].nodePort}')
</b></code></pre>

Determining the INGRESS_HOST value is dependent on the deployment technology being used.  Further information can be found in the sections in this document discussing the various Kubernetes distribitions and deployments, additional information can be found at [Determining the Istio Ingress IP and Ports](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/#determining-the-ingress-ip-and-ports).

## Helm 3 Kubernetes package manager

Helm is used by several packages that are deployed using Kubernetes.  Helm can be installed using instructions found at, https://helm.sh/docs/using\_helm/#installing-helm.  For snap based linux distributions the following can be used as a quick-start.

On installations not using Rancher Destop do the following:

<pre><code><b>sudo snap install helm --classic
</b></code></pre>

<pre><code><b>helm repo update
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
</b></code></pre>

## Encryption

This section describes how to secure traffic into the service mesh ingress.  Lets Encrypt is an internet based service for provisioning certificates.  If you are using a locally deployed mesh within KinD for example you will need to use the minica approach.

The domain name choosen here is one I own and is used only for illustrative purposes.  Substitute your own domains as appropriate.

### minica

Minica is a simple CA intended for use in situations where the CA operator also operates each host where a certificate will be used. It automatically generates both a key and a certificate when asked to produce a certificate. It does not offer OCSP or CRL services. Minica is appropriate, for instance, for generating certificates for RPC systems or microservices.

More information about minica can be found at, https://github.com/jsha/minica.

<pre><code>
$ go get -d github.com/jsha/minica
$ mkdir minica
$ cd minica
$ minica --domains platform-services.karlmutch.com
$ cd -
$ tree minica
minica
├── minica-key.pem
├── minica.pem
└── platform-services.karlmutch.com
    ├── cert.pem
    └── key.pem

1 directory, 4 filess
</code></pre>

### Lets Encrypt

letsencrypt is a public SSL/TLS certificate provider intended for use with the open internet that is being used to secure our service mesh for this project. The lets encrypt provisioning tool can be installed from github and accessed to produce TLS certificates for your service.  If you are deploying on KinD or a local installation then lets encrypt is not supported and using minica, https://github.com/jsha/minica, is recommended instead for testing purposes.

Prior to running the lets encrypt tools you should identify the desired DNS hostname and email you wish to use for your example service cluster.  In our example we have a domain registered as, karlmutch.com.  This domain is available to us as an administrator, and we have choosen to use the host name platform-service.karlmutch.com as the services hostname.

The first step is to add a registered hosts entry for the platform-services.karlmutch.com host into the DNS account, if the host is unknown add an IP address such as 127.0.0.1.  During the generation process you will be prompted to add a DNS TXT record into the custom resource records for the domain, this requires the dummy entry to be present.

Setting up and initiating this process can be done using the following:

<pre><code><b>git clone https://github.com/letsencrypt/letsencrypt</b>
Cloning into 'letsencrypt'...
remote: Enumerating objects: 255, done.
remote: Counting objects: 100% (255/255), done.
remote: Compressing objects: 100% (188/188), done.
remote: Total 71278 (delta 135), reused 110 (delta 65), pack-reused 71023
Receiving objects: 100% (71278/71278), 23.55 MiB | 26.53 MiB/s, done.
Resolving deltas: 100% (52331/52331), done.
<b>cd letsencrypt</b>
<b>./letsencrypt-auto certonly --rsa-key-size 4096 --agree-tos --manual --preferred-challenges=dns --email=karlmutch@cognizant.com -d platform-services.karlmutch.com</b>
</code></pre>

You will be prompted with the IP address logging when starting the script, you should choose 'Y' to enabled the logging as this assists auditing of DNS changes on the internet by registras and regulatory bodies.

<pre><code>
Are you OK with your IP being logged?
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: <b>Y
</b></code></pre>

After this step you will be asked to add a text record to your DNS records proving that you have control over the domain, this will appear much like the following:

<pre><code>
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Please deploy a DNS TXT record under the name
acme-challenge.platform-services.karlmutch.com with the following value:

mbUa9_gb4RVYhTumHy3zIi3PIXFh0k_oOgCie4NvhqQ

Before continuing, verify the record is deployed.
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Press Enter to Continue
</code></pre>

You should wait 10 to 20 minutes for the TXT record to appear in the database of your DNS provider before selecting continue otherwise the verification will fail and you will need to restart it.

<pre><code>
Waiting for verification...
Cleaning up challenges

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/platform-services.karlmutch.com/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/platform-services.karlmutch.com/privkey.pem
   Your cert will expire on 2020-02-23. To obtain a new or tweaked
   version of this certificate in the future, simply run
   letsencrypt-auto again. To non-interactively renew *all* of your
   certificates, run "letsencrypt-auto renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le
</code></pre>

Once the certificate generation is complete you will have the certificate saved at the location detailed in the 'IMPORTANT NOTES' section.  Keep a record of where this is you will need it later.

After the certificate has been issue feel free to delete the TXT record that served as proof of ownership as it is no longer needed.

## Configuration of secrets and cluster access

This project makes use of several secrets that are used to access resources under its control, including the Postgres Database, the Honeycomb service, and the lets encrypt issues certificate.

The experiment service Honeycomb observability solution uses a key to access Datasets defined by the Honeycomb account and store events in the same.  Configuring the service is done by creating a Kubernetes secret.  For now we can define the Honeycomb Dataset, and API key using an environment variable and when we deploy the secrets for the Postgres Database the secret for the API will be injected using the stencil tool.

<pre><code><b>export O11Y_KEY=a54d762df847474b22915
export O11Y_DATASET=platform-services
</b></code></pre>

The services also use a postgres Database instance to persist experiment data, this is installed later in the process.  The following shows an example of what should be defined for Postgres support prior to running the stencil command:

<pre><code><b>export PGRELEASE=$USER-poc
export PGNAMESPACE=postgresql-poc
export PGHOST=$PGRELEASE-postgresql.$PGNAMESPACE.svc.cluster.local
export PORT=5432
export PGUSER=postgres
export PGPASSWORD=p355w0rd
export PGDATABASE=platform
</b></code></pre>

<pre><code><b>
stencil < cmd/experimentsrv/secret.yaml | kubectl apply -f -
</b></code></pre>

The last set of secrets that need to be stored are related to securing the mesh for third party connections using TLS.  This secret contains the full certificate chain and private key needed to implement TLS on the gRPC connections exposed by the mesh.

If you used Lets Encrypt then the following applies:

<pre><code><b>
sudo kubectl create -n istio-system secret generic platform-services-tls-cert \
    --from-file=key=/etc/letsencrypt/live/platform-services.karlmutch.com/privkey.pem \
    --from-file=cert=/etc/letsencrypt/live/platform-services.karlmutch.com/fullchain.pem
</b></code></pre>

If minica was used then the following would be used:

<pre><code><b>
kubectl create -n istio-system secret generic platform-services-tls-cert \
    --from-file=key=./minica/platform-services.karlmutch.com/key.pem \
    --from-file=cert=./minica/platform-services.karlmutch.com/cert.pem
</b></code></pre>

If you are using AWS ACM to manage your certificates the platform-services-tls-cert secret is not required.

## Deploying the Istio Ingress configuration

Istio provides an ingress resource that can be used to secure the service using either secrets (certificates) or using cloud provider provisioned certificates.

Using your own secrets you will use the default ingress yaml that will point at the platform-services-tls-cert Kubernetes provisioned in the previous section.

<pre><code><b>
kubectl apply -f ingress.yaml
</b></code></pre>

### Deploying using AWS Route 53 and ACM Public CA

Create a DNS hosted zone for your cluster.  Add an A record for your clusters ingress that points at your cluster using the load balancer address identified via the following commands:

<pre><code><b>
kubectl get svc istio-ingressgateway -n istio-system</b>
NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP                                                               PORT(S)
                                       AGE
istio-ingressgateway   LoadBalancer   10.100.90.52   acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com   15021:31220/TCP,80:31317/TCP,443:31536/TCP,31400:31243/TCP,15443:30506/TCP   38s
<b>nslookup</b>
> <b>acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com</b>
Server:         127.0.0.53
Address:        127.0.0.53#53

Non-authoritative answer:
Name:   acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com
Address: 35.163.132.217
Name:   acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com
Address: 54.70.157.123
> <b>[ENTER]</b>
<b>cat <<EOF >changes.json
{
            "Comment": "Add LB to the Route53 records so that certificates can be validated against the host name ",
            "Changes": [{
            "Action": "UPSERT",
                        "ResourceRecordSet": {
                                    "Name": "platform-services.karlmutch.com",
                                    "Type": "A",
                                    "TTL": 300,
                                 "ResourceRecords": [{ "Value": "35.163.132.217"}, {"Value": "54.70.157.123"}]
}}]
}
EOF
aws route53 list-hosted-zones-by-name --dns-name karlmutch.com</b>

{
    "HostedZones": [
        {
            "Id": "/hostedzone/Z1UZMEEDVXVHH3",
            "Name": "karlmutch.com.",
            "CallerReference": "RISWorkflow-RD:97f5f182-4b86-41f1-9a24-46218ab70d25",
            "Config": {
                "Comment": "HostedZone created by Route53 Registrar",
                "PrivateZone": false
            },
            "ResourceRecordSetCount": 16
        }
    ],
    "DNSName": "karlmutch.com",
    "IsTruncated": false,
    "MaxItems": "100"
}
<b>aws route53 change-resource-record-sets --hosted-zone-id Z1UZMEEDVXVHH3  --change-batch file://changes.json</b>
{
    "ChangeInfo": {
        "Id": "/change/C029593935XMLAA0T5OHH",
        "Status": "PENDING",
        "SubmittedAt": "2021-04-07T22:12:17.619Z",
        "Comment": "Add LB to the Route53 records so that certificates can be validated against the host name "
    }
}
<b>aws route53  get-change --id /change/C029593935XMLAA0T5OHH</b>
{
    "ChangeInfo": {
        "Id": "/change/C029593935XMLAA0T5OHH",
        "Status": "INSYNC",
        "SubmittedAt": "2021-04-07T22:12:17.619Z",
        "Comment": "Add LB to the Route53 records so that certificates can be validated against the host name "
    }
}
</code></pre>

Now we are in a position to generate the certificate:
<pre><code>
<b>aws acm request-certificate --domain-name platform-services.karlmutch.com --validation-method DNS --idempotency-token 1234</b>
{
    "CertificateArn": "arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb"
}
{
    "Certificate": {
        "CertificateArn": "arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb",
        "DomainName": "platform-services.karlmutch.com",
        "SubjectAlternativeNames": [
            "platform-services.karlmutch.com"
        ],
        "DomainValidationOptions": [
            {
                "DomainName": "platform-services.karlmutch.com",
                "ValidationDomain": "platform-services.karlmutch.com",
                "ValidationStatus": "PENDING_VALIDATION",
                "ResourceRecord": {
                    "Name": "_a9d4b51d79d2b08121a0796cbfbb7a68.platform-services.karlmutch.com.",
                    "Type": "CNAME",
                    "Value": "_e327e9f51160630a9f0056fd3eb56a74.bbfvkzsszw.acm-validations.aws."
                },
                "ValidationMethod": "DNS"
            }
        ],
        "Subject": "CN=platform-services.karlmutch.com",
        "Issuer": "Amazon",
        "CreatedAt": 1617837918.0,
        "Status": "PENDING_VALIDATION",
        "KeyAlgorithm": "RSA-2048",
        "SignatureAlgorithm": "SHA256WITHRSA",
        "InUseBy": [],
        "Type": "AMAZON_ISSUED",
        "KeyUsages": [],
        "ExtendedKeyUsages": [],
        "RenewalEligibility": "INELIGIBLE",
        "Options": {
            "CertificateTransparencyLoggingPreference": "ENABLED"
        }
    }
}
<b>cat <<EOF >changes.json
{
            "Comment": "Add the certificate issuance validation to the Route53 records so that certificates can be validated",
            "Changes": [{
            "Action": "UPSERT",
                        "ResourceRecordSet": {
                                 "Name": "_a9d4b51d79d2b08121a0796cbfbb7a68.platform-services.karlmutch.com.",
                                 "Type": "CNAME",
                                 "TTL": 300,
                                 "ResourceRecords": [{ "Value": "_e327e9f51160630a9f0056fd3eb56a74.bbfvkzsszw.acm-validations.aws."}]
}}]
}
EOF
aws route53 change-resource-record-sets --hosted-zone-id Z1UZMEEDVXVHH3  --change-batch file://changes.json</b>
{
    "ChangeInfo": {
        "Id": "/change/C084369336HI4IZ12CDU9",
        "Status": "PENDING",
        "SubmittedAt": "2021-04-07T23:46:04.831Z",
        "Comment": "Add the certificate issuance validation to the Route53 records so that certificates can be validated"
    }
}
<b>aws route53  get-change --id /change/C084369336HI4IZ12CDU9</b>
{
    "ChangeInfo": {
        "Id": "/change/C084369336HI4IZ12CDU9",
        "Status": "INSYNC",
        "SubmittedAt": "2021-04-07T23:46:04.831Z",
        "Comment": "Add the certificate issuance validation to the Route53 records so that certificates can be validated"
    }
}
# Now we wait for the certificate to be issued:
<b>aws acm describe-certificate --certificate-arn arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb<b>
{
    "Certificate": {
        "CertificateArn": "arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb",
        "DomainName": "platform-services.karlmutch.com",
        "SubjectAlternativeNames": [
            "platform-services.karlmutch.com"
        ],
        "DomainValidationOptions": [
            {
                "DomainName": "platform-services.karlmutch.com",
                "ValidationDomain": "platform-services.karlmutch.com",
                "ValidationStatus": "SUCCESS",
                "ResourceRecord": {
                    "Name": "_a9d4b51d79d2b08121a0796cbfbb7a68.platform-services.karlmutch.com.",
                    "Type": "CNAME",
                    "Value": "_e327e9f51160630a9f0056fd3eb56a74.bbfvkzsszw.acm-validations.aws."
                },
                "ValidationMethod": "DNS"
            }
        ],
        "Serial": "07:df:33:6b:78:11:e2:a3:ee:f1:54:51:3f:81:78:28",
        "Subject": "CN=platform-services.karlmutch.com",
        "Issuer": "Amazon",
        "CreatedAt": 1617837918.0,
        "IssuedAt": 1617839861.0,
        "Status": "ISSUED",
        "NotBefore": 1617753600.0,
        "NotAfter": 1651881599.0,
        "KeyAlgorithm": "RSA-2048",
        "SignatureAlgorithm": "SHA256WITHRSA",
        "InUseBy": [],
        "Type": "AMAZON_ISSUED",
        "KeyUsages": [
            {
                "Name": "DIGITAL_SIGNATURE"
            },
            {
                "Name": "KEY_ENCIPHERMENT"
            }
        ],
        "ExtendedKeyUsages": [
            {
                "Name": "TLS_WEB_SERVER_AUTHENTICATION",
                "OID": "1.3.6.1.5.5.7.3.1"
            },
            {
                "Name": "TLS_WEB_CLIENT_AUTHENTICATION",
                "OID": "1.3.6.1.5.5.7.3.2"
            }
        ],
        "RenewalEligibility": "INELIGIBLE",
        "Options": {
            "CertificateTransparencyLoggingPreference": "ENABLED"
        }
    }
}
</code></pre>

If you are using AWS to provision certificates to secure the ingress connections then use aws-ingress.yaml

<pre><code><b>
kubectl apply -f aws-ingress.yaml
</b></code></pre>

Then we patch the ingress so that it uses our newly issued ACM certificate via the AWS ARN.

<pre><code>
<b>arn="arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb"</b>
<b>kubectl -n istio-system patch service istio-ingressgateway --patch "$(cat<<EOF
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: $arn
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: "https"
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "3600"
EOF
)"</b>
</code></pre>


A test of the certificate on an empty ingress will then appear as follows:

<pre><code>
<b>curl -Iv https://platform-services.karlmutch.com</b>
* Rebuilt URL to: https://platform-services.karlmutch.com/
*   Trying 35.163.132.217...
* TCP_NODELAY set
* Connected to platform-services.karlmutch.com (35.163.132.217) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256
* ALPN, server did not agree to a protocol
* Server certificate:
*  subject: CN=platform-services.karlmutch.com
*  start date: Apr  7 00:00:00 2021 GMT
*  expire date: May  6 23:59:59 2022 GMT
*  subjectAltName: host "platform-services.karlmutch.com" matched cert's "platform-services.karlmutch.com"
*  issuer: C=US; O=Amazon; OU=Server CA 1B; CN=Amazon
*  SSL certificate verify ok.
> HEAD / HTTP/1.1
> Host: platform-services.karlmutch.com
> User-Agent: curl/7.58.0
> Accept: */*
>
* TLSv1.2 (IN), TLS alert, Client hello (1):
* Empty reply from server
* Connection #0 to host platform-services.karlmutch.com left intact
curl: (52) Empty reply from server
</code></pre>

## Deploying the Observability proxy server

This proxy server is used to forward tracing and metrics from your istio mesh based deployment to the Honeycomb service.

<pre><code><b>
helm repo rm honeycomb || true
helm repo add honeycomb https://honeycombio.github.io/helm-charts
# The collector is used to aggregate logs, metrics, events etc from application components inside the system
helm install opentelemetry-collector honeycomb/opentelemetry-collector --set honeycomb.apiKey=$O11Y_KEY --set honeycomb.dataset=$O11Y_DATASET
# Install the Honeycomb Kubernetes Agent 
helm install honeycomb honeycomb/honeycomb --set honeycomb.apiKey=$O11Y_KEY --values honeycomb-agent-values.yaml
# An alternative is to use, kubectl apply -f <(stencil < honeycomb-agent.yaml), https://docs.honeycomb.io/getting-data-in/integrations/kubernetes/#getting-started-using-kubectl
</b></code></pre>

In order to instrument the base Kubernetes deployment for use with honeycomb you should follow the instructions found at https://docs.honeycomb.io/getting-data-in/integrations/kubernetes/.

The dataset used by the istio and services deployed within this project also needs configuration to allow the Honeycomb platform to identify import fields.  Once data begins flowing into the data set you can navigate to the definitions section for the dataset and set the 'Name' item to the name field, 'Parent span ID' item to parentId, 'Service name' to serviceName, 'Span duration' to durationMs, 'Span ID' to id, and finally 'Trace ID' to traceId.

### Postgres DB

To deploy the platform experiment service a database must be present.  The PoC is intended to use an in-cluster DB designed that is dropped after the service is destroyed.

If you wish to use Aurora then you will need to use the AWS CLI Tools or the Web Console to create your Postgres Database appropriately, and then to set your environment variables PGHOST, PGPORT, PGUSER, PGPASSWORD, and PGDATABASE appropriately.  You will also be expected to run the sql setup scripts yourself.

If you are using OSX to deploy into Rancher Desktop installing the postgres client can be done as followings

<pre><code><b>
brew install libpq
export PATH="/usr/local/opt/libpq/bin:$PATH"
PGHOST=127.0.0.1 PGDATABASE=platform psql -f sql/platform.sql -d postgres
</b></code></pre>

The first step is to install the Ubuntu postgres 11 client on your system and then to populate the schema on the remote database:

<pre><code><b></b>
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
sudo apt-get -y update
sudo apt-get -y --upgrade postgresql-client-11
</code></pre>

### Deploying an in-cluster Database

This section gives guidence on how to install an in-cluster database for use-cases where data persistence beyond a single deployment is not a concern.  These instructions are therefore limited to testing only scenarios.  For information concerning Kubernetes storage strategies you should consult other sources and read about stateful sets in Kubernetes.  In production using a single source of truth then cloud provider offerings such as AWS Aurora are recommended.

A secrets file containing host information, passwords and other secrets is assumed to have already been applied using the instructions several sections above.  The secrets are needed to allows access to the postgres DB, and/or other external resources.  YAML files will be needed to populate secrets into the service mesh, individual services document the secrets they require within their README.md files found on github and provide examples, for example https://github.com/fetchcore-forks/platform-services/cmd/experimentsrv/README.md.

In order to deploy Postgres this document describes a helm based approach.  The bitnami postgresql distribution can be installed using the following:

<pre><code><b>
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install $PGRELEASE bitnami/postgresql \
--namespace $PGNAMESPACE --create-namespace --set postgresqlPassword=$PGPASSWORD,postgresqlDatabase=postgres
</b>
NAME: karlmutch-poc
LAST DEPLOYED: Tue Oct 26 13:55:08 2021
NAMESPACE: postgresql-poc
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: postgresql
CHART VERSION: 10.12.7
APP VERSION: 11.13.0

** Please be patient while the chart is being deployed **

PostgreSQL can be accessed via port 5432 on the following DNS names from within your cluster:

    karlmutch-poc-postgresql.postgresql-poc.svc.cluster.local - Read/Write connection

To get the password for "postgres" run:

    export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql-poc karlmutch-poc-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)

To connect to your database run the following command:

    kubectl run karlmutch-poc-postgresql-client --rm --tty -i --restart='Never' --namespace postgresql-poc --image docker.io/bitnami/postgresql:11.13.0-debian-10-r67 --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host karlmutch-poc-postgresql -U postgres -d postgres -p 5432



To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace postgresql-poc svc/karlmutch-poc-postgresql 5432:5432 &
    PGPASSWORD="$POSTGRES_PASSWORD" psql --host 127.0.0.1 -U postgres -d postgres -p 5432
</code></pre>


Special note should be taken of the output from the helm command it has a lot of valuable information concerning your postgres deployment that will be needed when you load the database schema.

Setting up the proxy will be needed prior to running the SQL database provisioning scripts.  When doing this prior to running the postgres client set the PGHOST environment variable to 127.0.0.1 so that the proxy on the localhost is used.  The proxy will timeout after inactivity and shutdown so be prepared to restart it when needed.

<pre><code><b>
kubectl wait --for=condition=Ready --namespace $PGNAMESPACE pod/$PGRELEASE-postgresql-0
kubectl port-forward --namespace $PGNAMESPACE svc/$PGRELEASE-postgresql 5432:5432 &amp;
sleep 2
</b></code></pre>

If you are using OSX then you will need to add the postgres client using the following :

<pre><code><b>
brew install libpq
export PATH="/usr/local/opt/libpq/bin:$PATH"
</b></code></pre>

<pre><code><b>
PGHOST=127.0.0.1 PGDATABASE=platform psql -f sql/platform.sql -d postgres
</b></code></pre>

Further information about how to deployed the service specific database for the experiment service for example can be found in the cmd/experiment/README.md file.

## Deploying into the Istio mesh

This section describes the activities for deployment of the services provided by the Plaform PoC.  The first two sections provide a description of how to deploy the TLS secured ingress for the service mesh, the first being for cloud provisioned systems and the second for the localized KinD based deployment.

### Configuring the Rancher Desktop ingress (Critical Step)
Rancher Desktop will often provision gateway connections to your local host interfaces from the lima based virtual machine via a bridge. To determine if this has occurred you can use the netstat command to look for ports that are in a LISTEN state.

<pre><code><b>
kubectl get svc istio-ingressgateway -n istio-system</b>
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE
istio-ingressgateway   LoadBalancer   10.43.175.205   &lt;pending&gt;     15021:32219/TCP,80:31130/TCP,<b>443:31217/TCP</b>,31400:31228/TCP,15443:30344/TCP   3d23h
# At this point you can see the secured 443 port being mapped to 31217.  Now look at netstat ...
<b>netstat -a | grep LISTEN | grep 31217</b>
tcp4       0      0  \*.31217                \*.\*                    LISTEN
# Indeed there is a LISTEN port on your IPv4 interfaces, meaning that 127.0.0.1 can be used to access services in your cluster.
<b>
export INGRESS_HOST=127.0.0.1
export CLUSTER_INGRESS=$INGRESS_HOST:$SECURE_INGRESS_PORT
</b></code></pre>

### Configuring the KinD ingress (Critical Step)

<pre><code><b>
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export CLUSTER_INGRESS=$INGRESS_HOST:$SECURE_INGRESS_PORT
</b></code></pre>

### Configuring the service cloud based DNS (Critical Step)

When using this mesh instance with a TLS based deployment the DNS domain name used for the LetsEncrypt certificate (CN), will need to have its address record (A) updated to point at the AWS load balancer assigned to the Kubernetes cluster.  In AWS this is done via Route53:

<pre><code><b>
INGRESS_HOST=`kubectl get svc --namespace istio-system -o go-template='{{range .items}}{{range .status.loadBalancer.ingress}}{{.hostname}}{{printf "\n"}}{{end}}{{end}}'`
export CLUSTER_INGRESS=$INGRESS_HOST:$SECURE_INGRESS_PORT
dig +short $INGRESS_HOST
</b></code></pre>

Take the IP addresses from the above output and use these as the A record for the LetsEncrypt host name, platform-services.karlmutch.com, and this will enable accessing the mesh and validation of the common name (CN) in the certificate.  After several minutes, you should test this step by using the following command to verify that the certificate negotiation is completed.

### Service deployment overview

Platform services use Dockerfiles to encapsulate their build steps which are documented within their respective README.md files.  Building services are single step CLI operations and require only the installation of Docker, and any version of Go 1.7 or later.  Builds will produce containers and will upload these to your current AWS account users ECS docker registry.  Deployments are staged from this registry.  

<pre><code><b>kubectl get nodes</b>
NAME                                           STATUS    ROLES     AGE       VERSION
ip-172-20-118-127.us-west-2.compute.internal   Ready     node      17m       v1.9.3
ip-172-20-41-63.us-west-2.compute.internal     Ready     node      17m       v1.9.3
ip-172-20-55-189.us-west-2.compute.internal    Ready     master    18m       v1.9.3
</code></pre>

Once secrets are loaded individual services can be deployed from a checked out developer copy of the service repo using a command like the following :

<pre><code><b>cd ~/project/src/github.com/fetchcore-forks/platform-services</b>
<b>cd cmd/[service] ; kubectl apply -f \<(istioctl kube-inject -f <(stencil [service].yaml 2>/dev/null)); cd - </b>
</code></pre>

The next step is to determine the deployment locations for your images as a result of running the build script.

AWS - The stencil command will automatically fill in the AWS ECR account for you based on your AWS credentials

KinD, MicroK8s - The local host and port number for the internal registry is used.  Export this value as follows:

```
export IMAGE_REPO="localhost:5000/"
```

Rancher Desktop with nerdctl - leave the IMAGE_REPO environment variable undefined.

A minimal set of servers for testing includes the downstream and experiment servers as follows:

```
kubectl apply -f <(istioctl kube-inject -f <(cd cmd/downstream > /dev/null ; stencil < downstream.yaml)) && \
kubectl apply -f <(istioctl kube-inject -f <(cd cmd/experimentsrv > /dev/null ; stencil < experimentsrv.yaml))
```

In order to locate the image repository the stencil tool will test for the presence of AWS credentials and if found will use the account as the source of AWS ECR images.  In the case where the credentials are not present then the default localhost registry will be used for image deployment.

Once the application is deployed you can discover the gateway points within the kubernetes cluster by using the kubectl commands as documented in the cmd/experimentsrv/README.md file.

More information about deploying a real service and using the experimentsrv server can be found at, https://github.com/fetchcore-forks/platform-services/blob/master/cmd/experimentsrv/README.md.

For more information about using Rancher Desktop with nerdctl image wrangling please see the following, https://itnext.io/rancher-desktop-and-nerdctl-for-local-k8s-dev-d1348629932a.

### Testing connectivity using KinD, and Rancher Desktop

Once the initial services have been deployed the connectivity can be tested using the following:

<pre><code><b>curl -Iv --cacert minica/platform-services.karlmutch.com/../minica.pem --header "Host: platform-services.karlmutch.com"  --connect-to "platform-services.karlmutch.com:$SECURE_INGRESS_PORT:$CLUSTER_INGRESS" https://platform-services.karlmutch.com:$SECURE_INGRESS_PORT
</b>
* Rebuilt URL to: https://platform-services.karlmutch.com:30017/
* Connecting to hostname: 172.20.0.2
* Connecting to port: 30017
*   Trying 172.20.0.2...
* TCP\_NODELAY set
* Connected to 172.20.0.2 (172.20.0.2) port 30017 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: minica/platform-services.karlmutch.com/../minica.pem
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Unknown (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Client hello (1):
* TLSv1.3 (OUT), TLS Unknown, Certificate Status (22):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS\_AES\_256\_GCM\_SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: CN=platform-services.karlmutch.com
*  start date: Dec 23 18:10:54 2020 GMT
*  expire date: Jan 22 18:10:54 2023 GMT
*  subjectAltName: host "platform-services.karlmutch.com" matched cert's "platform-services.karlmutch.com"
*  issuer: CN=minica root ca 34d475
*  SSL certificate verify ok.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* Using Stream ID: 1 (easy handle 0x5584d85765c0)
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
\> HEAD / HTTP/2
\> Host: platform-services.karlmutch.com
\> User-Agent: curl/7.58.0
\> Accept: */*
\>
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS Unknown, Unknown (23):
* Connection state changed (MAX\_CONCURRENT\_STREAMS updated)!
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (IN), TLS Unknown, Unknown (23):
\< HTTP/2 404
HTTP/2 404
\< date: Mon, 04 Jan 2021 19:50:31 GMT
date: Mon, 04 Jan 2021 19:50:31 GMT
\< server: istio-envoy
server: istio-envoy

\<
* Connection #0 to host 172.20.0.2 left intact
</code></pre>

### Testing connectivity using Cloud based solutions

A very basic test of the TLS negotiation can be done using the curl command:

<pre><code><b>
curl -Iv https://helloworld.letsencrypt.org
curl -Iv https://platform-services.karlmutch.com:$SECURE_INGRESS_PORT
export INGRESS_HOST=platform-services.karlmutch.com:$SECURE_INGRESS_PORT
</b></code></pre>

### Debugging

There are several pages of debugging instructions that can be used for situations when grpc failures to occur without much context, this applies to unexplained GRPC errors that reference C++ files within envoy etc.  These pages can be found by using the search function on the Istio web site at, https://istio.io/search.html?q=debugging.

You might find the following use cases useful for avoiding using hard coded pod names etc when debugging.

The following example shows enabling debugging for http2 and rbac layers within the Ingress Envoy instance.

<pre><code><b>
kubectl exec $(kubectl get pods -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -n istio-system -- curl -X POST "localhost:15000/logging?rbac=debug" -s
kubectl exec $(kubectl get pods -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -n istio-system -- curl -X POST "localhost:15000/logging?filter=debug" -s
kubectl exec $(kubectl get pods -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -n istio-system -- curl -X POST "localhost:15000/logging?http2=debug" -s
</b></code></pre>

After making a test request the log can be retrieved using something like the following:

<pre><code><b>
kubectl logs $(kubectl get pods --namespace istio-system -l istio=ingressgateway -o jsonpath='{.items[0].metadata.name}') --namespace istio-system</b></code></pre>

When debugging the istio proxy side cars for services you can do the following to enable all of the modules within the proxy:

<pre><code><b>
kubectl exec $(kubectl get pods -l app=experiment -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -- curl -X POST "localhost:15000/logging?level=debug" -s</b></code></pre>

And then the logs can be captured during the testing using the following:

<pre><code><b>
kubectl logs $(kubectl get pods -l app=experiment -o jsonpath='{.items[0].metadata.name}') -c istio-proxy</b></code></pre>

# Logging and Observability

Currently the service mesh is deployed with Observability tools.  These instruction do not go into Observability at this time.  However we do address logging.

Individual services do offering logging using the systemd facilities and these logs are routed to Kubernetes.  Logs can be obtained from pods and containers. The 'kubectl get services' command can be used to identify the running platform services and the 'kubectl get pod' command can be used to get the health of services.  Once a pod isidentified with a running service instance the logs can be extract using a combination of the pod instance and the service name together, for example:

<pre><code><b>kuebctl get services</b>
NAME          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
experiments   ClusterIP   100.68.93.48   <none>        30001/TCP   12m
kubernetes    ClusterIP   100.64.0.1     <none>        443/TCP     1h
<b>kubectl get pod</b>
NAME                             READY     STATUS    RESTARTS   AGE
experiments-v1-bc46b5d68-tltg9   2/2       Running   0          12m
<b>kubectl logs experiments-v1-bc46b5d68-tltg9 experiments</b>
./experimentsrv built at 2018-01-18\_15:22:47+0000, against commit id 34e761994b895ac48cd832ac3048854a671256b0
2018-01-18T16:50:18+0000 INF experimentsrv git hash version 34e761994b895ac48cd832ac3048854a671256b0 \_: [host experiments-v1-bc46b5d68-tltg9]
2018-01-18T16:50:18+0000 INF experimentsrv database startup / recovery has been performed dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com:5432 name platform \_: [host experiments-v1-bc46b5d68-tltg9]
2018-01-18T16:50:18+0000 INF experimentsrv database has 1 connections  dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com:5432 name platform dbConnectionCount 1 \_: [host experiments-v1-bc46b5d68-tltg9]
2018-01-18T16:50:33+0000 INF experimentsrv database has 1 connections  dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com:5432 name platform dbConnectionCount 1 \_: [host experiments-v1-bc46b5d68-tltg9]
</code></pre>

The container name can also include the istio mesh and kubernetes installed system containers for indepth debugging purposes.

### Kubernetes Web UI and console  (Deprecated)

In addition to the kops information for a cluster being hosted on S3, the kubectl information for accessing the cluster is stored within the ~/.kube directory.  The web UI can be deployed using the instruction at https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/#deploying-the-dashboard-ui, the following set of instructions include the deployment as it stood at k8s 1.9.  Take the opportunity to also review the document at the above location.

Kubectl service accounts can be created at will and given access to cluster resources.  To create, authorize and then authenticate a service account the following steps can be used:

```
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/influxdb/influxdb.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/influxdb/heapster.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/influxdb/grafana.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/rbac/heapster-rbac.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml
kubectl create serviceaccount studioadmin
secret_name=`kubectl get serviceaccounts studioadmin -o json | jq '.secrets[] | [.name] | join("")' -r`
secret_kube=`kubectl get secret $secret_name -o json | jq '.data.token' -r | base64 --decode`
# The following will open up all service accounts for admin, review the k8s documentation specific to your
# install version of k8s to narrow the roles
kubectl create clusterrolebinding serviceaccounts-cluster-admin --clusterrole=cluster-admin --group=system:serviceaccounts
```

The value in secret kube can be used to login to the k8s web UI.  First start 'kube proxy' in a terminal window to create a proxy server for the cluster.  Use a browser to navigate to http://localhost:8001/ui.  Then use the value in the secret\_kube variable as your 'Token' (Service Account Bearer Token).

You will now have access to the Web UI for your cluster with full privs.


# AAA using Auth0

Platform services are secured using the Auth0.com service.  Auth0 is a service that provides support for headless machine to machine authentication.  Auth0 is being used initially to provide Bearer tokens for both headless and CLI clients to platform services proof of concept.

Auth0 supports the ability to create a hosted database for storing user account and credential information.  You should navigate to the Connections -> Database section and create a database with the name of "Username-Password-Authentication".  This will be used later when creating applications as your source of user information.

Auth0 authorizations can be done using a Demo auth0.com account.  To do this you will need to add a custom API to the Auth0 account, call it something like "Experiments API" and give it an identifier of "http://api.karlmutch.com/experimentsrv", you should also enable RBAC and the "Add Permissions in the Access Token" options.  Then use the save button to persist the new API.  Identifiers are used as the 'audience' setting when generating tokens via web calls against the AAA features of the Auth0 platform.

The next stop is to use the menu bar to select the "Permissions" tab.  This tab allows you to create a scope to be used for the permissions granted to user.  Create a scope called "all:experiments" with a description, and select the Add button.  This scope will become available for use by authenticated user roles to allow them to access the API.

Next, click the "Machine To Machine Applications" tab.  This should show that a new Test Application has been created and authorized against this API.  To the right of the Authorized switch is a drop down button that can be used to expose more detailed information related to the scopes that are permitted via this API.  You should see that the all:experiments scope is not yet selected, select it and then use the update button.

Now navigate using the left side panel to the Applications screen.  Click to select your new "Experiments (Test Application)".  The screen displayed as a result will show the Client ID, and the Client secret that will be needed later on, take a note of thes values as they will be needed during AAA operation.  Go to the bottom of the page and you will be able to expose some advanced settings".  Inside the advanced settings you will see a ribbon bar with a "Grant Types" tab that can be clicked on revealing the selections available for grant type, ensure that the "password" radio button is selected to enable passwords for authentication, and click on the Save Changes button to save the selection.

The first API added by the Auth0 platform will be the client that accesses the Auth0 service itself providing per user authentication and token generation. When you begin creating a client you will be able to select the "Auth0 Management API" as one of the APIs you wish to secure.

The next step is to create a User and assign the user a role.  The left hand panel has a "Users Management" menu.  Using the menu you can select the "User" option and then use the "CREATE USER" button on the right side of the screen.  This where the real power of the Auth0 platform comes into play as you can use your real email address and perform operations related to identity management and passwords resets without needing to implement these features yourself.  When creating the user the connection field should be filled in with the Database connection that you created initially in these instructions. "Username-Password-Authentication".  After creating your User you can go to the Users panel and click on the email, then click on the permissions tab.  Add the all:experiments permission to the users prodile using the "ASSIGN PERMISSIONS" button.

Having done all of this you will need to go to the account settings and set the default provider for your account to be the Username-Password-Authentication connection.

You can now use various commands to manipulate the APIs outside of what will exist in the application code, this is a distinct advantage over directly using enterprise tools such as Okta.  Should you wish to use Okta as an Identity provider, or backend, to Auth0 then this can be done however you will need help from our Tech Ops department to do this and is an expensive option.  At this time the user and passwords being used for securing APIs can be managed through the Auth0 dashboard including the ability to invite users to become admins.

<pre><code><b>
export AUTH0_DOMAIN=karlmutch-ai.auth0.com
export AUTH0_CLIENT_ID=pL3iSUmOB7EPiXae4gPfuEasccV7PATs
export AUTH0_CLIENT_SECRET=KHSCFuFumudWGKISCYD79ZkwF2YFCiQYurhjik0x6OKYyOb7TkfGKJrHKXXADzqG
export AUTH0_REQUEST=$(printf '{"client_id": "%s", "client_secret": "%s", "audience":"https://api.karlmutch.com/","grant_type":"password", "username": "karlmutch@gmail.com", "password": "Ap9ss2345f", "scope": "all:experiments", "realm": "Username-Password-Authentication" }' "$AUTH0_CLIENT_ID" "$AUTH0_CLIENT_SECRET")
export AUTH0_TOKEN=$(curl -s --request POST --url https://karlmutch.auth0.com/oauth/token --header 'content-type: application/json' --data "$AUTH0_REQUEST" | jq -r '"\(.access_token)"')

</b>
c.f. https://auth0.com/docs/quickstart/backend/golang/02-using#obtaining-an-access-token-for-testing.
</code></pre>

If you are using the test API and you are either running a kubectl port-forward or have a local instance of the postgres DB, you can do something like:

<pre><code><b>kubectl port-forward --namespace $PGNAMESPACE svc/$PGRELEASE-postgresql 5432:5432 &
cd cmd/downstream
go run . --ip-port=":30008" &
cd ../..
cd cmd/experimentsrv
export AUTH0_DOMAIN=karlmutch.auth0.com
export AUTH0_CLIENT_ID=pL3iSUmOB7EPiXae4gPfuEasccV7PATs
export AUTH0_CLIENT_SECRET=KHSCFuFumudWGKISCYD79ZkwF2YFCiQYurhjik0x6OKYyOb7TkfGKJrHKXXADzqG
export AUTH0_REQUEST=$(printf '{"client_id": "%s", "client_secret": "%s", "audience":"https://api.karlmutch.com/","grant_type":"password", "username": "karlmutch@gmail.com", "password": "Passw0rd!", "scope": "all:experiments", "realm": "Username-Password-Authentication" }' "$AUTH0_CLIENT_ID" "$AUTH0_CLIENT_SECRET")
export AUTH0_TOKEN=$(curl -s --request POST --url https://karlmutch.auth0.com/oauth/token --header 'content-type: application/json' --data "$AUTH0_REQUEST" | jq -r '"\(.access_token)"')
go test -v --dbaddr=localhost:5432 -ip-port="[::]:30007" -dbname=platform -downstream="[::]:30008"
</b></code></pre>

## Auth0 claims extensibility

Auth0 can be configured to include additional headers with user metadata such as email addresses etc using custom rules in the Auth0 rules configuration.  Header that are added can be queried and extracted from gRPC HTTP authorization header meta data as shown in the experimentsrv server.go file. An example of a rule is as follows:

<pre><code>
function (user, context, callback) {
  context.accessToken["http://karlmutch.com/user"] = user.email;
  callback(null, user, context);
 }</code></pre>

 An example of extracting this item on the gRPC client side can be found in cmd/experimentsrv/server.go in the GetUserFromClaims function.

# Manually invoking and using production services with TLS

When using the gRPC services within a secured cluster these instructions can be used to access and exercise the services.

An example client for running a simple ping style test against the cluster is provided in the cmd/cli-experiment directory.  This client acts as a test for the presence of the service.  If the commands to obtain a JWT have been followed then this command can be run against the cluster as follows:

<pre><code><b>
cd cmd/cli-experiment
go run . --server-addr=platform-services.karlmutch.com:443 --auth0-token="$AUTH0_TOKEN"</b>
(\*com_karlmtch_ai_experiment.CheckResponse)(0xc00003c280)(modules:"downstream" )
</code></pre>

Once valid transactions are being performed you should go back to the section on Honeycomb and configure the relevant fields inside the definitions panel for your Dataset.

If You are using rancher desktop with a minica implementation then you would use something such as the following:

<pre><code><b>
cd cmd/cli-experiment
kubectl get svc istio-ingressgateway -n istio-system</b>
NAME                   TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)                                                                      AGE
istio-ingressgateway   LoadBalancer   10.43.28.86   <pending>     15021:31027/TCP,80:32104/TCP,<b>443:30552/TCP</b>,31400:30652/TCP,15443:32029/TCP   2d1h
# The following command is using a server-name value that satisfies the IP SAN requirements of the SSL generated by the minica tool.
<b>go run . --server-addr=127.0.0.1:30552 --ca-cert ../../minica/minica.pem --server-name=platform-services.karlmutch.com --auth0-token="$AUTH0_TOKEN"</b>
</code></pre>

# Manually invoking and using services without TLS

When using the gRPC services within an unsecured cluster these instructions can be used to access and exercise the services.

A pre-requiste of manually invoking GRPC servcies is that the grpc_cli tooling is installed.  The instructions for doing this can be found within the grpc source code repository at, https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md.

The following instructions identify a $INGRESS_HOST value for cases where a LoadBalancer is being used.  If you are using Kubernetes in Docker and the cluster is hosted locally then the INGRESS_HOST value should be 127.0.0.1 for the following instructions.

Services used within the platform require that not only is the link integrity and security is maintained using mTLS but that an authorization block is also supplied to verify the user requesting a service.  The authorization can be supplied when using the gRPC command line tool using the metadata options.  First we retrieve a token using curl and then make a request against the service, run in this case as a local docker container, as follows:

<pre><code><b>grpc_cli call localhost:30001 com.karlmutch.experiment.Service.Get "id: 'test'" --metadata authorization:"Bearer $AUTH0_TOKEN"
</b></code></pre>

The services used within the platform all support reflection when using gRPC.  To examine calls available for a server you should first identify the endpoint through which the gateway is being routed, in this case as part of an Istio cluster on AWS with TLS enabled, for example:

<pre><code><b>export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export CLUSTER_INGRESS=$INGRESS_HOST:$SECURE_INGRESS_PORT
<b>grpc_cli ls $CLUSTER_INGRESS -l --channel_creds_type=ssl</b>
filename: experimentsrv.proto
package: com.karlmutch.experiment;
service Service {
  rpc Create(com.karlmutch.experiment.CreateRequest) returns (com.karlmutch.experiment.CreateResponse) {}
  rpc Get(com.karlmutch.experiment.GetRequest) returns (com.karlmutch.experiment.GetResponse) {}
  rpc MeshCheck(com.karlmutch.experiment.CheckRequest) returns (com.karlmutch.experiment.CheckResponse) {}
}

filename: grpc/health/v1/health.proto
package: grpc.health.v1;
service Health {
  rpc Check(grpc.health.v1.HealthCheckRequest) returns (grpc.health.v1.HealthCheckResponse) {}
  rpc Watch(grpc.health.v1.HealthCheckRequest) returns (stream grpc.health.v1.HealthCheckResponse) {}
}

filename: grpc_reflection_v1alpha/reflection.proto
package: grpc.reflection.v1alpha;
service ServerReflection {
  rpc ServerReflectionInfo(stream grpc.reflection.v1alpha.ServerReflectionRequest) returns (stream grpc.reflection.v1alpha.ServerReflectionResponse) {}
}
</code></pre>


In circumstances where you have the CN name validation enabled then the host name should reflect the common name for the host, for example:

<pre><code><b>grpc_cli call platform-services.karlmutch.com:$SECURE_INGRESS_PORT com.karlmutch.experiment.Service.Get "id: 'test'" --metadata authorization:"Bearer $AUTH0_TOKEN" --channel_creds_type=ssl
</b></code></pre>


To drill further into interfaces and examine the types being used within calls you can perform commands such as:

<pre><code><b>grpc_cli type $CLUSTER_INGRESS com.karlmutch.experiment.CreateRequest -l --channel_creds_type=ssl</b>
message CreateRequest {
.com.karlmutch.experiment.Experiment experiment = 1[json_name = "experiment"];
}
<b>grpc_cli type $CLUSTER_INGRESS com.karlmutch.experiment.Experiment -l --channel_creds_type=ssl</b>
message Experiment {
string uid = 1[json_name = "uid"];
string name = 2[json_name = "name"];
string description = 3[json_name = "description"];
.google.protobuf.Timestamp created = 4[json_name = "created"];
map&lt;uint32, .com.karlmutch.experiment.InputLayer&gt; inputLayers = 5[json_name = "inputLayers"];
map&lt;uint32, .com.karlmutch.experiment.OutputLayer&gt; outputLayers = 6[json_name = "outputLayers"];
}
<b>grpc_cli type $CLUSTER_INGRESS com.karlmutch.experiment.InputLayer -l --channel_creds_type=ssl</b>
message InputLayer {
enum Type {
	Unknown = 0;
	Enumeration = 1;
	Time = 2;
	Raw = 3;
}
string name = 1[json_name = "name"];
.com.karlmutch.experiment.InputLayer.Type type = 2[json_name = "type"];
repeated string values = 3[json_name = "values"];
}
</code></pre>

# Notes of using Kiali

A quick a dirty means of using Kiali is shown below.

```
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/addons/prometheus.yaml
kubectl apply -f ${ISTIO_HOME}/samples/addons/kiali.yaml
kubectl port-forward svc/kiali 20001:20001 -n istio-system &
```

Then browse to http://localhost:20001/...

# Shutting down a service, or cluster

<pre><code><b>kubectl delete -f experimentsrv.yaml
</b></code></pre>

<pre><code><b>kops delete cluster $CLUSTER_NAME --yes
</b></code></pre>
