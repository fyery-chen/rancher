package auth

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/types/slice"
	v13 "github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/cloud.huawei.com/v3"
	typesrbacv1 "github.com/rancher/types/apis/rbac.authorization.k8s.io/v1"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

func newRTBLifecycles(management *config.ManagementContext) (*prtbLifecycle) {
	crbInformer := management.RBAC.ClusterRoleBindings("").Controller().Informer()
	crbIndexers := map[string]cache.IndexFunc{
		rbByRoleAndSubjectIndex: rbByRoleAndSubject,
	}
	crbInformer.AddIndexers(crbIndexers)

	rbInformer := management.RBAC.RoleBindings("").Controller().Informer()
	rbIndexers := map[string]cache.IndexFunc{
		rbByOwnerIndex:          rbByOwner,
		rbByRoleAndSubjectIndex: rbByRoleAndSubject,
	}
	rbInformer.AddIndexers(rbIndexers)

	mgr := &manager{
		mgmt:          management,
		crbLister:     management.RBAC.ClusterRoleBindings("").Controller().Lister(),
		crLister:      management.RBAC.ClusterRoles("").Controller().Lister(),
		rLister:       management.RBAC.Roles("").Controller().Lister(),
		rbLister:      management.RBAC.RoleBindings("").Controller().Lister(),
		rtLister:      management.Business.BusinessRoleTemplates("").Controller().Lister(),
		nsLister:      management.Core.Namespaces("").Controller().Lister(),
		rbIndexer:     rbInformer.GetIndexer(),
		crbIndexer:    crbInformer.GetIndexer(),
		userMGR:       management.UserManager,
	}
	prtb := &prtbLifecycle{
		mgr:           mgr,
		businessLister: management.Business.Businesses("").Controller().Lister(),
	}
	return prtb
}

type manager struct {
	crLister      typesrbacv1.ClusterRoleLister
	rLister       typesrbacv1.RoleLister
	rbLister      typesrbacv1.RoleBindingLister
	crbLister     typesrbacv1.ClusterRoleBindingLister
	rtLister      v3.BusinessRoleTemplateLister
	nsLister      v13.NamespaceLister
	rbIndexer     cache.Indexer
	crbIndexer    cache.Indexer
	mgmt          *config.ManagementContext
	userMGR       user.Manager
}

// When a PRTB is created that gives a subject some permissions in a project or cluster, we need to create a "membership" binding
// that gives the subject access to the the project/cluster custom resource itself
func (m *manager) ensureBusinessMembershipBinding(roleName, rtbUID, namespace string, business *v3.Business, makeOwner bool, subject v1.Subject) error {
	if err := m.createBusinessMembershipRole(roleName, business, makeOwner); err != nil {
		return err
	}

	key := rbRoleSubjectKey(roleName, subject)
	set := labels.Set(map[string]string{rtbUID: membershipBindingOwner})
	rbs, err := m.rbLister.List("", set.AsSelector())
	if err != nil {
		return err
	}
	var rb *v1.RoleBinding
	for _, iRB := range rbs {
		if len(iRB.Subjects) != 1 {
			iKey := rbRoleSubjectKey(iRB.RoleRef.Name, iRB.Subjects[0])
			if iKey == key {
				rb = iRB
				continue
			}
		}
		if err := m.reconcileBusinessMembershipBindingForDelete(namespace, roleName, rtbUID); err != nil {
			return err
		}
	}

	if rb != nil {
		return nil
	}

	objs, err := m.crbIndexer.ByIndex(rbByRoleAndSubjectIndex, key)
	if err != nil {
		return err
	}

	if len(objs) == 0 {
		_, err := m.mgmt.RBAC.RoleBindings(namespace).Create(&v1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "rolebinding-",
				Labels: map[string]string{
					rtbUID: membershipBindingOwner,
				},
			},
			Subjects: []v1.Subject{subject},
			RoleRef: v1.RoleRef{
				Kind: "Role",
				Name: roleName,
			},
		})
		return err
	}

	rb, _ = objs[0].(*v1.RoleBinding)
	for owner := range rb.Labels {
		if rtbUID == owner {
			return nil
		}
	}

	rb = rb.DeepCopy()
	if rb.Labels == nil {
		rb.Labels = map[string]string{}
	}
	rb.Labels[rtbUID] = membershipBindingOwner
	_, err = m.mgmt.RBAC.RoleBindings(namespace).Update(rb)
	return err
}

