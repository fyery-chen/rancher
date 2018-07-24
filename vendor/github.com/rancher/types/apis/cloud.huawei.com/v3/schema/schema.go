package schema

import (
	"github.com/rancher/norman/types"
	m "github.com/rancher/norman/types/mapper"
	"github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/factory"
)

var (
	Version = types.APIVersion{
		Version: "v3",
		Group:   "cloud.huawei.com",
		Path:    "/v3",
	}

	Schemas = factory.Schemas(&Version).
		Init(authzTypes).
		Init(schemaTypes).
		Init(businessTypes)
)

func businessTypes(schema *types.Schemas) *types.Schemas {
	return schema.
		AddMapperForType(&Version, v3.Business{},
			m.DisplayName{}).
		MustImport(&Version, v3.Business{}).
		MustImport(&Version, v3.BusinessQuotaCheck{}).
		MustImport(&Version, v3.SubnetInfo{}).
		MustImport(&Version, v3.VpcInfo{}).
		MustImport(&Version, v3.NodeFlavor{}).
		MustImport(&Version, v3.SshKeyInfo{}).
		MustImport(&Version, v3.SshKey{}).
		MustImport(&Version, v3.AvailabilityZone{}).
		MustImport(&Version, v3.ZoneState{}).
		MustImport(&Version, v3.HighwaySubnet{}).
		MustImport(&Version, v3.BusinessQuotaCheckOutput{}).
		MustImport(&Version, v3.HuaweiCloudApiInformationOutput{}).
		MustImport(&Version, v3.HuaweiCloudApiInformationInput{}).
		MustImportAndCustomize(&Version, v3.Business{}, func(schema *types.Schema) {
			schema.ResourceActions["checkout"] = types.Action{
				Input:  "businessQuotaCheck",
				Output: "businessQuotaCheckOutput",
			}
			schema.ResourceActions["getHuaweiCloudApiInfo"] = types.Action{
				Input:  "huaweiCloudApiInformationInput",
				Output: "huaweiCloudApiInformationOutput",
			}
		})
}

func schemaTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		MustImport(&Version, v3.DynamicSchema{})
}

func authzTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v3.HuaweiCloudAccount{}, m.DisplayName{}).
		AddMapperForType(&Version, v3.BusinessGlobalRole{}, m.DisplayName{}).
		AddMapperForType(&Version, v3.BusinessRoleTemplate{}, m.DisplayName{}).
		MustImport(&Version, v3.HuaweiCloudAccount{}).
		MustImport(&Version, v3.BusinessGlobalRole{}).
		MustImport(&Version, v3.BusinessGlobalRoleBinding{}).
		MustImport(&Version, v3.BusinessRoleTemplate{}).
		MustImport(&Version, v3.BusinessRoleTemplateBinding{}).
		MustImport(&Version, v3.BusinessGlobalRoleBinding{})
}
