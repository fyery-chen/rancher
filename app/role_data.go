package app

import (
	"github.com/pkg/errors"
	"github.com/rancher/types/apis/cloud.huawei.com/v3"
	managementv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"golang.org/x/crypto/bcrypt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var defaultAdminLabel = map[string]string{"authz.management.cattle.io/bootstrapping": "admin-user"}

func addRoles(management *config.ManagementContext) (string, error) {
	rb := newRoleBuilder()

	rb.addRole("Create Business", "businesses-create").addRule().apiGroups("cloud.huawei.com").resources("businesses").verbs("create").
		addRule().apiGroups("management.cattle.io").resources("clusters").verbs("create")

	// TODO user should be dynamically authorized to only see herself
	// TODO Need "self-service" for nodetemplates such that a user can create them, but only RUD their own
	// TODO enable when groups are "in". they need to be self-service

	if err := rb.reconcileGlobalRoles(management); err != nil {
		return "", errors.Wrap(err, "problem reconciling globl roles")
	}

	// RoleTemplates to be used inside of clusters
	rb = newRoleBuilder()

	// Cluster roles
	rb.addRoleTemplate("Business Owner", "business-owner", "business", true, false, false).
		addRule().apiGroups("*").resources("*").verbs("*").
		addRule().apiGroups().nonResourceURLs("*").verbs("*")

	rb.addRoleTemplate("Business Member", "business-member", "businesses", true, false, false).
		addRule().apiGroups("cloud.huawei.com").resources("businessroletemplatebindings").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("clusterroletemplatebindings").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("projects").verbs("create").
		addRule().apiGroups("management.cattle.io").resources("nodes", "nodepools").verbs("get", "list", "watch").
		addRule().apiGroups("*").resources("nodes").verbs("get", "list", "watch").
		addRule().apiGroups("*").resources("persistentvolumes").verbs("get", "list", "watch").
		addRule().apiGroups("*").resources("storageclasses").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("clusterevents").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("clusterpipelines").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("clusterloggings").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("clusteralerts").verbs("get", "list", "watch").
		addRule().apiGroups("management.cattle.io").resources("notifiers").verbs("get", "list", "watch")

	//rb.addRoleTemplate("Create Clusters", "clusters-create", "business", true, false, false).
	//	addRule().apiGroups("management.cattle.io").resources("clusters").verbs("create")
	//
	//rb.addRoleTemplate("View All Clusters", "cluster-view", "business", true, false, false).
	//	addRule().apiGroups("management.cattle.io").resources("clusterroletemplatebindings").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("nodes", "nodepools").verbs("get", "list", "watch").
	//	addRule().apiGroups("*").resources("nodes").verbs("get", "list", "watch").
	//	addRule().apiGroups("*").resources("persistentvolumes").verbs("get", "list", "watch").
	//	addRule().apiGroups("*").resources("storageclasses").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("clusterevents").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("clusterpipelines").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("clusterloggings").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("clusteralerts").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("notifiers").verbs("get", "list", "watch")

	// Not specific to project or cluster
	// TODO When clusterevents has value, consider adding this back in
	//rb.addRoleTemplate("View Events", "events-view", "", true, false, false).
	//	addRule().apiGroups("*").resources("events").verbs("get", "list", "watch").
	//	addRule().apiGroups("management.cattle.io").resources("clusterevents").verbs("get", "list", "watch")

	if err := rb.reconcileRoleTemplates(management); err != nil {
		return "", errors.Wrap(err, "problem reconciling role templates")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	set := labels.Set(defaultAdminLabel)
	admins, err := management.Management.Users("").List(v1.ListOptions{LabelSelector: set.String()})
	if err != nil {
		return "", err
	}

	// TODO This logic is going to be a problem in an HA setup because a race will cause more than one admin user to be created
	var admin *managementv3.User
	if len(admins.Items) == 0 {
		admin, err = management.Management.Users("").Create(&managementv3.User{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "user-",
				Labels:       defaultAdminLabel,
			},
			DisplayName:        "Default Admin",
			Username:           "admin",
			Password:           string(hash),
			MustChangePassword: true,
		})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return "", errors.Wrap(err, "can not ensure admin user exists")
		}

	} else {
		admin = &admins.Items[0]
	}

	bindings, err := management.Business.BusinessGlobalRoleBindings("").List(v1.ListOptions{LabelSelector: set.String()})
	if err != nil {
		return "", err
	}
	if len(bindings.Items) == 0 {
		management.Business.BusinessGlobalRoleBindings("").Create(
			&v3.BusinessGlobalRoleBinding{
				ObjectMeta: v1.ObjectMeta{
					GenerateName: "globalrolebinding-",
					Labels:       defaultAdminLabel,
				},
				UserName:       admin.Name,
				GlobalRoleName: "admin",
			})
	}

	return admin.Name, nil
}
