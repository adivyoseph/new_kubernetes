/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubeadm

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconfigv1beta1 "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig/v1beta1"
	kubeproxyconfigv1alpha1 "k8s.io/kubernetes/pkg/proxy/apis/kubeproxyconfig/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MasterConfiguration contains a list of elements which make up master's
// configuration object.
type MasterConfiguration struct {
	metav1.TypeMeta

	// API holds configuration for the k8s apiserver.
	API API
	// KubeProxy holds configuration for the k8s service proxy.
	KubeProxy KubeProxy
	// Etcd holds configuration for etcd.
	Etcd Etcd
	// KubeletConfiguration holds configuration for the kubelet.
	KubeletConfiguration KubeletConfiguration
	// Networking holds configuration for the networking topology of the cluster.
	Networking Networking
	// KubernetesVersion is the target version of the control plane.
	KubernetesVersion string
	// NodeName is the name of the node that will host the k8s control plane.
	// Defaults to the hostname if not provided.
	NodeName string
	// NoTaintMaster will, if set, suppress the tainting of the
	// master node allowing workloads to be run on it (e.g. in
	// single node configurations).
	NoTaintMaster bool

	// Token is used for establishing bidirectional trust between nodes and masters.
	// Used for joining nodes in the cluster.
	Token string
	// TokenTTL defines the ttl for Token. Defaults to 24h.
	TokenTTL *metav1.Duration
	// TokenUsages describes the ways in which this token can be used.
	TokenUsages []string
	// Extra groups that this token will authenticate as when used for authentication
	TokenGroups []string

	// CRISocket is used to retrieve container runtime info.
	CRISocket string

	// APIServerExtraArgs is a set of extra flags to pass to the API Server or override
	// default ones in form of <flagname>=<value>.
	// TODO: This is temporary and ideally we would like to switch all components to
	// use ComponentConfig + ConfigMaps.
	APIServerExtraArgs map[string]string
	// ControllerManagerExtraArgs is a set of extra flags to pass to the Controller Manager
	// or override default ones in form of <flagname>=<value>
	// TODO: This is temporary and ideally we would like to switch all components to
	// use ComponentConfig + ConfigMaps.
	ControllerManagerExtraArgs map[string]string
	// SchedulerExtraArgs is a set of extra flags to pass to the Scheduler or override
	// default ones in form of <flagname>=<value>
	// TODO: This is temporary and ideally we would like to switch all components to
	// use ComponentConfig + ConfigMaps.
	SchedulerExtraArgs map[string]string

	// APIServerExtraVolumes is an extra set of host volumes mounted to the API server.
	APIServerExtraVolumes []HostPathMount
	// ControllerManagerExtraVolumes is an extra set of host volumes mounted to the
	// Controller Manager.
	ControllerManagerExtraVolumes []HostPathMount
	// SchedulerExtraVolumes is an extra set of host volumes mounted to the scheduler.
	SchedulerExtraVolumes []HostPathMount

	// APIServerCertSANs sets extra Subject Alternative Names for the API Server
	// signing cert.
	APIServerCertSANs []string
	// CertificatesDir specifies where to store or look for all required certificates.
	CertificatesDir string

	// ImageRepository is the container registry to pull control plane images from.
	ImageRepository string

	// CIImageRepository is the container registry for core images generated by CI.
	// Useful for running kubeadm with images from CI builds.
	// +k8s:conversion-gen=false
	CIImageRepository string

	// UnifiedControlPlaneImage specifies if a specific container image should be
	// used for all control plane components.
	UnifiedControlPlaneImage string

	// AuditPolicyConfiguration defines the options for the api server audit system.
	AuditPolicyConfiguration AuditPolicyConfiguration

	// FeatureGates enabled by the user.
	FeatureGates map[string]bool

	// The cluster name
	ClusterName string
}

