package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Keysaas is a specification for a Keysaas resource
// +k8s:openapi-gen=true
type Keysaas struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeysaasSpec   `json:"spec"`
	Status KeysaasStatus `json:"status"`
}

// KeysaasSpec is the spec for a KeysaasSpec resource
// +k8s:openapi-gen=true
type KeysaasSpec struct {
	//MySQL Service name
	MySQLServiceName string `json:"mySQLServiceName"`
	//MySQL Username
	MySQLUserName string `json:"mySQLUserName"`
	//MySQL Password
	MySQLUserPassword string `json:"mySQLUserPassword"`
	//Keysaas Admin Email
	KeysaasAdminEmail string `json:"keysaasAdminEmail"`
	//PVC Volume Name
	PvcVolumeName string `json:"pvcVolumeName"`
	//Domain Name
	DomainName string `json:"domainName"`
	//TLS Flag
	Tls string `json:"tls"`
}

// KeysaasStatus is the status for a Keysaas resource
// +k8s:openapi-gen=true
type KeysaasStatus struct {
	PodName    string `json:"podName"`
	SecretName string `json:"secretName"`
	Status     string `json:"status"`
	Url        string `json:"url"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// KeysaasList is a list of Keysaas resources
type KeysaasList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Keysaas `json:"items"`
}
