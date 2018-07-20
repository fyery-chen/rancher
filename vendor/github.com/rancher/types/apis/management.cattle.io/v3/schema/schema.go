package schema

import (
	"net/http"

	"github.com/rancher/norman/types"
	m "github.com/rancher/norman/types/mapper"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/factory"
)

var (
	Version = types.APIVersion{
		Version: "v3",
		Group:   "management.cattle.io",
		Path:    "/v3",
	}

	Schemas = factory.Schemas(&Version).
		Init(authnTypes).
		Init(tokens).
		Init(schemaTypes).
		Init(userTypes).
		Init(globalTypes).
		Init(rkeTypes).
		Init(clusterTypes)

	TokenSchemas = factory.Schemas(&Version).
			Init(tokens)
)

func rkeTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.AddMapperForType(&Version, v3.BaseService{}, m.Drop{Field: "image"})
}

func clusterTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v3.Cluster{},
			&m.Embed{Field: "status"},
			m.DisplayName{},
		).
		AddMapperForType(&Version, v3.ClusterStatus{},
			m.Drop{Field: "serviceAccountToken"},
		).
		AddMapperForType(&Version, v3.ClusterEvent{}, &m.Move{
			From: "type",
			To:   "eventType",
		}).
		AddMapperForType(&Version, v3.ClusterRegistrationToken{},
			&m.Embed{Field: "status"},
		).
		MustImport(&Version, v3.Cluster{}).
		MustImport(&Version, v3.ClusterEvent{}).
		MustImport(&Version, v3.ClusterRegistrationToken{}).
		MustImport(&Version, v3.GenerateKubeConfigOutput{}).
		MustImport(&Version, v3.ImportClusterYamlInput{}).
		MustImport(&Version, v3.ImportYamlOutput{}).
		MustImport(&Version, v3.ExportOutput{}).
		MustImportAndCustomize(&Version, v3.Cluster{}, func(schema *types.Schema) {
			schema.MustCustomizeField("name", func(field types.Field) types.Field {
				field.Type = "dnsLabel"
				field.Nullable = true
				field.Required = false
				return field
			})
			schema.ResourceActions["generateKubeconfig"] = types.Action{
				Output: "generateKubeConfigOutput",
			}
			schema.ResourceActions["importYaml"] = types.Action{
				Input:  "importClusterYamlInput",
				Output: "importYamlOutput",
			}
			schema.ResourceActions["exportYaml"] = types.Action{
				Output: "exportOutput",
			}
		})
}

func schemaTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		MustImport(&Version, v3.DynamicSchema{})
}

func tokens(schemas *types.Schemas) *types.Schemas {
	return schemas.
		MustImportAndCustomize(&Version, v3.Token{}, func(schema *types.Schema) {
			schema.CollectionActions = map[string]types.Action{
				"logout": {},
			}
		})
}

func authnTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v3.User{}, m.DisplayName{}).
		AddMapperForType(&Version, v3.Group{}, m.DisplayName{}).
		MustImport(&Version, v3.Group{}).
		MustImport(&Version, v3.GroupMember{}).
		AddMapperForType(&Version, v3.Principal{}, m.DisplayName{}).
		MustImportAndCustomize(&Version, v3.Principal{}, func(schema *types.Schema) {
			schema.CollectionMethods = []string{http.MethodGet}
			schema.ResourceMethods = []string{http.MethodGet}
			schema.CollectionActions = map[string]types.Action{
				"search": {
					Input:  "searchPrincipalsInput",
					Output: "collection",
				},
			}
		}).
		MustImport(&Version, v3.SearchPrincipalsInput{}).
		MustImport(&Version, v3.ChangePasswordInput{}).
		MustImport(&Version, v3.SetPasswordInput{}).
		MustImportAndCustomize(&Version, v3.User{}, func(schema *types.Schema) {
			schema.ResourceActions = map[string]types.Action{
				"setpassword": {
					Input:  "setPasswordInput",
					Output: "user",
				},
			}
			schema.CollectionActions = map[string]types.Action{
				"changepassword": {
					Input: "changePasswordInput",
				},
			}
		}).
		MustImportAndCustomize(&Version, v3.AuthConfig{}, func(schema *types.Schema) {
			schema.CollectionMethods = []string{http.MethodGet}
		}).
		// Local Config
		MustImportAndCustomize(&Version, v3.LocalConfig{}, func(schema *types.Schema) {
			schema.BaseType = "authConfig"
			schema.CollectionMethods = []string{}
			schema.ResourceMethods = []string{http.MethodGet}
		}).
		//Github Config
		MustImportAndCustomize(&Version, v3.GithubConfig{}, func(schema *types.Schema) {
			schema.BaseType = "authConfig"
			schema.ResourceActions = map[string]types.Action{
				"disable": {},
				"configureTest": {
					Input:  "githubConfig",
					Output: "githubConfigTestOutput",
				},
				"testAndApply": {
					Input: "githubConfigApplyInput",
				},
			}
			schema.CollectionMethods = []string{}
			schema.ResourceMethods = []string{http.MethodGet, http.MethodPut}
		}).
		MustImport(&Version, v3.GithubConfigTestOutput{}).
		MustImport(&Version, v3.GithubConfigApplyInput{}).
		// Active Directory Config
		MustImportAndCustomize(&Version, v3.ActiveDirectoryConfig{}, func(schema *types.Schema) {
			schema.BaseType = "authConfig"
			schema.ResourceActions = map[string]types.Action{
				"disable": {},
				"testAndApply": {
					Input: "activeDirectoryTestAndApplyInput",
				},
			}
			schema.CollectionMethods = []string{}
			schema.ResourceMethods = []string{http.MethodGet, http.MethodPut}
		}).
		MustImport(&Version, v3.ActiveDirectoryTestAndApplyInput{})
}

func userTypes(schema *types.Schemas) *types.Schemas {
	return schema.
		MustImportAndCustomize(&Version, v3.Preference{}, func(schema *types.Schema) {
			schema.MustCustomizeField("name", func(f types.Field) types.Field {
				f.Required = true
				return f
			})
			schema.MustCustomizeField("namespaceId", func(f types.Field) types.Field {
				f.Required = false
				return f
			})
		})
}

func globalTypes(schema *types.Schemas) *types.Schemas {
	return schema.
		AddMapperForType(&Version, v3.ListenConfig{},
			m.DisplayName{},
			m.Drop{Field: "caKey"},
			m.Drop{Field: "caCert"},
		).
		MustImport(&Version, v3.ListenConfig{}).
		MustImportAndCustomize(&Version, v3.Setting{}, func(schema *types.Schema) {
			schema.MustCustomizeField("name", func(f types.Field) types.Field {
				f.Required = true
				return f
			})
		})
}
