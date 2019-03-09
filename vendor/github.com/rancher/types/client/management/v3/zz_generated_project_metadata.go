package client

const (
	ProjectMetadataType                    = "projectMetadata"
	ProjectMetadataFieldAutoScan           = "auto_scan"
	ProjectMetadataFieldEnableContentTrust = "enable_content_trust"
	ProjectMetadataFieldPreventVul         = "prevent_vul"
	ProjectMetadataFieldPublic             = "public"
	ProjectMetadataFieldSeverity           = "severity"
)

type ProjectMetadata struct {
	AutoScan           string `json:"auto_scan,omitempty" yaml:"auto_scan,omitempty"`
	EnableContentTrust string `json:"enable_content_trust,omitempty" yaml:"enable_content_trust,omitempty"`
	PreventVul         string `json:"prevent_vul,omitempty" yaml:"prevent_vul,omitempty"`
	Public             string `json:"public,omitempty" yaml:"public,omitempty"`
	Severity           string `json:"severity,omitempty" yaml:"severity,omitempty"`
}
