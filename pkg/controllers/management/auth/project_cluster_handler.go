package auth

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/rancher/norman/condition"
	corev1 "github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	creatorIDAnn                  = "field.cattle.io/creatorId"
	creatorOwnerBindingAnnotation = "authz.management.cattle.io/creator-owner-binding"
)

var defaultProjectLabels = labels.Set(map[string]string{"authz.management.cattle.io/default-project": "true"})
var crtbCeatorOwnerAnnotations = map[string]string{creatorOwnerBindingAnnotation: "true"}

func newPandCLifecycles(management *config.ManagementContext) (*businessLifecycle) {
	m := &mgr{
		mgmt:          management,
		nsLister:      management.Core.Namespaces("").Controller().Lister(),
		prtbLister:    management.Business.BusinessRoleTemplateBindings("").Controller().Lister(),
		businessLister: management.Business.Businesses("").Controller().Lister(),
	}
	p := &businessLifecycle{
		mgr: m,
	}
	return p
}

type businessLifecycle struct {
	mgr *mgr
}

func (l *businessLifecycle) sync(key string, orig *v3.Business) error {
	if orig == nil {
		return nil
	}

	obj := orig.DeepCopyObject()

	obj, err := l.mgr.reconcileResourceToNamespace(obj)
	if err != nil {
		return err
	}

	obj, err = l.mgr.reconcileCreatorRTB(obj)
	if err != nil {
		return err
	}

	// update if it has changed
	if obj != nil && !reflect.DeepEqual(orig, obj) {
		_, err = l.mgr.mgmt.Business.Businesses("").ObjectClient().Update(orig.Name, obj)
		if err != nil {
			return err
		}
	}
	return err
}

func (l *businessLifecycle) Create(obj *v3.Business) (*v3.Business, error) {
	// no-op because the sync function will take care of it
	return obj, nil
}

func (l *businessLifecycle) Updated(obj *v3.Business) (*v3.Business, error) {
	// no-op because the sync function will take care of it
	return obj, nil
}

func (l *businessLifecycle) Remove(obj *v3.Business) (*v3.Business, error) {
	err := l.mgr.deleteNamespace(obj)
	return obj, err
}

type clusterLifecycle struct {
	mgr *mgr
}

type mgr struct {
	mgmt          *config.ManagementContext
	nsLister      corev1.NamespaceLister
	businessLister v3.BusinessLister
	prtbLister v3.BusinessRoleTemplateBindingLister
}

func (m *mgr) reconcileCreatorRTB(obj runtime.Object) (runtime.Object, error) {
	return v3.CreatorMadeOwner.DoUntilTrue(obj, func() (runtime.Object, error) {
		metaAccessor, err := meta.Accessor(obj)
		if err != nil {
			return obj, err
		}

		typeAccessor, err := meta.TypeAccessor(obj)
		if err != nil {
			return obj, err
		}

		creatorID, ok := metaAccessor.GetAnnotations()[creatorIDAnn]
		if !ok {
			logrus.Warnf("%v %v has no creatorId annotation. Cannot add creator as owner", typeAccessor.GetKind(), metaAccessor.GetName())
			return obj, nil
		}

		rtbName := "creator"
		om := v1.ObjectMeta{
			Name:      rtbName,
			Namespace: metaAccessor.GetName(),
		}

		switch typeAccessor.GetKind() {
		case v3.BusinessGroupVersionKind.Kind:
			if rtb, _ := m.prtbLister.Get(metaAccessor.GetName(), rtbName); rtb != nil {
				return obj, nil
			}
			if _, err := m.mgmt.Business.BusinessRoleTemplateBindings(metaAccessor.GetName()).Create(&v3.BusinessRoleTemplateBinding{
				ObjectMeta:       om,
				BusinessName:      metaAccessor.GetNamespace() + ":" + metaAccessor.GetName(),
				RoleTemplateName: "project-owner",
				UserName:         creatorID,
			}); err != nil {
				return obj, err
			}
		}

		return obj, nil
	})
}

func (m *mgr) deleteNamespace(obj runtime.Object) error {
	o, err := meta.Accessor(obj)
	if err != nil {
		return condition.Error("MissingMetadata", err)
	}

	nsClient := m.mgmt.K8sClient.CoreV1().Namespaces()
	ns, err := nsClient.Get(o.GetName(), v1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if ns.Status.Phase != v12.NamespaceTerminating {
		err = nsClient.Delete(o.GetName(), nil)
		if apierrors.IsNotFound(err) {
			return nil
		}
	}
	return err
}

func (m *mgr) reconcileResourceToNamespace(obj runtime.Object) (runtime.Object, error) {
	return v3.NamespaceBackedResource.Do(obj, func() (runtime.Object, error) {
		o, err := meta.Accessor(obj)
		if err != nil {
			return obj, condition.Error("MissingMetadata", err)
		}
		t, err := meta.TypeAccessor(obj)
		if err != nil {
			return obj, condition.Error("MissingTypeMetadata", err)
		}

		ns, _ := m.nsLister.Get("", o.GetName())
		if ns == nil {
			nsClient := m.mgmt.K8sClient.CoreV1().Namespaces()
			_, err := nsClient.Create(&v12.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: o.GetName(),
					Annotations: map[string]string{
						"management.cattle.io/system-namespace": "true",
					},
				},
			})
			if err != nil {
				return obj, condition.Error("NamespaceCreationFailure", errors.Wrapf(err, "failed to create namespace for %v %v", t.GetKind(), o.GetName()))
			}
		}

		return obj, nil
	})
}