func (m *manager) createBusinessMembershipRole(roleName string, business *v3.Business, makeOwner bool) error {
	if cr, _ := m.rLister.Get("", roleName); cr == nil {
		return m.createMembershipRole(businessResource, roleName, makeOwner, business, m.mgmt.RBAC.Roles("").ObjectClient())
	}
	return nil
}

func (m *manager) createMembershipRole(resourceType, roleName string, makeOwner bool, ownerObject interface{}, client *objectclient.ObjectClient) error {
	metaObj, err := meta.Accessor(ownerObject)
	if err != nil {
		return err
	}
	typeMeta, err := meta.TypeAccessor(ownerObject)
	if err != nil {
		return err
	}
	rules := []v1.PolicyRule{
		{
			APIGroups:     []string{"cloud.huawei.com"},
			Resources:     []string{resourceType},
			ResourceNames: []string{metaObj.GetName()},
			Verbs:         []string{"get"},
		},
	}

	if makeOwner {
		rules[0].Verbs = []string{"*"}
	} else {
		rules[0].Verbs = []string{"get"}
	}
	_, err = client.Create(&v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleName,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: typeMeta.GetAPIVersion(),
					Kind:       typeMeta.GetKind(),
					Name:       metaObj.GetName(),
					UID:        metaObj.GetUID(),
				},
			},
		},
		Rules: rules,
	})
	return err
}

// The PRTB has been deleted, either delete or update the project membership binding so that the subject
// is removed from the project if they should be
func (m *manager) reconcileBusinessMembershipBindingForDelete(namespace, roleToKeep, rtbUID string) error {
	list := func(ns string, selector labels.Selector) ([]runtime.Object, error) {
		rbs, err := m.rbLister.List(ns, selector)
		if err != nil {
			return nil, err
		}

		var items []runtime.Object
		for _, rb := range rbs {
			items = append(items, rb.DeepCopy())
		}
		return items, nil
	}

	convert := func(i interface{}) string {
		rb, _ := i.(*v1.RoleBinding)
		return rb.RoleRef.Name
	}

	return m.reconcileMembershipBindingForDelete("", roleToKeep, rtbUID, list, convert, m.mgmt.RBAC.RoleBindings("").ObjectClient())
}

type listFn func(ns string, selector labels.Selector) ([]runtime.Object, error)
type convertFn func(i interface{}) string

func (m *manager) reconcileMembershipBindingForDelete(namespace, roleToKeep, rtbUID string, list listFn, convert convertFn, client *objectclient.ObjectClient) error {
	set := labels.Set(map[string]string{rtbUID: membershipBindingOwner})
	roleBindings, err := list(namespace, set.AsSelector())
	if err != nil {
		return err
	}

	for _, rb := range roleBindings {
		objMeta, err := meta.Accessor(rb)
		if err != nil {
			return err
		}

		roleName := convert(rb)
		if roleName == roleToKeep {
			continue
		}

		for k, v := range objMeta.GetLabels() {
			if k == rtbUID && v == membershipBindingOwner {
				delete(objMeta.GetLabels(), k)
			}
		}

		if len(objMeta.GetLabels()) == 0 {
			if err := client.Delete(objMeta.GetName(), &metav1.DeleteOptions{}); err != nil {
				if apierrors.IsNotFound(err) {
					continue
				}
				return err
			}
		} else {
			if _, err := client.Update(objMeta.GetName(), rb); err != nil {
				return err
			}
		}
	}

	return nil
}

