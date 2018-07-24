package client

const (
	SshKeyType         = "sshKey"
	SshKeyFieldKeypair = "keypair"
)

type SshKey struct {
	Keypair *SshKeyInfo `json:"keypair,omitempty" yaml:"keypair,omitempty"`
}
