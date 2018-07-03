package businessstored

import (
	"context"
	"net/http"
	"sync"

	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/types"
	ccluster "github.com/rancher/rancher/pkg/api/customization/cluster"
	"github.com/rancher/rancher/pkg/api/customization/business"
	"github.com/rancher/rancher/pkg/api/store/scoped"
	"github.com/rancher/rancher/pkg/clustermanager"

	businessschema "github.com/rancher/types/apis/cloud.huawei.com/v1/schema"
	"github.com/rancher/types/client/cloud/v1"
	businessclient  "github.com/rancher/types/client/cloud/v1"
	"github.com/rancher/types/config"
)

func Setup(ctx context.Context, apiContext *config.ScaledContext, clusterManager *clustermanager.Manager,
	k8sProxy http.Handler) error {
	// Here we setup all types that will be stored in the business cluster
	schemas := apiContext.Schemas

	wg := &sync.WaitGroup{}
	factory := &crd.Factory{ClientGetter: apiContext.ClientGetter}

	createCrd(ctx, wg, factory, schemas, &businessschema.Version,
		businessclient.BusinessQuotaType, businessclient.BusinessQuotaType)

	wg.Wait()

	businessQuota(schemas, apiContext, clusterManager, k8sProxy)
	setupScopedTypes(schemas)

	return nil
}

func setupScopedTypes(schemas *types.Schemas) {
	for _, schema := range schemas.Schemas() {
		if schema.Scope != types.NamespaceScope || schema.Store == nil || schema.Store.Context() != config.ManagementStorageContext {
			continue
		}

		for _, key := range []string{"projectId", "clusterId"} {
			ns, ok := schema.ResourceFields["namespaceId"]
			if !ok {
				continue
			}

			if _, ok := schema.ResourceFields[key]; !ok {
				continue
			}

			schema.Store = scoped.NewScopedStore(key, schema.Store)
			ns.Required = false
			schema.ResourceFields["namespaceId"] = ns
			break
		}
	}
}

func businessQuota(schemas *types.Schemas, managementContext *config.ScaledContext, clusterManager *clustermanager.Manager, k8sProxy http.Handler) {
	handler := business.NewbusinessHandler(managementContext)

	schema := schemas.Schema(&businessschema.Version, client.BusinessQuotaType)
	schema.Formatter = ccluster.Formatter
	schema.ActionHandler = handler.Checkout
}