package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/parse"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rancher/pkg/clustermanager"
	monitorutil "github.com/rancher/rancher/pkg/monitoring"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	mgmtclientv3 "github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config/dialer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewIstioGraphHandler(dialerFactory dialer.Factory, clustermanager *clustermanager.Manager) *IstioGraphHandler {
	return &IstioGraphHandler{
		dialerFactory:  dialerFactory,
		clustermanager: clustermanager,
	}
}

type IstioGraphHandler struct {
	dialerFactory  dialer.Factory
	clustermanager *clustermanager.Manager
}

func (h *IstioGraphHandler) QuerySeriesAction(actionName string, action *types.Action, apiContext *types.APIContext) error {
	var queryGraphInput v3.QueryGraphInput
	actionInput, err := parse.ReadBody(apiContext.Request)
	if err != nil {
		return err
	}

	if err = convert.ToObj(actionInput, &queryGraphInput); err != nil {
		return err
	}

	inputParser := newClusterGraphInputParser(queryGraphInput)
	if err = inputParser.parse(); err != nil {
		return err
	}

	clusterName := inputParser.ClusterName
	userContext, err := h.clustermanager.UserContext(clusterName)
	if err != nil {
		return fmt.Errorf("get usercontext failed, %v", err)
	}

	prometheusName, prometheusNamespace := monitorutil.IstioMonitoringInfo()
	token, err := getAuthToken(userContext, prometheusName, prometheusNamespace)
	if err != nil {
		return err
	}

	reqContext, cancel := context.WithTimeout(context.Background(), prometheusReqTimeout)
	defer cancel()

	svcName, svcNamespace, svcPort := monitorutil.IstioPrometheusEndpoint()
	prometheusQuery, err := NewPrometheusQuery(reqContext, clusterName, token, svcNamespace, svcName, svcPort, h.dialerFactory)
	if err != nil {
		return err
	}

	mgmtClient := h.clustermanager.ScaledContext.Management
	nodeLister := mgmtClient.Nodes(metav1.NamespaceAll).Controller().Lister()

	nodeMap, err := getNodeName2InternalIPMap(nodeLister, clusterName)
	if err != nil {
		return err
	}

	var graphs []mgmtclientv3.ClusterMonitorGraph
	if err = access.List(apiContext, apiContext.Version, mgmtclientv3.IstioMonitorGraphType, &types.QueryOptions{Conditions: inputParser.Conditions}, &graphs); err != nil {
		return err
	}

	var queries []*PrometheusQuery
	for _, graph := range graphs {
		g := graph
		graphName := getRefferenceGraphName(g.ClusterID, g.Name)
		monitorMetrics, err := graph2Metrics(userContext, mgmtClient, clusterName, g.ResourceType, graphName, g.MetricsSelector, g.DetailsMetricsSelector, inputParser.Input.MetricParams, inputParser.Input.IsDetails)
		if err != nil {
			return err
		}

		queries = append(queries, metrics2PrometheusQuery(monitorMetrics, inputParser.Start, inputParser.End, inputParser.Step, isInstanceGraph(g.GraphType))...)
	}

	seriesSlice, err := prometheusQuery.Do(queries)
	if err != nil {
		return fmt.Errorf("query series failed, %v", err)
	}

	if seriesSlice == nil {
		apiContext.WriteResponse(http.StatusNoContent, nil)
		return nil
	}

	collection := v3.QueryClusterGraphOutput{Type: "collection"}
	for k, v := range seriesSlice {
		graphName, resourceType, _ := parseID(k)
		series := convertInstance(v, nodeMap, resourceType)
		queryGraph := v3.QueryClusterGraph{
			GraphName: graphName,
			Series:    series,
		}
		collection.Data = append(collection.Data, queryGraph)
	}

	res, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("marshal query series result failed, %v", err)
	}

	apiContext.Response.Write(res)
	return nil
}