// Certain resources (projects, machines, prtbs, crtbs, clusterevents, etc) exist in the mangement plane but are scoped to clusters or
// projects. They need special RBAC handling because the need to be authorized just inside of the namespace that backs the project
// or cluster they belong to.
func (m *manager) grantManagementPlanePrivileges(roleTemplateName string, resources []string, subject v1.Subject, binding interface{}) error {
	bindingMeta, err := meta.Accessor(binding)
	if err != nil {
		return err
	}
	bindingTypeMeta, err := meta.TypeAccessor(binding)
	if err != nil {
		return err
	}
	namespace := bindingMeta.GetNamespace()

	roles, err := m.gatherAndDedupeRoles(roleTemplateName)
	if err != nil {
		return err
	}

	desiredRBs := map[string]*v1.RoleBinding{}
	roleBindings := m.mgmt.RBAC.RoleBindings(namespace)
	for _, role := range roles {
		for _, resource := range resources {
			verbs, err := m.checkForManagementPlaneRules(role, resource)
			if err != nil {
				return err
			}
			if len(verbs) > 0 {
				if err := m.reconcileManagementPlaneRole(namespace, resource, role, verbs); err != nil {
					return err
				}

				bindingName := bindingMeta.GetName() + "-" + role.Name
				if _, ok := desiredRBs[bindingName]; !ok {
					desiredRBs[bindingName] = &v1.RoleBinding{
						ObjectMeta: metav1.ObjectMeta{
							Name: bindingName,
							OwnerReferences: []metav1.OwnerReference{
								{
									APIVersion: bindingTypeMeta.GetAPIVersion(),
									Kind:       bindingTypeMeta.GetKind(),
									Name:       bindingMeta.GetName(),
									UID:        bindingMeta.GetUID(),
								},
							},
						},
						Subjects: []v1.Subject{subject},
						RoleRef: v1.RoleRef{
							Kind: "Role",
							Name: role.Name,
						},
					}
				}
			}
		}
	}

	currentRBs := map[string]*v1.RoleBinding{}
	current, err := m.rbIndexer.ByIndex(rbByOwnerIndex, string(bindingMeta.GetUID()))
	if err != nil {
		return err
	}
	for _, c := range current {
		rb := c.(*v1.RoleBinding)
		currentRBs[rb.Name] = rb
	}

	return m.reconcileDesiredMGMTPlaneRoleBindings(currentRBs, desiredRBs, roleBindings)
}

// grantManagementClusterScopedPrivilegesInProjectNamespace ensures that rolebindings for roles like cluster-owner (that should be able to fully
// manage all projects in a cluster) grant proper permissions to project-scoped resources. Specifically, this satisfies the use case that
// a cluster owner should be able to manage the members of all projects in their cluster
func (m *manager) grantManagementClusterScopedPrivilegesInProjectNamespace(roleTemplateName, projectNamespace string, resources []string,
	subject v1.Subject, binding *v3.BusinessRoleTemplateBinding) error {
	roles, err := m.gatherAndDedupeRoles(roleTemplateName)
	if err != nil {
		return err
	}

	desiredRBs := map[string]*v1.RoleBinding{}
	roleBindings := m.mgmt.RBAC.RoleBindings(projectNamespace)
	for _, role := range roles {
		for _, resource := range resources {
			verbs, err := m.checkForManagementPlaneRules(role, resource)
			if err != nil {
				return err
			}
			if len(verbs) > 0 {
				if err := m.reconcileManagementPlaneRole(projectNamespace, resource, role, verbs); err != nil {
					return err
				}

				bindingName := binding.Name + "-" + role.Name
				if _, ok := desiredRBs[bindingName]; !ok {
					desiredRBs[bindingName] = &v1.RoleBinding{
						ObjectMeta: metav1.ObjectMeta{
							Name: bindingName,
							Labels: map[string]string{
								string(binding.UID): crtbInProjectBindingOwner,
							},
						},
						Subjects: []v1.Subject{subject},
						RoleRef: v1.RoleRef{
							Kind: "Role",
							Name: role.Name,
						},
					}
				}
			}
		}
	}

	currentRBs := map[string]*v1.RoleBinding{}
	set := labels.Set(map[string]string{string(binding.UID): crtbInProjectBindingOwner})
	current, err := m.rbLister.List(projectNamespace, set.AsSelector())
	if err != nil {
		return err
	}
	for _, rb := range current {
		currentRBs[rb.Name] = rb
	}

	return m.reconcileDesiredMGMTPlaneRoleBindings(currentRBs, desiredRBs, roleBindings)
}

