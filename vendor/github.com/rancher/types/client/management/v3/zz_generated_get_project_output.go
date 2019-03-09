package client

const (
	GetProjectOutputType             = "getProjectOutput"
	GetProjectOutputFieldMessage     = "message"
	GetProjectOutputFieldProjectInfo = "projectInfo"
)

type GetProjectOutput struct {
	Message     string              `json:"message,omitempty" yaml:"message,omitempty"`
	ProjectInfo []HarborProjectInfo `json:"projectInfo,omitempty" yaml:"projectInfo,omitempty"`
}
