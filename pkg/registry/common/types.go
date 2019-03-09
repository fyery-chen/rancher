package common

const (
	Project    = "harborProject"
	Repository = "harborRepository"
	Tag        = "harborTag"
	All        = "harborAll"
)

type APIClientConfig struct {
	RegistryServer string
	Username       string
	Password       string
	RootCA         string
	ClientCert     string
	ClientKey      string
	Proxy          string
	RequestType    string
	Page           int
	PageSize       int
	ProjectID      int
	RepositoryName string
}

type ProjectMetadata struct {
	Public             string `json:"public,omitempty"`
	EnableContentTrust string `json:"enable_content_trust,omitempty"`
	PreventVul         string `json:"prevent_vul,omitempty"`
	Severity           string `json:"severity,omitempty"`
	AutoScan           string `json:"auto_scan,omitempty"`
}

type HarborProject struct {
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
	ID           int               `json:"id,omitempty"`
	Name         string            `json:"name,omitempty"`
	ProjectID    int               `json:"project_id,omitempty"`
	Description  string            `json:"description,omitempty"`
	PullCount    int               `json:"pull_count,omitempty"`
	StarCount    int               `json:"star_count,omitempty"`
	TagsCount    int               `json:"tags_count,omitempty"`
	Labels       []RepositoryLabel `json:"labels,omitempty"`
	CreationTime string            `json:"creation_time,omitempty"`
	UpdateTime   string            `json:"update_time,omitempty"`
}

type RepositoryTag struct {
	Digest string `json:"digest,omitempty"`
	Name   string `json:"name,omitempty"`
	Size   string `json:"size,omitempty"`
}