// grantManagementProjectScopedPrivilegesInClusterNamespace ensures that project roles grant permissions to certain cluster-scoped
// resources(notifier, clusterpipelines). These resources exists in cluster namespace but need to be shared between projects.
func (m *manager) grantManagementProjectScopedPrivilegesInClusterNamespace(roleTemplateName, clusterNamespace string, resources []string,
	subject v1.Subject, binding *v3.BusinessRoleTemplateBinding) error {
	roles, err := m.gatherAndDedupeRoles(roleTemplateName)
	if err != nil {
		return err
	}

	desiredRBs := map[string]*v1.RoleBinding{}
	roleBindings := m.mgmt.RBAC.RoleBindings(clusterNamespace)
	for _, role := range roles {
		for _, resource := range resources {
			verbs, err := m.checkForManagementPlaneRules(role, resource)
			if err != nil {
				return err
			}
			if len(verbs) > 0 {
				if err := m.reconcileManagementPlaneRole(clusterNamespace, resource, role, verbs); err != nil {
					return err
				}

				bindingName := binding.Name + "-" + role.Name
				if _, ok := desiredRBs[bindingName]; !ok {
					desiredRBs[bindingName] = &v1.RoleBinding{
						ObjectMeta: metav1.ObjectMeta{
							Name: bindingName,
							Labels: map[string]string{
								string(binding.UID): prtbInClusterBindingOwner,
							},
						},
						Subjects: []v1.Subject{subject},
						RoleRef: v1.RoleRef{
							Kind: "Role",
							Name: role.Name,
						},
					}
				}
			}
		}
	}

	currentRBs := map[string]*v1.RoleBinding{}
	set := labels.Set(map[string]string{string(binding.UID): prtbInClusterBindingOwner})
	current, err := m.rbLister.List(clusterNamespace, set.AsSelector())
	if err != nil {
		return err
	}
	for _, rb := range current {
		currentRBs[rb.Name] = rb
	}

	return m.reconcileDesiredMGMTPlaneRoleBindings(currentRBs, desiredRBs, roleBindings)
}

func (m *manager) gatherAndDedupeRoles(roleTemplateName string) (map[string]*v3.BusinessRoleTemplate, error) {
	rt, err := m.rtLister.Get("", roleTemplateName)
	if err != nil {
		return nil, err
	}
	allRoles := map[string]*v3.BusinessRoleTemplate{}
	if err := m.gatherRoleTemplates(rt, allRoles); err != nil {
		return nil, err
	}

	//de-dupe
	roles := map[string]*v3.BusinessRoleTemplate{}
	for _, role := range allRoles {
		roles[role.Name] = role
	}
	return roles, nil
}

