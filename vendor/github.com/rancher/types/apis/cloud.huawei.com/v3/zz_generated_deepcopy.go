package v3

import (
	reflect "reflect"

	v1 "k8s.io/api/rbac/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
//
// Deprecated: deepcopy registration will go away when static deepcopy is fully implemented.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Action).DeepCopyInto(out.(*Action))
			return nil
		}, InType: reflect.TypeOf(&Action{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*AvailableZone).DeepCopyInto(out.(*AvailableZone))
			return nil
		}, InType: reflect.TypeOf(&AvailableZone{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Business).DeepCopyInto(out.(*Business))
			return nil
		}, InType: reflect.TypeOf(&Business{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessGlobalRole).DeepCopyInto(out.(*BusinessGlobalRole))
			return nil
		}, InType: reflect.TypeOf(&BusinessGlobalRole{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessGlobalRoleBinding).DeepCopyInto(out.(*BusinessGlobalRoleBinding))
			return nil
		}, InType: reflect.TypeOf(&BusinessGlobalRoleBinding{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessGlobalRoleBindingList).DeepCopyInto(out.(*BusinessGlobalRoleBindingList))
			return nil
		}, InType: reflect.TypeOf(&BusinessGlobalRoleBindingList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessGlobalRoleList).DeepCopyInto(out.(*BusinessGlobalRoleList))
			return nil
		}, InType: reflect.TypeOf(&BusinessGlobalRoleList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessList).DeepCopyInto(out.(*BusinessList))
			return nil
		}, InType: reflect.TypeOf(&BusinessList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessQuotaCheck).DeepCopyInto(out.(*BusinessQuotaCheck))
			return nil
		}, InType: reflect.TypeOf(&BusinessQuotaCheck{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessQuotaCheckOutput).DeepCopyInto(out.(*BusinessQuotaCheckOutput))
			return nil
		}, InType: reflect.TypeOf(&BusinessQuotaCheckOutput{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessQuotaSpec).DeepCopyInto(out.(*BusinessQuotaSpec))
			return nil
		}, InType: reflect.TypeOf(&BusinessQuotaSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessRoleTemplate).DeepCopyInto(out.(*BusinessRoleTemplate))
			return nil
		}, InType: reflect.TypeOf(&BusinessRoleTemplate{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessRoleTemplateBinding).DeepCopyInto(out.(*BusinessRoleTemplateBinding))
			return nil
		}, InType: reflect.TypeOf(&BusinessRoleTemplateBinding{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessRoleTemplateBindingList).DeepCopyInto(out.(*BusinessRoleTemplateBindingList))
			return nil
		}, InType: reflect.TypeOf(&BusinessRoleTemplateBindingList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessRoleTemplateList).DeepCopyInto(out.(*BusinessRoleTemplateList))
			return nil
		}, InType: reflect.TypeOf(&BusinessRoleTemplateList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BusinessStatus).DeepCopyInto(out.(*BusinessStatus))
			return nil
		}, InType: reflect.TypeOf(&BusinessStatus{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*DynamicSchema).DeepCopyInto(out.(*DynamicSchema))
			return nil
		}, InType: reflect.TypeOf(&DynamicSchema{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*DynamicSchemaList).DeepCopyInto(out.(*DynamicSchemaList))
			return nil
		}, InType: reflect.TypeOf(&DynamicSchemaList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*DynamicSchemaSpec).DeepCopyInto(out.(*DynamicSchemaSpec))
			return nil
		}, InType: reflect.TypeOf(&DynamicSchemaSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*DynamicSchemaStatus).DeepCopyInto(out.(*DynamicSchemaStatus))
			return nil
		}, InType: reflect.TypeOf(&DynamicSchemaStatus{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Field).DeepCopyInto(out.(*Field))
			return nil
		}, InType: reflect.TypeOf(&Field{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Filter).DeepCopyInto(out.(*Filter))
			return nil
		}, InType: reflect.TypeOf(&Filter{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*HuaweiCloudAccount).DeepCopyInto(out.(*HuaweiCloudAccount))
			return nil
		}, InType: reflect.TypeOf(&HuaweiCloudAccount{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*HuaweiCloudAccountList).DeepCopyInto(out.(*HuaweiCloudAccountList))
			return nil
		}, InType: reflect.TypeOf(&HuaweiCloudAccountList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*HuaweiCloudAccountSpec).DeepCopyInto(out.(*HuaweiCloudAccountSpec))
			return nil
		}, InType: reflect.TypeOf(&HuaweiCloudAccountSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*HuaweiCloudApiInformationInput).DeepCopyInto(out.(*HuaweiCloudApiInformationInput))
			return nil
		}, InType: reflect.TypeOf(&HuaweiCloudApiInformationInput{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*HuaweiCloudApiInformationOutput).DeepCopyInto(out.(*HuaweiCloudApiInformationOutput))
			return nil
		}, InType: reflect.TypeOf(&HuaweiCloudApiInformationOutput{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ListOpts).DeepCopyInto(out.(*ListOpts))
			return nil
		}, InType: reflect.TypeOf(&ListOpts{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*NodeFlavor).DeepCopyInto(out.(*NodeFlavor))
			return nil
		}, InType: reflect.TypeOf(&NodeFlavor{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*SubnetInfo).DeepCopyInto(out.(*SubnetInfo))
			return nil
		}, InType: reflect.TypeOf(&SubnetInfo{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Values).DeepCopyInto(out.(*Values))
			return nil
		}, InType: reflect.TypeOf(&Values{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*VpcInfo).DeepCopyInto(out.(*VpcInfo))
			return nil
		}, InType: reflect.TypeOf(&VpcInfo{})},
	)
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Action) DeepCopyInto(out *Action) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Action.
func (in *Action) DeepCopy() *Action {
	if in == nil {
		return nil
	}
	out := new(Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailableZone) DeepCopyInto(out *AvailableZone) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailableZone.
func (in *AvailableZone) DeepCopy() *AvailableZone {
	if in == nil {
		return nil
	}
	out := new(AvailableZone)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Business) DeepCopyInto(out *Business) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Business.
func (in *Business) DeepCopy() *Business {
	if in == nil {
		return nil
	}
	out := new(Business)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Business) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessGlobalRole) DeepCopyInto(out *BusinessGlobalRole) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]v1.PolicyRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessGlobalRole.
