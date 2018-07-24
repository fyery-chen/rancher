package client

const (
	SshKeyInfoType      = "sshKeyInfo"
	SshKeyInfoFieldName = "name"
)

type SshKeyInfo struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
