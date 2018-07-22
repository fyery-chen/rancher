package managementstored

import (
	"context"
	"sync"

	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/api/customization/authn"
	"github.com/rancher/rancher/pkg/api/customization/roletemplatebinding"
	"github.com/rancher/rancher/pkg/api/customization/huaweicloudapi"
	"github.com/rancher/rancher/pkg/api/customization/setting"
	"github.com/rancher/rancher/pkg/api/store/preference"
	"github.com/rancher/rancher/pkg/auth/principals"
	"github.com/rancher/rancher/pkg/auth/providers"
	managementschema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	businessschema "github.com/rancher/types/apis/cloud.huawei.com/v3/schema"
	"github.com/rancher/types/client/management/v3"
	businessclient "github.com/rancher/types/client/cloud/v3"
	"github.com/rancher/types/config"
)

func Setup(ctx context.Context, apiContext *config.ScaledContext) error {
	// Here we setup all types that will be stored in the Management cluster
	schemas := apiContext.Schemas

	wg := &sync.WaitGroup{}
	factory := &crd.Factory{ClientGetter: apiContext.ClientGetter}

	createCrd(ctx, wg, factory, schemas, &managementschema.Version,
		client.AuthConfigType,
		client.DynamicSchemaType,
		client.GroupMemberType,
		client.GroupType,
		client.ListenConfigType,
		client.PreferenceType,
		client.SettingType,
		client.TokenType,
		client.UserType)
	createCrd(ctx, wg, factory, schemas, &businessschema.Version,
		businessclient.HuaweiCloudAccountType,
		businessclient.BusinessRoleTemplateBindingType,
		businessclient.BusinessRoleTemplateType,
		businessclient.BusinessGlobalRoleBindingType,
		businessclient.BusinessGlobalRoleType,
		businessclient.BusinessType)

	wg.Wait()

	User(schemas, apiContext)
	Setting(schemas)
	Preference(schemas, apiContext)
	BusinessRoleTemplateBinding(schemas, apiContext)
	HuaweiCloudApi(schemas, apiContext)

	principals.Schema(ctx, apiContext, schemas)
	providers.SetupAuthConfig(ctx, apiContext, schemas)
	authn.SetUserStore(schemas.Schema(&managementschema.Version, client.UserType), apiContext)
	authn.SetRTBStore(ctx, schemas.Schema(&businessschema.Version, businessclient.BusinessRoleTemplateBindingType), apiContext)

	return nil
}

func User(schemas *types.Schemas, management *config.ScaledContext) {
	schema := schemas.Schema(&managementschema.Version, client.UserType)
	schema.Formatter = authn.UserFormatter
	schema.CollectionFormatter = authn.CollectionFormatter
	handler := &authn.Handler{
		UserClient: management.Management.Users(""),
	}
	schema.ActionHandler = handler.Actions
}

func Preference(schemas *types.Schemas, management *config.ScaledContext) {
	schema := schemas.Schema(&managementschema.Version, client.PreferenceType)
	schema.Store = preference.NewStore(management.Core.Namespaces(""), schema.Store)
}

func Setting(schemas *types.Schemas) {
	schema := schemas.Schema(&managementschema.Version, client.SettingType)
	schema.Formatter = setting.Formatter
}

func BusinessRoleTemplateBinding(schemas *types.Schemas, management *config.ScaledContext) {
	schema := schemas.Schema(&businessschema.Version, businessclient.BusinessRoleTemplateBindingType)
	schema.Validator = roletemplatebinding.NewCRTBValidator(management)
}

func HuaweiCloudApi(schemas *types.Schemas, managementContext *config.ScaledContext) {
	handler := huaweicloudapi.NewHandler(managementContext)

	schema := schemas.Schema(&businessschema.Version, businessclient.BusinessType)
	schema.Formatter = huaweicloudapi.Formatter
	schema.ActionHandler = handler.HuaweiCloudActionHandler
}

