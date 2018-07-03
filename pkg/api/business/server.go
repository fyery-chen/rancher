package business

import (
	"context"
	"net/http"

	normanapi "github.com/rancher/norman/api"
	"github.com/rancher/norman/api/builtin"
	"github.com/rancher/norman/pkg/subscribe"
	"github.com/rancher/rancher/pkg/api/controllers/dynamicschema"
	"github.com/rancher/rancher/pkg/api/controllers/settings"
	"github.com/rancher/rancher/pkg/api/controllers/whitelistproxy"

	"github.com/rancher/rancher/pkg/api/business/businessstored"
	"github.com/rancher/rancher/pkg/clustermanager"
	businessSchema "github.com/rancher/types/apis/cloud.huawei.com/v1/schema"
	"github.com/rancher/types/config"
)

func New(ctx context.Context, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager,
	k8sProxy http.Handler) (http.Handler, error) {
	subscribe.Register(&builtin.Version, scaledContext.Schemas)
	subscribe.Register(&businessSchema.Version, scaledContext.Schemas)

	if err := businessstored.Setup(ctx, scaledContext, clusterManager, k8sProxy); err != nil {
		return nil, err
	}


	server := normanapi.NewAPIServer()
	server.AccessControl = scaledContext.AccessControl

	if err := server.AddSchemas(scaledContext.Schemas); err != nil {
		return nil, err
	}

	dynamicschema.Register(scaledContext, server.Schemas)
	whitelistproxy.Register(scaledContext)
	err := settings.Register(scaledContext)

	return server, err
}
