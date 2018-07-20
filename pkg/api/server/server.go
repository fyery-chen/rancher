package server

import (
	"context"
	"net/http"

	normanapi "github.com/rancher/norman/api"
	"github.com/rancher/norman/api/builtin"
	"github.com/rancher/norman/pkg/subscribe"
	"github.com/rancher/rancher/pkg/api/controllers/dynamicschema"
	"github.com/rancher/rancher/pkg/api/controllers/settings"
	"github.com/rancher/rancher/pkg/api/server/managementstored"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	businessSchema "github.com/rancher/types/apis/cloud.huawei.com/v3/schema"
	"github.com/rancher/types/config"
)

func New(ctx context.Context, scaledContext *config.ScaledContext) (http.Handler, error) {
	subscribe.Register(&builtin.Version, scaledContext.Schemas)
	subscribe.Register(&managementSchema.Version, scaledContext.Schemas)
	subscribe.Register(&businessSchema.Version, scaledContext.Schemas)

	if err := managementstored.Setup(ctx, scaledContext); err != nil {
		return nil, err
	}

	server := normanapi.NewAPIServer()
	server.AccessControl = scaledContext.AccessControl

	if err := server.AddSchemas(scaledContext.Schemas); err != nil {
		return nil, err
	}

	dynamicschema.Register(scaledContext, server.Schemas)
	err := settings.Register(scaledContext)

	return server, err
}
