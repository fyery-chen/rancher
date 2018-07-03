package schema

import (
	"github.com/rancher/norman/types"
	m "github.com/rancher/norman/types/mapper"
	"github.com/rancher/types/apis/cloud.huawei.com/v1"
	"github.com/rancher/types/factory"
)

var (
	Version = types.APIVersion{
		Version: "v1",
		Group:   "cloud.huawei.com",
		Path:    "/v1",
	}

	Schemas = factory.Schemas(&Version).
		Init(businessTypes)
)

func businessTypes(schema *types.Schemas) *types.Schemas {
	return schema.
		AddMapperForType(&Version, v1.BusinessQuota{},
			m.DisplayName{}).
		MustImport(&Version, v1.BusinessQuota{}).
		MustImportAndCustomize(&Version, v1.BusinessQuota{}, func(schema *types.Schema) {
			schema.ResourceActions = map[string]types.Action{
				"checkout": {
					Input: "businessQuotaCheck",
				},
			}
		}).
		MustImport(&Version, v1.BusinessQuotaCheck{})
}
