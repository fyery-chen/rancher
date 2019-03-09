package client

const (
	RepositoryTagType        = "repositoryTag"
	RepositoryTagFieldDigest = "digest"
	RepositoryTagFieldName   = "name"
	RepositoryTagFieldSize   = "size"
)

type RepositoryTag struct {
	Digest string `json:"digest,omitempty" yaml:"digest,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Size   string `json:"size,omitempty" yaml:"size,omitempty"`
}