func (in *BusinessGlobalRole) DeepCopy() *BusinessGlobalRole {
	if in == nil {
		return nil
	}
	out := new(BusinessGlobalRole)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessGlobalRole) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessGlobalRoleBinding) DeepCopyInto(out *BusinessGlobalRoleBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessGlobalRoleBinding.
func (in *BusinessGlobalRoleBinding) DeepCopy() *BusinessGlobalRoleBinding {
	if in == nil {
		return nil
	}
	out := new(BusinessGlobalRoleBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessGlobalRoleBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessGlobalRoleBindingList) DeepCopyInto(out *BusinessGlobalRoleBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BusinessGlobalRoleBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessGlobalRoleBindingList.
func (in *BusinessGlobalRoleBindingList) DeepCopy() *BusinessGlobalRoleBindingList {
	if in == nil {
		return nil
	}
	out := new(BusinessGlobalRoleBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessGlobalRoleBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessGlobalRoleList) DeepCopyInto(out *BusinessGlobalRoleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BusinessGlobalRole, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessGlobalRoleList.
func (in *BusinessGlobalRoleList) DeepCopy() *BusinessGlobalRoleList {
	if in == nil {
		return nil
	}
	out := new(BusinessGlobalRoleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessGlobalRoleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessList) DeepCopyInto(out *BusinessList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Business, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessList.
func (in *BusinessList) DeepCopy() *BusinessList {
	if in == nil {
		return nil
	}
	out := new(BusinessList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessQuotaCheck) DeepCopyInto(out *BusinessQuotaCheck) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessQuotaCheck.
func (in *BusinessQuotaCheck) DeepCopy() *BusinessQuotaCheck {
	if in == nil {
		return nil
	}
	out := new(BusinessQuotaCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessQuotaCheckOutput) DeepCopyInto(out *BusinessQuotaCheckOutput) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessQuotaCheckOutput.
func (in *BusinessQuotaCheckOutput) DeepCopy() *BusinessQuotaCheckOutput {
	if in == nil {
		return nil
	}
	out := new(BusinessQuotaCheckOutput)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessQuotaSpec) DeepCopyInto(out *BusinessQuotaSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessQuotaSpec.
func (in *BusinessQuotaSpec) DeepCopy() *BusinessQuotaSpec {
	if in == nil {
		return nil
	}
	out := new(BusinessQuotaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessRoleTemplate) DeepCopyInto(out *BusinessRoleTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]v1.PolicyRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RoleTemplateNames != nil {
		in, out := &in.RoleTemplateNames, &out.RoleTemplateNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessRoleTemplate.
func (in *BusinessRoleTemplate) DeepCopy() *BusinessRoleTemplate {
	if in == nil {
		return nil
	}
	out := new(BusinessRoleTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessRoleTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessRoleTemplateBinding) DeepCopyInto(out *BusinessRoleTemplateBinding) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessRoleTemplateBinding.
func (in *BusinessRoleTemplateBinding) DeepCopy() *BusinessRoleTemplateBinding {
	if in == nil {
		return nil
	}
	out := new(BusinessRoleTemplateBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessRoleTemplateBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessRoleTemplateBindingList) DeepCopyInto(out *BusinessRoleTemplateBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BusinessRoleTemplateBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessRoleTemplateBindingList.
func (in *BusinessRoleTemplateBindingList) DeepCopy() *BusinessRoleTemplateBindingList {
	if in == nil {
		return nil
	}
	out := new(BusinessRoleTemplateBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessRoleTemplateBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessRoleTemplateList) DeepCopyInto(out *BusinessRoleTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BusinessRoleTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessRoleTemplateList.
func (in *BusinessRoleTemplateList) DeepCopy() *BusinessRoleTemplateList {
	if in == nil {
		return nil
	}
	out := new(BusinessRoleTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessRoleTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessStatus) DeepCopyInto(out *BusinessStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessStatus.
func (in *BusinessStatus) DeepCopy() *BusinessStatus {
	if in == nil {
		return nil
	}
	out := new(BusinessStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicSchema) DeepCopyInto(out *DynamicSchema) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicSchema.
func (in *DynamicSchema) DeepCopy() *DynamicSchema {
	if in == nil {
		return nil
	}
	out := new(DynamicSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DynamicSchema) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicSchemaList) DeepCopyInto(out *DynamicSchemaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DynamicSchema, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicSchemaList.
func (in *DynamicSchemaList) DeepCopy() *DynamicSchemaList {
	if in == nil {
		return nil
	}
	out := new(DynamicSchemaList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DynamicSchemaList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicSchemaSpec) DeepCopyInto(out *DynamicSchemaSpec) {
	*out = *in
	if in.ResourceMethods != nil {
		in, out := &in.ResourceMethods, &out.ResourceMethods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ResourceFields != nil {
		in, out := &in.ResourceFields, &out.ResourceFields
		*out = make(map[string]Field, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.ResourceActions != nil {
		in, out := &in.ResourceActions, &out.ResourceActions
		*out = make(map[string]Action, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.CollectionMethods != nil {
		in, out := &in.CollectionMethods, &out.CollectionMethods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.CollectionFields != nil {
		in, out := &in.CollectionFields, &out.CollectionFields
		*out = make(map[string]Field, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.CollectionActions != nil {
		in, out := &in.CollectionActions, &out.CollectionActions
		*out = make(map[string]Action, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.CollectionFilters != nil {
		in, out := &in.CollectionFilters, &out.CollectionFilters
		*out = make(map[string]Filter, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.IncludeableLinks != nil {
		in, out := &in.IncludeableLinks, &out.IncludeableLinks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicSchemaSpec.
func (in *DynamicSchemaSpec) DeepCopy() *DynamicSchemaSpec {
	if in == nil {
		return nil
	}
	out := new(DynamicSchemaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicSchemaStatus) DeepCopyInto(out *DynamicSchemaStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicSchemaStatus.
func (in *DynamicSchemaStatus) DeepCopy() *DynamicSchemaStatus {
	if in == nil {
		return nil
	}
	out := new(DynamicSchemaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Field) DeepCopyInto(out *Field) {
	*out = *in
	in.Default.DeepCopyInto(&out.Default)
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Field.
func (in *Field) DeepCopy() *Field {
	if in == nil {
		return nil
	}
	out := new(Field)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Filter) DeepCopyInto(out *Filter) {
	*out = *in
	if in.Modifiers != nil {
		in, out := &in.Modifiers, &out.Modifiers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Filter.
func (in *Filter) DeepCopy() *Filter {
	if in == nil {
		return nil
	}
	out := new(Filter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HuaweiCloudAccount) DeepCopyInto(out *HuaweiCloudAccount) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HuaweiCloudAccount.
func (in *HuaweiCloudAccount) DeepCopy() *HuaweiCloudAccount {
	if in == nil {
		return nil
	}
	out := new(HuaweiCloudAccount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HuaweiCloudAccount) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HuaweiCloudAccountList) DeepCopyInto(out *HuaweiCloudAccountList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HuaweiCloudAccount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HuaweiCloudAccountList.
func (in *HuaweiCloudAccountList) DeepCopy() *HuaweiCloudAccountList {
	if in == nil {
		return nil
	}
	out := new(HuaweiCloudAccountList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HuaweiCloudAccountList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HuaweiCloudAccountSpec) DeepCopyInto(out *HuaweiCloudAccountSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HuaweiCloudAccountSpec.
func (in *HuaweiCloudAccountSpec) DeepCopy() *HuaweiCloudAccountSpec {
	if in == nil {
		return nil
	}
	out := new(HuaweiCloudAccountSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HuaweiCloudApiInformationInput) DeepCopyInto(out *HuaweiCloudApiInformationInput) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HuaweiCloudApiInformationInput.
func (in *HuaweiCloudApiInformationInput) DeepCopy() *HuaweiCloudApiInformationInput {
	if in == nil {
		return nil
	}
	out := new(HuaweiCloudApiInformationInput)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HuaweiCloudApiInformationOutput) DeepCopyInto(out *HuaweiCloudApiInformationOutput) {
	*out = *in
	if in.VpcInfo != nil {
		in, out := &in.VpcInfo, &out.VpcInfo
		*out = make([]VpcInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SshKeyName != nil {
		in, out := &in.SshKeyName, &out.SshKeyName
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NodeFlavor != nil {
		in, out := &in.NodeFlavor, &out.NodeFlavor
		*out = make([]NodeFlavor, len(*in))
		copy(*out, *in)
	}
	if in.AvailableZone != nil {
		in, out := &in.AvailableZone, &out.AvailableZone
		*out = make([]AvailableZone, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HuaweiCloudApiInformationOutput.
func (in *HuaweiCloudApiInformationOutput) DeepCopy() *HuaweiCloudApiInformationOutput {
	if in == nil {
		return nil
	}
	out := new(HuaweiCloudApiInformationOutput)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ListOpts) DeepCopyInto(out *ListOpts) {
	*out = *in
	if in.Filters != nil {
		in, out := &in.Filters, &out.Filters
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListOpts.
func (in *ListOpts) DeepCopy() *ListOpts {
	if in == nil {
		return nil
	}
	out := new(ListOpts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeFlavor) DeepCopyInto(out *NodeFlavor) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeFlavor.
func (in *NodeFlavor) DeepCopy() *NodeFlavor {
	if in == nil {
		return nil
	}
	out := new(NodeFlavor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubnetInfo) DeepCopyInto(out *SubnetInfo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubnetInfo.
func (in *SubnetInfo) DeepCopy() *SubnetInfo {
	if in == nil {
		return nil
	}
	out := new(SubnetInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Values) DeepCopyInto(out *Values) {
	*out = *in
	if in.StringSliceValue != nil {
		in, out := &in.StringSliceValue, &out.StringSliceValue
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Values.
func (in *Values) DeepCopy() *Values {
	if in == nil {
		return nil
	}
	out := new(Values)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VpcInfo) DeepCopyInto(out *VpcInfo) {
	*out = *in
	if in.SubnetInfo != nil {
		in, out := &in.SubnetInfo, &out.SubnetInfo
		*out = make([]SubnetInfo, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VpcInfo.
func (in *VpcInfo) DeepCopy() *VpcInfo {
	if in == nil {
		return nil
	}
	out := new(VpcInfo)
	in.DeepCopyInto(out)
	return out
}
