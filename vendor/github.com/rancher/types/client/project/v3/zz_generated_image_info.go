package client

const (
	ImageInfoType              = "imageInfo"
	ImageInfoFieldInternalPath = "internal_path"
	ImageInfoFieldName         = "name"
	ImageInfoFieldNamespace    = "namespace"
	ImageInfoFieldPath         = "path"
	ImageInfoFieldTags         = "tags"
)

type ImageInfo struct {
	InternalPath string   `json:"internal_path,omitempty" yaml:"internal_path,omitempty"`
	Name         string   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace    string   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Path         string   `json:"path,omitempty" yaml:"path,omitempty"`
	Tags         []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}
