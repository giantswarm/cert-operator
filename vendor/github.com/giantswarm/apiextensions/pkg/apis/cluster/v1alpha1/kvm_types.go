package v1alpha1

import (
	"net"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewKVMConfigCRD returns a new custom resource definition for KVMConfig. This
// might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: kvmconfigs.cluster.giantswarm.io
//     spec:
//       group: cluster.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: KVMConfig
//         plural: kvmconfigs
//         singular: kvmconfig
//
func NewKVMConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kvmconfigs.cluster.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "cluster.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "KVMConfig",
				Plural:   "kvmconfigs",
				Singular: "kvmconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KVMConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              KVMConfigSpec `json:"spec"`
}

type KVMConfigSpec struct {
	Cluster       KVMConfigSpecCluster       `json:"cluster" yaml:"cluster"`
	KVMConfig     KVMConfigSpecKVM           `json:"kvm" yaml:"kvm"`
	VersionBundle KVMConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type KVMConfigSpecCluster struct {
	Calico     KVMConfigSpecClusterCalico     `json:"calico" yaml:"calico"`
	Customer   KVMConfigSpecClusterCustomer   `json:"customer" yaml:"customer"`
	Docker     KVMConfigSpecClusterDocker     `json:"docker" yaml:"docker"`
	Etcd       KVMConfigSpecClusterEtcd       `json:"etcd" yaml:"etcd"`
	ID         string                         `json:"id" yaml:"id"`
	Kubernetes KVMConfigSpecClusterKubernetes `json:"kubernetes" yaml:"kubernetes"`
	Masters    []KVMConfigSpecClusterNode     `json:"masters" yaml:"masters"`
	Vault      KVMConfigSpecClusterVault      `json:"vault" yaml:"vault"`
	Workers    []KVMConfigSpecClusterNode     `json:"workers" yaml:"workers"`
}

type KVMConfigSpecClusterCalico struct {
	CIDR   int    `json:"cidr" yaml:"cidr"`
	Domain string `json:"domain" yaml:"domain"`
	MTU    int    `json:"mtu" yaml:"mtu"`
	Subnet string `json:"subnet" yaml:"subnet"`
}

type KVMConfigSpecClusterCustomer struct {
	ID string `json:"id" yaml:"id"`
}

type KVMConfigSpecClusterDocker struct {
	Daemon KVMConfigSpecClusterDockerDaemon `json:"daemon" yaml:"daemon"`
}

type KVMConfigSpecClusterDockerDaemon struct {
	CIDR      string `json:"cidr" yaml:"cidr"`
	ExtraArgs string `json:"extraArgs" yaml:"extraArgs"`
}

type KVMConfigSpecClusterEtcd struct {
	AltNames string `json:"altNames" yaml:"altNames"`
	Domain   string `json:"domain" yaml:"domain"`
	Port     int    `json:"port" yaml:"port"`
	Prefix   string `json:"prefix" yaml:"prefix"`
}

type KVMConfigSpecClusterKubernetes struct {
	API               KVMConfigSpecClusterKubernetesAPI               `json:"api" yaml:"api"`
	DNS               KVMConfigSpecClusterKubernetesDNS               `json:"dns" yaml:"dns"`
	Domain            string                                          `json:"domain" yaml:"domain"`
	Hyperkube         KVMConfigSpecClusterKubernetesHyperkube         `json:"hyperkube" yaml:"hyperkube"`
	IngressController KVMConfigSpecClusterKubernetesIngressController `json:"ingressController" yaml:"ingressController"`
	Kubelet           KVMConfigSpecClusterKubernetesKubelet           `json:"kubelet" yaml:"kubelet"`
	NetworkSetup      KVMConfigSpecClusterKubernetesNetworkSetup      `json:"networkSetup" yaml:"networkSetup"`
	SSH               KVMConfigSpecClusterKubernetesSSH               `json:"ssh" yaml:"ssh"`
}

type KVMConfigSpecClusterKubernetesAPI struct {
	AltNames       string `json:"altNames" yaml:"altNames"`
	ClusterIPRange string `json:"clusterIPRange" yaml:"clusterIPRange"`
	Domain         string `json:"domain" yaml:"domain"`
	IP             net.IP `json:"ip" yaml:"ip"`
	InsecurePort   int    `json:"insecurePort" yaml:"insecurePort"`
	SecurePort     int    `json:"securePort" yaml:"securePort"`
}

type KVMConfigSpecClusterKubernetesDNS struct {
	IP net.IP `json:"ip" yaml:"ip"`
}

type KVMConfigSpecClusterKubernetesHyperkube struct {
	Docker KVMConfigSpecClusterKubernetesHyperkubeDocker `json:"docker" yaml:"docker"`
}

type KVMConfigSpecClusterKubernetesHyperkubeDocker struct {
	Image string `json:"image" yaml:"image"`
}

type KVMConfigSpecClusterKubernetesIngressController struct {
	Docker         KVMConfigSpecClusterKubernetesIngressControllerDocker `json:"docker" yaml:"docker"`
	Domain         string                                                `json:"domain" yaml:"domain"`
	WildcardDomain string                                                `json:"wildcardDomain" yaml:"wildcardDomain"`
	InsecurePort   int                                                   `json:"insecurePort" yaml:"insecurePort"`
	SecurePort     int                                                   `json:"securePort" yaml:"securePort"`
}

type KVMConfigSpecClusterKubernetesIngressControllerDocker struct {
	Image string `json:"image" yaml:"image"`
}

type KVMConfigSpecClusterKubernetesKubelet struct {
	AltNames string `json:"altNames" yaml:"altNames"`
	Domain   string `json:"domain" yaml:"domain"`
	Labels   string `json:"labels" yaml:"labels"`
	Port     int    `json:"port" yaml:"port"`
}

type KVMConfigSpecClusterKubernetesNetworkSetup struct {
	Docker KVMConfigSpecClusterKubernetesNetworkSetupDocker `json:"docker" yaml:"docker"`
}

type KVMConfigSpecClusterKubernetesNetworkSetupDocker struct {
	Image string `json:"image" yaml:"image"`
}

type KVMConfigSpecClusterKubernetesSSH struct {
	UserList []KVMConfigSpecClusterKubernetesSSHUser `json:"userList" yaml:"userList"`
}

type KVMConfigSpecClusterKubernetesSSHUser struct {
	Name      string `json:"name" yaml:"name"`
	PublicKey string `json:"publicKey" yaml:"publicKey"`
}

type KVMConfigSpecClusterNode struct {
	ID string `json:"id" yaml:"id"`
}

type KVMConfigSpecClusterVault struct {
	Address string `json:"address" yaml:"address"`
	Token   string `json:"token" yaml:"token"`
}

type KVMConfigSpecKVM struct {
	EndpointUpdater KVMConfigSpecKVMEndpointUpdater `json:"endpointUpdater" yaml:"endpointUpdater"`
	K8sKVM          KVMConfigSpecKVMK8sKVM          `json:"k8sKVM" yaml:"k8sKVM"`
	Masters         []KVMConfigSpecKVMNode          `json:"masters" yaml:"masters"`
	Network         KVMConfigSpecKVMNetwork         `json:"network" yaml:"network"`
	NodeController  KVMConfigSpecKVMNodeController  `json:"nodeController" yaml:"nodeController"`
	Workers         []KVMConfigSpecKVMNode          `json:"workers" yaml:"workers"`
}

type KVMConfigSpecKVMEndpointUpdater struct {
	Docker KVMConfigSpecKVMEndpointUpdaterDocker `json:"docker" yaml:"docker"`
}

type KVMConfigSpecKVMEndpointUpdaterDocker struct {
	Image string `json:"image" yaml:"image"`
}

type KVMConfigSpecKVMK8sKVM struct {
	Docker      KVMConfigSpecKVMK8sKVMDocker `json:"docker" yaml:"docker"`
	StorageType string                       `json:"storageType" yaml:"storageType"`
}

type KVMConfigSpecKVMK8sKVMDocker struct {
	Image string `json:"image" yaml:"image"`
}

type KVMConfigSpecKVMNode struct {
	CPUs   int     `json:"cpus" yaml:"cpus"`
	Disk   float64 `json:"disk" yaml:"disk"`
	Memory string  `json:"memory" yaml:"memory"`
}

type KVMConfigSpecKVMNetwork struct {
	Flannel KVMConfigSpecKVMNetworkFlannel `json:"flannel" yaml:"flannel"`
}

type KVMConfigSpecKVMNetworkFlannel struct {
	VNI int `json:"vni" yaml:"vni"`
}

type KVMConfigSpecKVMNodeController struct {
	Docker KVMConfigSpecKVMNodeControllerDocker `json:"docker" yaml:"docker"`
}

type KVMConfigSpecKVMNodeControllerDocker struct {
	Image string `json:"image" yaml:"image"`
}

type KVMConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KVMConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []KVMConfig `json:"items"`
}
