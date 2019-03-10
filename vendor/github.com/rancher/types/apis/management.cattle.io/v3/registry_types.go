package v3

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GlobalRegistry struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec GlobalRegistrySpec `json:"spec"`
	// Most recent observed status of the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status RegistryStatus `json:"status"`
}

type ClusterRegistry struct {
	types.Namespaced

	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec ClusterRegistrySpec `json:"spec"`
	// Most recent observed status of the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status RegistryStatus `json:"status"`
}

type ProjectRegistry struct {
	types.Namespaced

	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec ProjectRegistrySpec `json:"spec"`
	// Most recent observed status of the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status RegistryStatus `json:"status"`
}

type RegistryCommon struct {
	DisplayName   string `json:"displayName,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty"`
	UserName      string `json:"userName,omitempty"`
	Password      string `json:"password,omitempty"`
	RegistryType  string `json:"registryType,omitempty"`
}

type GlobalRegistrySpec struct {
	RegistryCommon
}

type ClusterRegistrySpec struct {
	RegistryCommon
	ClusterName            string `json:"clusterName" norman:"type=reference[cluster]"`
	IncludeSystemComponent *bool  `json:"includeSystemComponent,omitempty" norman:"default=true"`
}

type ProjectRegistrySpec struct {
	RegistryCommon
	ProjectName string `json:"projectName" norman:"type=reference[project]"`
}

type RegistryStatus struct {
	RegistryState string `json:"RegistryState,omitempty" norman:"options=active|inactive,default=active"`
}

type GetRegistryInfoCommon struct {
	RegistryName string `json:"registryName,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
}

type GetProjectInput struct {
	GetRegistryInfoCommon
}

type GetRepositoryInput struct {
	GetRegistryInfoCommon
	ProjectID int  `json:"projectId,omitempty"`
	IsAll     bool `json:"isAll,omitempty"`
}

type GetRepositoryTagsInput struct {
	GetRegistryInfoCommon
	RepositoryName string `json:"repositoryName,omitempty"`
}

type ClusterRegistryTestInput struct {
	ClusterName string `json:"clusterId" norman:"required,type=reference[cluster]"`
	RegistryCommon
}

type ProjectRegistryTestInput struct {
	ProjectName string `json:"projectId" norman:"required,type=reference[project]"`
	RegistryCommon
}

type GlobalRegistryTestInput struct {
	RegistryCommon
}

type ProjectMetadata struct {
	Public             string `json:"public,omitempty"`
	EnableContentTrust string `json:"enable_content_trust,omitempty"`
	PreventVul         string `json:"prevent_vul,omitempty"`
	Severity           string `json:"severity,omitempty"`
	AutoScan           string `json:"auto_scan,omitempty"`
}

type HarborProjectInfo struct {
	ProjectID         int             `json:"project_id,omitempty"`
	OwnerID           int             `json:"owner_id,omitempty"`
	Name              string          `json:"name,omitempty"`
	CreationTime      string          `json:"creation_time,omitempty"`
	UpdateTime        string          `json:"update_time,omitempty"`
	Deleted           bool            `json:"deleted,omitempty"`
	OwnerName         string          `json:"owner_name,omitempty"`
	Togglable         bool            `json:"togglable,omitempty"`
	CurrentUserRoleID int             `json:"current_user_role_id,omitempty"`
	RepoCount         int             `json:"repo_count,omitempty"`
	ChartCount        int             `json:"chart_count,omitempty"`
	Metadata          ProjectMetadata `json:"metadata,omitempty"`
}

type GetProjectOutput struct {
	ProjectInfo []HarborProjectInfo `json:"projectInfo,omitempty"`
	Message     string              `json:"message,omitempty"`
}

type RepositoryLabel struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Color        string `json:"color,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ProjectID    int    `json:"project_id,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	Deleted      bool   `json:"deleted,omitempty"`
}

type HarborRepository struct {
	//ID           int               `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	//ProjectID    int               `json:"project_id,omitempty"`
	Description string `json:"description,omitempty"`
	PullCount   int    `json:"pull_count,omitempty"`
	//StarCount    int               `json:"star_count,omitempty"`
	TagsCount int `json:"tags_count,omitempty"`
	//Labels       []RepositoryLabel `json:"labels,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
}

type GetRepositoryOutput struct {
	RepositoryInfo []HarborRepository `json:"repositoryInfo,omitempty"`
	Message        string             `json:"message,omitempty"`
}

type RepositoryTag struct {
	Digest string `json:"digest,omitempty"`
	Name   string `json:"name,omitempty"`
	Size   int    `json:"size,omitempty"`
}

type GetRepositoryTagsOutput struct {
	RepositoryTagsInfo []RepositoryTag `json:"repositoryTagsInfo,omitempty"`
	Message            string          `json:"message,omitempty"`
}
