package client

const (
	HarborProjectInfoType                   = "harborProjectInfo"
	HarborProjectInfoFieldChartCount        = "chart_count"
	HarborProjectInfoFieldCreationTime      = "creation_time"
	HarborProjectInfoFieldCurrentUserRoleID = "current_user_role_id"
	HarborProjectInfoFieldDeleted           = "deleted"
	HarborProjectInfoFieldMetadata          = "metadata"
	HarborProjectInfoFieldName              = "name"
	HarborProjectInfoFieldOwnerID           = "owner_id"
	HarborProjectInfoFieldOwnerName         = "owner_name"
	HarborProjectInfoFieldProjectID         = "project_id"
	HarborProjectInfoFieldRepoCount         = "repo_count"
	HarborProjectInfoFieldTogglable         = "togglable"
	HarborProjectInfoFieldUpdateTime        = "update_time"
)

type HarborProjectInfo struct {
	ChartCount        int64            `json:"chart_count,omitempty" yaml:"chart_count,omitempty"`
	CreationTime      string           `json:"creation_time,omitempty" yaml:"creation_time,omitempty"`
	CurrentUserRoleID int64            `json:"current_user_role_id,omitempty" yaml:"current_user_role_id,omitempty"`
	Deleted           bool             `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Metadata          *ProjectMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Name              string           `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerID           int64            `json:"owner_id,omitempty" yaml:"owner_id,omitempty"`
	OwnerName         string           `json:"owner_name,omitempty" yaml:"owner_name,omitempty"`
	ProjectID         int64            `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	RepoCount         int64            `json:"repo_count,omitempty" yaml:"repo_count,omitempty"`
	Togglable         bool             `json:"togglable,omitempty" yaml:"togglable,omitempty"`
	UpdateTime        string           `json:"update_time,omitempty" yaml:"update_time,omitempty"`
}
