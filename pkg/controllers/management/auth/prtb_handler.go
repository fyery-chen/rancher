package auth

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/types/apis/cloud.huawei.com/v3"
)

const (
	businessResource = "businesses"
	clusterResource           = "clusters"
	membershipBindingOwner    = "memberhsip-binding-owner"
	crtbInProjectBindingOwner = "crtb-in-project-binding-owner"
	prtbInClusterBindingOwner = "prtb-in-cluster-binding-owner"
	rbByOwnerIndex            = "auth.management.cattle.io/rb-by-owner"
	rbByRoleAndSubjectIndex   = "auth.management.cattle.io/crb-by-role-and-subject"
)

var businessManagmentPlaneResources = []string{"businessroletemplatebindings"}

type prtbLifecycle struct {
	mgr           *manager
	businessLister v3.BusinessLister
}

func (p *prtbLifecycle) Create(obj *v3.BusinessRoleTemplateBinding) (*v3.BusinessRoleTemplateBinding, error) {
	obj, err := p.reconcileSubject(obj)
	if err != nil {
		return nil, err
	}
	err = p.reconcileBindings(obj)
	return obj, err
}

func (p *prtbLifecycle) Updated(obj *v3.BusinessRoleTemplateBinding) (*v3.BusinessRoleTemplateBinding, error) {
	obj, err := p.reconcileSubject(obj)
	if err != nil {
		return nil, err
	}
	err = p.reconcileBindings(obj)
	return obj, err
}

func (p *prtbLifecycle) Remove(obj *v3.BusinessRoleTemplateBinding) (*v3.BusinessRoleTemplateBinding, error) {
	parts := strings.SplitN(obj.BusinessName, ":", 2)
	if len(parts) < 2 {
		return nil, errors.Errorf("cannot determine project and cluster from %v", obj.BusinessName)
	}
	clusterName := parts[0]
	err := p.mgr.reconcileBusinessMembershipBindingForDelete(clusterName, "", string(obj.UID))
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (p *prtbLifecycle) reconcileSubject(binding *v3.BusinessRoleTemplateBinding) (*v3.BusinessRoleTemplateBinding, error) {
	if binding.UserName != "" || binding.GroupName != "" || binding.GroupPrincipalName != "" {
		return binding, nil
	}

	if binding.UserPrincipalName != "" && binding.UserName == "" {
		displayName := binding.Annotations["auth.cattle.io/principal-display-name"]
		user, err := p.mgr.userMGR.EnsureUser(binding.UserPrincipalName, displayName)
		if err != nil {
			return binding, err
		}

		binding.UserName = user.Name
		return binding, nil
	}

	return nil, errors.Errorf("Binding %v has no subject", binding.Name)
}

// When a PRTB is created or updated, translate it into several k8s roles and bindings to actually enforce the RBAC.
// Specifically:
// - ensure the subject can see the project and its parent cluster in the mgmt API
// - if the subject was granted owner permissions for the project, ensure they can create/update/delete the project
// - if the subject was granted privileges to mgmt plane resources that are scoped to the project, enforce those rules in the project's mgmt plane namespace
func (p *prtbLifecycle) reconcileBindings(binding *v3.BusinessRoleTemplateBinding) error {
	if binding.UserName == "" && binding.GroupPrincipalName == "" && binding.GroupName == "" {
		return nil
	}

	businessName := binding.BusinessName
	business, err := p.businessLister.Get("", businessName)
	if err != nil {
		return err
	}
	if business == nil {
		return errors.Errorf("cannot create binding because business %v was not found", businessName)
	}

	isOwnerRole := binding.RoleTemplateName == "business-owner"
	var businessRoleName string
	if isOwnerRole {
		businessRoleName = strings.ToLower(fmt.Sprintf("%v-businessowner", businessName))
	} else {
		businessRoleName = strings.ToLower(fmt.Sprintf("%v-businessmember", businessName))
	}

	subject, err := buildSubjectFromRTB(binding)
	if err != nil {
		return err
	}

	if err := p.mgr.ensureBusinessMembershipBinding(businessRoleName, string(binding.UID), "", business, false, subject); err != nil {
		return err
	}
	//if err := p.mgr.grantManagementProjectScopedPrivilegesInClusterNamespace(binding.RoleTemplateName, proj.Namespace, prtbClusterManagmentPlaneResources, subject, binding); err != nil {
	//	return err
	//}
	return p.mgr.grantManagementPlanePrivileges(binding.RoleTemplateName, businessManagmentPlaneResources, subject, binding)
}