// API struct contains elements of API server address.
type API struct {
	// AdvertiseAddress sets the IP address for the API server to advertise.
	AdvertiseAddress string
	// ControlPlaneEndpoint sets a stable IP address or DNS name for the control plane; it
	// can be a valid IP address or a RFC-1123 DNS subdomain, both with optional TCP port.
	// In case the ControlPlaneEndpoint is not specified, the AdvertiseAddress + BindPort
	// are used; in case the ControlPlaneEndpoint is specified but without a TCP port,
	// the BindPort is used.
	// Possible usages are:
	// e.g. In an cluster with more than one control plane instances, this field should be
	// assigned the address of the external load balancer in front of the
	// control plane instances.
	// e.g.  in environments with enforced node recycling, the ControlPlaneEndpoint
	// could be used for assigning a stable DNS to the control plane.
	ControlPlaneEndpoint string
	// BindPort sets the secure port for the API Server to bind to.
	// Defaults to 6443.
	BindPort int32
}

// TokenDiscovery contains elements needed for token discovery.
type TokenDiscovery struct {
	// ID is the first part of a bootstrap token. Considered public information.
	// It is used when referring to a token without leaking the secret part.
	ID string
	// Secret is the second part of a bootstrap token. Should only be shared
	// with trusted parties.
	Secret string
	// TODO: Seems unused. Remove?
	// Addresses []string
}

// Networking contains elements describing cluster's networking configuration.
type Networking struct {
	// ServiceSubnet is the subnet used by k8s services. Defaults to "10.96.0.0/12".
	ServiceSubnet string
	// PodSubnet is the subnet used by pods.
	PodSubnet string
	// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
	DNSDomain string
}

