package auth

import (
	"context"

	"github.com/rancher/types/config"
)

func RegisterEarly(ctx context.Context, management *config.ManagementContext) {
	prtb := newRTBLifecycles(management)
	gr := newGlobalRoleLifecycle(management)
	grb := newGlobalRoleBindingLifecycle(management)
	p := newPandCLifecycles(management)
	u := newUserLifecycle(management)

	management.Business.BusinessRoleTemplateBindings("").AddLifecycle("mgmt-auth-prtb-controller", prtb)
	management.Business.BusinessGlobalRoles("").AddLifecycle("mgmt-auth-gr-controller", gr)
	management.Business.BusinessGlobalRoleBindings("").AddLifecycle("mgmt-auth-grb-controller", grb)
	management.Business.Businesses("").AddHandler("mgmt-cluster-rbac-delete", p.sync)
	management.Management.Users("").AddLifecycle("mgmt-auth-users-controller", u)
}

func RegisterLate(ctx context.Context, management *config.ManagementContext) {
	p := newPandCLifecycles(management)
	management.Business.Businesses("").AddLifecycle("mgmt-cluster-rbac-remove", p)
}
