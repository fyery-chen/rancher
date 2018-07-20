package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	managementapi "github.com/rancher/rancher/pkg/api/server"
	authrequests "github.com/rancher/rancher/pkg/auth/requests"
	"github.com/rancher/rancher/pkg/auth/tokens"
	//"github.com/rancher/rancher/server/ui"
	businessSchema "github.com/rancher/types/apis/cloud.huawei.com/v3/schema"
	"github.com/rancher/types/config"
	"github.com/rancher/rancher/pkg/dynamiclistener"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
)

func Start(ctx context.Context, httpPort, httpsPort int, scaledContext *config.ScaledContext) error {
	tokenAPI, err := tokens.NewAPIHandler(ctx, scaledContext)
	if err != nil {
		return err
	}

	managementAPI, err := managementapi.New(ctx, scaledContext)
	if err != nil {
		return err
	}

	root := mux.NewRouter()
	root.UseEncodedPath()

	//app.DefaultProxyDialer = utilnet.DialFunc(scaledContext.Dialer.LocalClusterDialer())

	rawAuthedAPIs := newAuthed(tokenAPI, managementAPI, nil)

	authedHandler, err := authrequests.NewAuthenticationFilter(ctx, scaledContext, rawAuthedAPIs)
	if err != nil {
		return err
	}

	root.PathPrefix("/v1").Handler(authedHandler)
	root.PathPrefix("/v3").Handler(authedHandler)
	root.PathPrefix("/k8s/clusters/").Handler(authedHandler)
	root.PathPrefix("/meta").Handler(authedHandler)
	//root.NotFoundHandler = ui.UI(http.NotFoundHandler())

	registerHealth(root)

	dynamiclistener.Start(ctx, scaledContext, httpPort, httpsPort, root)
	return nil
}

func newAuthed(tokenAPI http.Handler, managementAPI http.Handler, k8sproxy http.Handler) *mux.Router {
	authed := mux.NewRouter()
	authed.UseEncodedPath()
	authed.PathPrefix("/v3/token").Handler(tokenAPI)
	authed.PathPrefix(businessSchema.Version.Path).Handler(managementAPI)
	authed.PathPrefix(managementSchema.Version.Path).Handler(managementAPI)

	return authed
}