func (m *manager) reconcileDesiredMGMTPlaneRoleBindings(currentRBs, desiredRBs map[string]*v1.RoleBinding, roleBindings typesrbacv1.RoleBindingInterface) error {
	rbsToDelete := map[string]bool{}
	processed := map[string]bool{}
	for _, rb := range currentRBs {
		// protect against an rb being in the list more than once (shouldn't happen, but just to be safe)
		if ok := processed[rb.Name]; ok {
			continue
		}
		processed[rb.Name] = true

		if _, ok := desiredRBs[rb.Name]; ok {
			delete(desiredRBs, rb.Name)
		} else {
			rbsToDelete[rb.Name] = true
		}
	}

	for _, rb := range desiredRBs {
		_, err := roleBindings.Create(rb)
		if err != nil {
			return err
		}
	}

	for name := range rbsToDelete {
		if err := roleBindings.Delete(name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}

// If the roleTemplate has rules granting access to a management plane resource, return the verbs for those rules
func (m *manager) checkForManagementPlaneRules(role *v3.BusinessRoleTemplate, managementPlaneResource string) (map[string]bool, error) {
	var rules []v1.PolicyRule
	if role.External {
		externalRole, err := m.crLister.Get("", role.Name)
		if err != nil && !apierrors.IsNotFound(err) {
			// dont error if it doesnt exist
			return nil, err
		}
		if externalRole != nil {
			rules = externalRole.Rules
		}
	} else {
		rules = role.Rules
	}

	verbs := map[string]bool{}
	for _, rule := range rules {
		if (slice.ContainsString(rule.Resources, managementPlaneResource) || slice.ContainsString(rule.Resources, "*")) && len(rule.ResourceNames) == 0 {
			for _, v := range rule.Verbs {
				verbs[v] = true
			}
		}
	}

	return verbs, nil
}

func (m *manager) reconcileManagementPlaneRole(namespace, resource string, rt *v3.BusinessRoleTemplate, newVerbs map[string]bool) error {
	roleCli := m.mgmt.RBAC.Roles(namespace)
	if role, err := m.rLister.Get(namespace, rt.Name); err == nil && role != nil {
		currentVerbs := map[string]bool{}
		for _, rule := range role.Rules {
			if slice.ContainsString(rule.Resources, resource) {
				for _, v := range rule.Verbs {
					currentVerbs[v] = true
				}
			}
		}

		if !reflect.DeepEqual(currentVerbs, newVerbs) {
			role = role.DeepCopy()
			added := false
			for i, rule := range role.Rules {
				if slice.ContainsString(rule.Resources, resource) {
					role.Rules[i] = buildRule(resource, newVerbs)
					added = true
				}
			}
			if !added {
				role.Rules = append(role.Rules, buildRule(resource, newVerbs))
			}
			_, err := roleCli.Update(role)
			return err
		}
		return nil
	}

	rules := []v1.PolicyRule{buildRule(resource, newVerbs)}
	_, err := roleCli.Create(&v1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: rt.Name,
		},
		Rules: rules,
	})
	if err != nil {
		return errors.Wrapf(err, "couldn't create role %v", rt.Name)
	}

	return nil
}

func (m *manager) gatherRoleTemplates(rt *v3.BusinessRoleTemplate, roleTemplates map[string]*v3.BusinessRoleTemplate) error {
	roleTemplates[rt.Name] = rt

	for _, rtName := range rt.RoleTemplateNames {
		subRT, err := m.rtLister.Get("", rtName)
		if err != nil {
			return errors.Wrapf(err, "couldn't get RoleTemplate %s", rtName)
		}
		if err := m.gatherRoleTemplates(subRT, roleTemplates); err != nil {
			return errors.Wrapf(err, "couldn't gather RoleTemplate %s", rtName)
		}
	}

	return nil
}

func buildRule(resource string, verbs map[string]bool) v1.PolicyRule {
	var vs []string
	for v := range verbs {
		vs = append(vs, v)
	}
	return v1.PolicyRule{
		Resources: []string{resource},
		Verbs:     vs,
		APIGroups: []string{"*"},
	}
}

func buildSubjectFromRTB(binding interface{}) (v1.Subject, error) {
	var userName, groupPrincipalName, groupName, name, kind string
	if rtb, ok := binding.(*v3.BusinessRoleTemplateBinding); ok {
		userName = rtb.UserName
		groupPrincipalName = rtb.GroupPrincipalName
		groupName = rtb.GroupName
	} else {
		return v1.Subject{}, errors.Errorf("unrecognized roleTemplateBinding type: %v", binding)
	}


	if userName != "" {
		name = userName
		kind = "User"
	}

	if groupPrincipalName != "" {
		if name != "" {
			return v1.Subject{}, errors.Errorf("roletemplatebinding has more than one subject fields set: %v", binding)
		}
		name = groupPrincipalName
		kind = "Group"
	}

	if groupName != "" {
		if name != "" {
			return v1.Subject{}, errors.Errorf("roletemplatebinding has more than one subject fields set: %v", binding)
		}
		name = groupName
		kind = "Group"
	}

	if name == "" {
		return v1.Subject{}, errors.Errorf("roletemplatebinding doesn't have any subject fields set: %v", binding)
	}

	return v1.Subject{
		Kind: kind,
		Name: name,
	}, nil
}