// Etcd contains elements describing Etcd configuration.
type Etcd struct {
	// Endpoints of etcd members. Useful for using external etcd.
	// If not provided, kubeadm will run etcd in a static pod.
	Endpoints []string
	// CAFile is an SSL Certificate Authority file used to secure etcd communication.
	CAFile string
	// CertFile is an SSL certification file used to secure etcd communication.
	CertFile string
	// KeyFile is an SSL key file used to secure etcd communication.
	KeyFile string
	// DataDir is the directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string
	// ExtraArgs are extra arguments provided to the etcd binary
	// when run inside a static pod.
	ExtraArgs map[string]string
	// Image specifies which container image to use for running etcd.
	// If empty, automatically populated by kubeadm using the image
	// repository and default etcd version.
	Image string
	// ServerCertSANs sets extra Subject Alternative Names for the etcd server
	// signing cert. This is currently used for the etcd static-pod.
	ServerCertSANs []string
	// PeerCertSANs sets extra Subject Alternative Names for the etcd peer
	// signing cert. This is currently used for the etcd static-pod.
	PeerCertSANs []string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeConfiguration contains elements describing a particular node.
// TODO: This struct should be replaced by dynamic kubelet configuration.
type NodeConfiguration struct {
	metav1.TypeMeta

	// CACertPath is the path to the SSL certificate authority used to
	// secure comunications between node and master.
	// Defaults to "/etc/kubernetes/pki/ca.crt".
	CACertPath string
	// DiscoveryFile is a file or url to a kubeconfig file from which to
	// load cluster information.
	DiscoveryFile string
	// DiscoveryToken is a token used to validate cluster information
	// fetched from the master.
	DiscoveryToken string
	// DiscoveryTokenAPIServers is a set of IPs to API servers from which info
	// will be fetched. Currently we only pay attention to one API server but
	// hope to support >1 in the future.
	DiscoveryTokenAPIServers []string
	// DiscoveryTimeout modifies the discovery timeout
	DiscoveryTimeout *metav1.Duration
	// NodeName is the name of the node to join the cluster. Defaults
	// to the name of the host.
	NodeName string
	// TLSBootstrapToken is a token used for TLS bootstrapping.
	// Defaults to Token.
	TLSBootstrapToken string
	// Token is used for both discovery and TLS bootstrapping.
	Token string
	// CRISocket is used to retrieve container runtime info.
	CRISocket string
	// The cluster name
	ClusterName string

	// DiscoveryTokenCACertHashes specifies a set of public key pins to verify
	// when token-based discovery is used. The root CA found during discovery
	// must match one of these values. Specifying an empty set disables root CA
	// pinning, which can be unsafe. Each hash is specified as "<type>:<value>",
	// where the only currently supported type is "sha256". This is a hex-encoded
	// SHA-256 hash of the Subject Public Key Info (SPKI) object in DER-encoded
	// ASN.1. These hashes can be calculated using, for example, OpenSSL:
	// openssl x509 -pubkey -in ca.crt openssl rsa -pubin -outform der 2>&/dev/null | openssl dgst -sha256 -hex
	DiscoveryTokenCACertHashes []string

	// DiscoveryTokenUnsafeSkipCAVerification allows token-based discovery
	// without CA verification via DiscoveryTokenCACertHashes. This can weaken
	// the security of kubeadm since other nodes can impersonate the master.
	DiscoveryTokenUnsafeSkipCAVerification bool

	// FeatureGates enabled by the user.
	FeatureGates map[string]bool
}

// KubeletConfiguration contains elements describing initial remote configuration of kubelet.
type KubeletConfiguration struct {
	BaseConfig *kubeletconfigv1beta1.KubeletConfiguration
}

// GetControlPlaneImageRepository returns name of image repository
// for control plane images (API,Controller Manager,Scheduler and Proxy)
// It will override location with CI registry name in case user requests special
// Kubernetes version from CI build area.
// (See: kubeadmconstants.DefaultCIImageRepository)
func (cfg *MasterConfiguration) GetControlPlaneImageRepository() string {
	if cfg.CIImageRepository != "" {
		return cfg.CIImageRepository
	}
	return cfg.ImageRepository
}

// HostPathMount contains elements describing volumes that are mounted from the
// host.
type HostPathMount struct {
	// Name of the volume inside the pod template.
	Name string
	// HostPath is the path in the host that will be mounted inside
	// the pod.
	HostPath string
	// MountPath is the path inside the pod where hostPath will be mounted.
	MountPath string
	// Writable controls write access to the volume
	Writable bool
	// PathType is the type of the HostPath.
	PathType v1.HostPathType
}

// KubeProxy contains elements describing the proxy configuration.
type KubeProxy struct {
	Config *kubeproxyconfigv1alpha1.KubeProxyConfiguration
}

// AuditPolicyConfiguration holds the options for configuring the api server audit policy.
type AuditPolicyConfiguration struct {
	// Path is the local path to an audit policy.
	Path string
	// LogDir is the local path to the directory where logs should be stored.
	LogDir string
	// LogMaxAge is the number of days logs will be stored for. 0 indicates forever.
	LogMaxAge *int32
	// WebhookConfigPath is the local path to webhook policy.
	WebhookConfigPath string
	// WebhookInitialBackoff is the time to wait (in seconds) before retrying the first failed request.
	WebhookInitialBackoff string //defaults to 10s if not provided
	//TODO(chuckha) add other options for audit policy.
}

// CommonConfiguration defines the list of common configuration elements and the getter
// methods that must exist for both the MasterConfiguration and NodeConfiguration objects.
// This is used internally to deduplicate the kubeadm preflight checks.
type CommonConfiguration interface {
	GetCRISocket() string
	GetNodeName() string
	GetKubernetesVersion() string
}

// GetCRISocket will return the CRISocket that is defined for the MasterConfiguration.
// This is used internally to deduplicate the kubeadm preflight checks.
func (cfg *MasterConfiguration) GetCRISocket() string {
	return cfg.CRISocket
}

// GetNodeName will return the NodeName that is defined for the MasterConfiguration.
// This is used internally to deduplicate the kubeadm preflight checks.
func (cfg *MasterConfiguration) GetNodeName() string {
	return cfg.NodeName
}

// GetKubernetesVersion will return the KubernetesVersion that is defined for the MasterConfiguration.
// This is used internally to deduplicate the kubeadm preflight checks.
func (cfg *MasterConfiguration) GetKubernetesVersion() string {
	return cfg.KubernetesVersion
}

// GetCRISocket will return the CRISocket that is defined for the NodeConfiguration.
// This is used internally to deduplicate the kubeadm preflight checks.
func (cfg *NodeConfiguration) GetCRISocket() string {
	return cfg.CRISocket
}

// GetNodeName will return the NodeName that is defined for the NodeConfiguration.
// This is used internally to deduplicate the kubeadm preflight checks.
func (cfg *NodeConfiguration) GetNodeName() string {
	return cfg.NodeName
}

// GetKubernetesVersion will return an empty string since KubernetesVersion is not a
// defined property for NodeConfiguration. This will just cause the regex validation
// of the defined version to be skipped during the preflight checks.
// This is used internally to deduplicate the kubeadm preflight checks.
func (cfg *NodeConfiguration) GetKubernetesVersion() string {
	return ""
}
