package v3

import (
	"github.com/rancher/norman/condition"
	"github.com/rancher/norman/types"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	NamespaceBackedResource     condition.Cond = "BackingNamespaceCreated"
	CreatorMadeOwner            condition.Cond = "CreatorMadeOwner"
	DefaultNetworkPolicyCreated condition.Cond = "DefaultNetworkPolicyCreated"
)

type HuaweiCloudAccountSpec struct {
	DisplayName string `json:"displayName,omitempty" norman:"required"`
	AccessKey   string `json:"accessKey,omitempty"`
	SecretKey   string `json:"secretKey,omitempty"`
	UserName    string `json:"userName,omitempty" norman:"required,type=reference[user]"`
}

type HuaweiCloudAccount struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HuaweiCloudAccountSpec `json:"spec,omitempty"`
}

type BusinessGlobalRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	DisplayName string              `json:"displayName,omitempty" norman:"required"`
	Description string              `json:"description"`
	Rules       []rbacv1.PolicyRule `json:"rules,omitempty"`
	Builtin     bool                `json:"builtin" norman:"nocreate,noupdate"`
}

type BusinessGlobalRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	UserName       string `json:"userName,omitempty" norman:"required,type=reference[user]"`
	GlobalRoleName string `json:"globalRoleName,omitempty" norman:"required,type=reference[businessGlobalRole]"`
}

type BusinessRoleTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	DisplayName       string              `json:"displayName,omitempty" norman:"required"`
	Description       string              `json:"description"`
	Rules             []rbacv1.PolicyRule `json:"rules,omitempty"`
	Builtin           bool                `json:"builtin" norman:"nocreate,noupdate"`
	External          bool                `json:"external"`
	Hidden            bool                `json:"hidden"`
	Locked            bool                `json:"locked,omitempty" norman:"type=boolean"`
	Context           string              `json:"context" norman:"type=string,options=project|cluster|business"`
	RoleTemplateNames []string            `json:"roleTemplateNames,omitempty" norman:"type=array[reference[businessRoleTemplate]]"`
}

type BusinessRoleTemplateBinding struct {
	types.Namespaced
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	UserName           string `json:"userName,omitempty" norman:"type=reference[user]"`
	UserPrincipalName  string `json:"userPrincipalName,omitempty" norman:"type=reference[principal]"`
	GroupName          string `json:"groupName,omitempty" norman:"type=reference[group]"`
	GroupPrincipalName string `json:"groupPrincipalName,omitempty" norman:"type=reference[principal]"`
	BusinessName       string `json:"businessName,omitempty" norman:"required,type=reference[business]"`
	RoleTemplateName   string `json:"roleTemplateName,omitempty" norman:"required,type=reference[businessRoleTemplate]"`
}
