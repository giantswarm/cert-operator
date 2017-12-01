package v1alpha1

import (
	"net"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewAWSConfigCRD returns a new custom resource definition for AWSConfig. This
// might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: awsconfigs.cluster.giantswarm.io
//     spec:
//       group: cluster.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: AWSConfig
//         plural: awsconfigs
//         singular: awsconfig
//
func NewAWSConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "awsconfigs.cluster.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "cluster.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "AWSConfig",
				Plural:   "awsconfigs",
				Singular: "awsconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AWSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AWSConfigSpec `json:"spec"`
}

type AWSConfigSpec struct {
	Cluster       AWSConfigSpecCluster       `json:"cluster" yaml:"cluster"`
	AWS           AWSConfigSpecAWS           `json:"aws" yaml:"aws"`
	VersionBundle AWSConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type AWSConfigSpecCluster struct {
	Calico     AWSConfigSpecClusterCalico     `json:"calico" yaml:"calico"`
	Customer   AWSConfigSpecClusterCustomer   `json:"customer" yaml:"customer"`
	Docker     AWSConfigSpecClusterDocker     `json:"docker" yaml:"docker"`
	Etcd       AWSConfigSpecClusterEtcd       `json:"etcd" yaml:"etcd"`
	ID         string                         `json:"id" yaml:"id"`
	Kubernetes AWSConfigSpecClusterKubernetes `json:"kubernetes" yaml:"kubernetes"`
	Masters    []AWSConfigSpecClusterNode     `json:"masters" yaml:"masters"`
	Vault      AWSConfigSpecClusterVault      `json:"vault" yaml:"vault"`
	Workers    []AWSConfigSpecClusterNode     `json:"workers" yaml:"workers"`
}

type AWSConfigSpecClusterCalico struct {
	CIDR   int    `json:"cidr" yaml:"cidr"`
	Domain string `json:"domain" yaml:"domain"`
	MTU    int    `json:"mtu" yaml:"mtu"`
	Subnet string `json:"subnet" yaml:"subnet"`
}

type AWSConfigSpecClusterCustomer struct {
	ID string `json:"id" yaml:"id"`
}

type AWSConfigSpecClusterDocker struct {
	Daemon AWSConfigSpecClusterDockerDaemon `json:"daemon" yaml:"daemon"`
}

type AWSConfigSpecClusterDockerDaemon struct {
	CIDR      string `json:"cidr" yaml:"cidr"`
	ExtraArgs string `json:"extraArgs" yaml:"extraArgs"`
}

type AWSConfigSpecClusterEtcd struct {
	AltNames string `json:"altNames" yaml:"altNames"`
	Domain   string `json:"domain" yaml:"domain"`
	Port     int    `json:"port" yaml:"port"`
	Prefix   string `json:"prefix" yaml:"prefix"`
}

type AWSConfigSpecClusterKubernetes struct {
	API               AWSConfigSpecClusterKubernetesAPI               `json:"api" yaml:"api"`
	DNS               AWSConfigSpecClusterKubernetesDNS               `json:"dns" yaml:"dns"`
	Domain            string                                          `json:"domain" yaml:"domain"`
	Hyperkube         AWSConfigSpecClusterKubernetesHyperkube         `json:"hyperkube" yaml:"hyperkube"`
	IngressController AWSConfigSpecClusterKubernetesIngressController `json:"ingressController" yaml:"ingressController"`
	Kubelet           AWSConfigSpecClusterKubernetesKubelet           `json:"kubelet" yaml:"kubelet"`
	NetworkSetup      AWSConfigSpecClusterKubernetesNetworkSetup      `json:"networkSetup" yaml:"networkSetup"`
	SSH               AWSConfigSpecClusterKubernetesSSH               `json:"ssh" yaml:"ssh"`
}

type AWSConfigSpecClusterKubernetesAPI struct {
	AltNames       string `json:"altNames" yaml:"altNames"`
	ClusterIPRange string `json:"clusterIPRange" yaml:"clusterIPRange"`
	Domain         string `json:"domain" yaml:"domain"`
	IP             net.IP `json:"ip" yaml:"ip"`
	InsecurePort   int    `json:"insecurePort" yaml:"insecurePort"`
	SecurePort     int    `json:"securePort" yaml:"securePort"`
}

type AWSConfigSpecClusterKubernetesDNS struct {
	IP net.IP `json:"ip" yaml:"ip"`
}

type AWSConfigSpecClusterKubernetesHyperkube struct {
	Docker AWSConfigSpecClusterKubernetesHyperkubeDocker `json:"docker" yaml:"docker"`
}

type AWSConfigSpecClusterKubernetesHyperkubeDocker struct {
	Image string `json:"image" yaml:"image"`
}

type AWSConfigSpecClusterKubernetesIngressController struct {
	Docker         AWSConfigSpecClusterKubernetesIngressControllerDocker `json:"docker" yaml:"docker"`
	Domain         string                                                `json:"domain" yaml:"domain"`
	WildcardDomain string                                                `json:"wildcardDomain" yaml:"wildcardDomain"`
	InsecurePort   int                                                   `json:"insecurePort" yaml:"insecurePort"`
	SecurePort     int                                                   `json:"securePort" yaml:"securePort"`
}

type AWSConfigSpecClusterKubernetesIngressControllerDocker struct {
	Image string `json:"image" yaml:"image"`
}

type AWSConfigSpecClusterKubernetesKubelet struct {
	AltNames string `json:"altNames" yaml:"altNames"`
	Domain   string `json:"domain" yaml:"domain"`
	Labels   string `json:"labels" yaml:"labels"`
	Port     int    `json:"port" yaml:"port"`
}

type AWSConfigSpecClusterKubernetesNetworkSetup struct {
	Docker AWSConfigSpecClusterKubernetesNetworkSetupDocker `json:"docker" yaml:"docker"`
}

type AWSConfigSpecClusterKubernetesNetworkSetupDocker struct {
	Image string `json:"image" yaml:"image"`
}

type AWSConfigSpecClusterKubernetesSSH struct {
	UserList []AWSConfigSpecClusterKubernetesSSHUser `json:"userList" yaml:"userList"`
}

type AWSConfigSpecClusterKubernetesSSHUser struct {
	Name      string `json:"name" yaml:"name"`
	PublicKey string `json:"publicKey" yaml:"publicKey"`
}

type AWSConfigSpecClusterNode struct {
	ID string `json:"id" yaml:"id"`
}

type AWSConfigSpecClusterVault struct {
	Address string `json:"address" yaml:"address"`
	Token   string `json:"token" yaml:"token"`
}

type AWSConfigSpecAWS struct {
	API     AWSConfigSpecAWSAPI     `json:"api" yaml:"api"`
	AZ      string                  `json:"az" yaml:"az"`
	Etcd    AWSConfigSpecAWSEtcd    `json:"etcd" yaml:"etcd"`
	Ingress AWSConfigSpecAWSIngress `json:"ingress" yaml:"ingress"`
	Masters []AWSConfigSpecAWSNode  `json:"masters" yaml:"masters"`
	Region  string                  `json:"region" yaml:"region"`
	VPC     AWSConfigSpecAWSVPC     `json:"vpc" yaml:"vpc"`
	Workers []AWSConfigSpecAWSNode  `json:"workers" yaml:"workers"`
}

type AWSConfigSpecAWSAPI struct {
	HostedZones string                 `json:"hostedZones" yaml:"hostedZones"`
	ELB         AWSConfigSpecAWSAPIELB `json:"elb" yaml:"elb"`
}

type AWSConfigSpecAWSAPIELB struct {
	IdleTimeoutSeconds int `json:"idleTimeoutSeconds" yaml:"idleTimeoutSeconds"`
}

type AWSConfigSpecAWSEtcd struct {
	HostedZones string                  `json:"hostedZones" yaml:"hostedZones"`
	ELB         AWSConfigSpecAWSEtcdELB `json:"elb" yaml:"elb"`
}

type AWSConfigSpecAWSEtcdELB struct {
	IdleTimeoutSeconds int `json:"idleTimeoutSeconds" yaml:"idleTimeoutSeconds"`
}

type AWSConfigSpecAWSIngress struct {
	HostedZones string                     `json:"hostedZones" yaml:"hostedZones"`
	ELB         AWSConfigSpecAWSIngressELB `json:"elb" yaml:"elb"`
}

type AWSConfigSpecAWSIngressELB struct {
	IdleTimeoutSeconds int `json:"idleTimeoutSeconds" yaml:"idleTimeoutSeconds"`
}

type AWSConfigSpecAWSNode struct {
	ImageID      string `json:"imageID" yaml:"imageID"`
	InstanceType string `json:"instanceType" yaml:"instanceType"`
}

type AWSConfigSpecAWSVPC struct {
	CIDR              string   `json:"cidr" yaml:"cidr"`
	PrivateSubnetCIDR string   `json:"privateSubnetCidr" yaml:"privateSubnetCidr"`
	PublicSubnetCIDR  string   `json:"publicSubnetCidr" yaml:"publicSubnetCidr"`
	RouteTableNames   []string `json:"routeTableNames" yaml:"routeTableNames"`
	PeerID            string   `json:"peerId" yaml:"peerId"`
}

type AWSConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AWSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AWSConfig `json:"items"`
}
