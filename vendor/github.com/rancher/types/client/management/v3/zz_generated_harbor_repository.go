package client

const (
	HarborRepositoryType              = "harborRepository"
	HarborRepositoryFieldCreationTime = "creation_time"
	HarborRepositoryFieldDescription  = "description"
	HarborRepositoryFieldName         = "name"
	HarborRepositoryFieldPullCount    = "pull_count"
	HarborRepositoryFieldTagsCount    = "tags_count"
	HarborRepositoryFieldUpdateTime   = "update_time"
)

type HarborRepository struct {
	CreationTime string `json:"creation_time,omitempty" yaml:"creation_time,omitempty"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty"`
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	PullCount    int64  `json:"pull_count,omitempty" yaml:"pull_count,omitempty"`
	TagsCount    int64  `json:"tags_count,omitempty" yaml:"tags_count,omitempty"`
	UpdateTime   string `json:"update_time,omitempty" yaml:"update_time,omitempty"`
}
