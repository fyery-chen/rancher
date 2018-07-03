package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/auth/providers/publicapi"
	authrequests "github.com/rancher/rancher/pkg/auth/requests"
	"github.com/rancher/rancher/pkg/auth/tokens"
	"github.com/rancher/rancher/pkg/clustermanager"
	"github.com/rancher/rancher/pkg/controllers/user/pipeline/hooks"
	rancherdialer "github.com/rancher/rancher/pkg/dialer"
	"github.com/rancher/rancher/pkg/dynamiclistener"
	"github.com/rancher/rancher/pkg/httpproxy"
	"github.com/rancher/rancher/pkg/rkenodeconfigserver"
	"github.com/rancher/rancher/server/ui"
	"github.com/rancher/rancher/server/whitelist"
	//managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	businessSchema "github.com/rancher/types/apis/cloud.huawei.com/v1/schema"
	businessAPI "github.com/rancher/rancher/pkg/api/business"
	"github.com/rancher/types/config"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/kubernetes/cmd/kube-apiserver/app"
)

func Start(ctx context.Context, httpPort, httpsPort int, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager) error {
	tokenAPI, err := tokens.NewAPIHandler(ctx, scaledContext)
	if err != nil {
		return err
	}

	publicAPI, err := publicapi.NewHandler(ctx, scaledContext)
	if err != nil {
		return err
	}
/*
	k8sProxy := k8sProxyPkg.New(scaledContext, scaledContext.Dialer)


	managementAPI, err := managementapi.New(ctx, scaledContext, clusterManager, k8sProxy)
	if err != nil {
		return err
	}
*/

	businessAPI, err := businessAPI.New(ctx, scaledContext, clusterManager, nil)
	if err != nil {
		return err
	}
	root := mux.NewRouter()
	root.UseEncodedPath()

	app.DefaultProxyDialer = utilnet.DialFunc(scaledContext.Dialer.LocalClusterDialer())

	rawAuthedAPIs := newAuthed(tokenAPI, nil, businessAPI, nil)

	authedHandler, err := authrequests.NewAuthenticationFilter(ctx, scaledContext, rawAuthedAPIs)
	if err != nil {
		return err
	}

	webhookHandler := hooks.New(scaledContext)

	root.PathPrefix("/v3-public").Handler(publicAPI)
	//root.PathPrefix("/v3").Handler(authedHandler)
	root.PathPrefix("/v1").Handler(authedHandler)
	root.PathPrefix("/hooks").Handler(webhookHandler)
	root.NotFoundHandler = ui.UI(http.NotFoundHandler())

	registerHealth(root)

	dynamiclistener.Start(ctx, scaledContext, httpPort, httpsPort, root)
	return nil
}

func newAuthed(tokenAPI http.Handler, managementAPI http.Handler, businessAPI http.Handler, k8sproxy http.Handler) *mux.Router {
	authed := mux.NewRouter()
	authed.UseEncodedPath()
	//authed.PathPrefix(managementSchema.Version.Path).Handler(managementAPI)
	authed.PathPrefix(businessSchema.Version.Path).Handler(businessAPI)

	return authed
}

func connectHandlers(scaledContext *config.ScaledContext) (http.Handler, http.Handler) {
	if f, ok := scaledContext.Dialer.(*rancherdialer.Factory); ok {
		return f.TunnelServer, rkenodeconfigserver.Handler(f.TunnelAuthorizer, scaledContext)
	}

	return http.NotFoundHandler(), http.NotFoundHandler()
}

func newProxy() http.Handler {
	return httpproxy.NewProxy("/proxy/", whitelist.Proxy.Get)
}
