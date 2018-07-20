package roletemplatebinding

import (
	"fmt"
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/config"
	"github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/client/cloud/v3"
)

func NewCRTBValidator(management *config.ScaledContext) types.Validator {
	return newValidator(management, client.BusinessRoleTemplateBindingFieldRoleTemplateId)
}

func newValidator(management *config.ScaledContext, field string) types.Validator {
	validator := &Validator{
		roleTemplateLister: management.Business.BusinessRoleTemplates("").Controller().Lister(),
		field:              field,
	}

	return validator.Validator
}

type Validator struct {
	roleTemplateLister v3.BusinessRoleTemplateLister
	field              string
}

func (v *Validator) Validator(request *types.APIContext, schema *types.Schema, data map[string]interface{}) error {
	return v.ValidateRoleTemplateBinding(data[v.field])
}

func (v *Validator) ValidateRoleTemplateBinding(obj interface{}) error {
	roleTemplateID, ok := obj.(string)
	if !ok {
		return httperror.NewAPIError(httperror.MissingRequired, "Request does not have a valid roleTemplateId")
	}

	roleTemplate, err := v.roleTemplateLister.Get("", roleTemplateID)
	if err != nil {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error getting role template: %v", err))
	}

	if roleTemplate.Locked {
		return httperror.NewAPIError(httperror.InvalidState, "Role is locked and cannot be assigned")
	}

	return nil
}
