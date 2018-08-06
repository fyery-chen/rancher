package cce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fyery-chen/cce-sdk/signer"
	"github.com/fyery-chen/cce-sdk/common"
	"github.com/rancher/kontainer-engine/types"
	"github.com/sirupsen/logrus"
	"crypto/tls"
	"encoding/base64"
	"github.com/rancher/kontainer-engine/drivers/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

type Driver struct {
	types.UnimplementedClusterSizeAccess
	types.UnimplementedVersionAccess

	driverCapabilities types.Capabilities

}

type state struct {
	ClusterName   	string
	DisplayName   	string
	Description     string
	ProjectID     	string
	Zone            string
	ClusterType   	string
	ClusterFlavor 	string
	ClusterVersion 	string
	ClusterBillingMode int64
	ClusterLabels   map[string]string
	ClientID      	string
	ClientSecret  	string
	ContainerNetworkMode string
	ContainerNetworkCidr string
	VpcID         	string
	SubnetID      	string
	HighwaySubnet 	string
	AuthenticatingProxyCa string
	ClusterID       string
	ExternalServerEnabled bool
	ClusterEipId    string
	ClusterJobId    string
	NodeConfig      *common.NodeConfig

	ClusterInfo types.ClusterInfo
}

const (
	retries = 5
	pollInterval = 30
	defaultNamespace = "cattle-system"
)

var isConsumerCloudMember = false

func NewDriver() types.Driver {
	driver := &Driver{
		driverCapabilities: types.Capabilities{
			Capabilities: make(map[int64]bool),
		},
	}

	return driver
}

func (d *Driver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["display-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the name of the cluster that should be displayed to the user",
	}
	driverFlag.Options["project-id"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the ID of your project to use when creating a cluster",
	}
	driverFlag.Options["zone"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The zone to launch the cluster",
		Value: "cn-north-1",
	}
	driverFlag.Options["client-id"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The CCE Client ID to use",
	}
	driverFlag.Options["client-secret"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The CCE Client Secret associated with the Client ID",
	}
	driverFlag.Options["cluster-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The cluster type, VirtualMachine or BareMetal",
		Value: "VirtualMachine",
	}
	driverFlag.Options["cluster-flavor"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The cluster flavor",
		Value: "cce.s2.small",
	}
	driverFlag.Options["cluster-billing-mode"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The bill mode of the cluster",
		Value: "0",
	}
	driverFlag.Options["description"] = &types.Flag{
		Type:  types.StringType,
		Usage: "An optional description of this cluster",
	}
	driverFlag.Options["master-version"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The kubernetes master version",
		Value: "v1.9.2-r1",
	}
	driverFlag.Options["node-count"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The number of nodes to create in this cluster",
		Value: "3",
	}
	driverFlag.Options["vpc-id"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The id of existing vpc",
	}
	driverFlag.Options["subnet-id"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The id of existing subnet",
	}
	driverFlag.Options["highway-subnet"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The id of existing highway subnet when the cluster-type is BareMetal",
	}
	driverFlag.Options["container-network-mode"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The network mode of container",
		Value: "overlay_l2",
	}
	driverFlag.Options["container-network-cidr"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The network cidr of container",
		Value: "172.16.0.0/16",
	}
	driverFlag.Options["cluster-labels"] = &types.Flag{
		Type:  types.StringSliceType,
		Usage: "The map of Kubernetes labels (key/value pairs) to be applied to cluster",
	}
	driverFlag.Options["node-flavor"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The node flavor",
		Value: "s3.large.2",
	}
	driverFlag.Options["available-zone"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The available zone which the nodes in",
		Value: "cn-north-1a",
	}
	driverFlag.Options["ssh-key"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The name of ssh key-pair",
	}
	driverFlag.Options["authenticating-proxy-ca"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The CA for authenticating proxy",
	}
	driverFlag.Options["root-volume-size"] = &types.Flag{
		Type:  types.IntType,
		Usage: "Size of the system disk attached to each node",
		Value: "40",
	}
	driverFlag.Options["root-volume-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Type of the system disk attached to each node",
		Value: "SATA",
	}
	driverFlag.Options["data-volume-size"] = &types.Flag{
		Type:  types.IntType,
		Usage: "Size of the data disk attached to each node",
		Value: "100",
	}
	driverFlag.Options["data-volume-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "Type of the data disk attached to each node",
		Value: "SATA",
	}
	driverFlag.Options["billing-mode"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The bill mode of the node",
		Value: "0",
	}
	driverFlag.Options["node-operation-system"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The operation system of nodes",
		Value: "EulerOS 2.2",
	}
	driverFlag.Options["bms-period-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The period type",
		Value: "month",
	}
	driverFlag.Options["bms-period-num"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The number of period",
		Value: "1",
	}
	driverFlag.Options["bms-is-auto-renew"] = &types.Flag{
		Type:  types.StringType,
		Usage: "If the period is auto renew",
		Value: "false",
	}
	driverFlag.Options["external-server-enabled"] = &types.Flag{
		Type:  types.BoolType,
		Usage: "To enable cluster elastic IP",
	}
	driverFlag.Options["cluster-eip-id"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The id of cluster eip",
	}
	driverFlag.Options["eip-ids"] = &types.Flag{
		Type:  types.StringSliceType,
		Usage: "The list of the exist EIPs",
	}
	driverFlag.Options["eip-count"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The number of eips to be created",
		Value: "3",
	}
	driverFlag.Options["eip-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The type of bandwidth",
		Value: "5-bgp",
	}
	driverFlag.Options["eip-charge-mode"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The charge mode of the bandwidth",
		Value: "traffic",
	}
	driverFlag.Options["eip-bandwidth-size"] = &types.Flag{
		Type:  types.IntType,
		Usage: "The size of bandwidth",
		Value: "10",
	}
	driverFlag.Options["eip-share-type"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The share type of bandwidth",
		Value: "PER",
	}
	driverFlag.Options["node-labels"] = &types.Flag{
		Type:  types.StringSliceType,
		Usage: "The map of Kubernetes labels (key/value pairs) to be applied to each node",
	}

	return &driverFlag, nil
}

func (d *Driver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["description"] = &types.Flag{
		Type:  types.StringType,
		Usage: "An optional description of this cluster",
	}

	return &driverFlag, nil
}

func getStateFromOptions(driverOptions *types.DriverOptions) (state, error) {
	state := state{
		NodeConfig: &common.NodeConfig{
			NodeLabels: map[string]string{},
		},
		ClusterLabels: map[string]string{},
	}
	state.ClusterName = getValueFromDriverOptions(driverOptions, types.StringType, "name").(string)
	state.DisplayName = getValueFromDriverOptions(driverOptions, types.StringType, "display-name", "displayName").(string)
	state.ProjectID = getValueFromDriverOptions(driverOptions, types.StringType, "project-id", "projectId").(string)
	state.Zone = getValueFromDriverOptions(driverOptions, types.StringType, "zone").(string)
	state.Description = getValueFromDriverOptions(driverOptions, types.StringType, "description").(string)
	state.ClusterType = getValueFromDriverOptions(driverOptions, types.StringType, "cluster-type", "clusterType").(string)
	state.ClusterFlavor = getValueFromDriverOptions(driverOptions, types.StringType, "cluster-flavor", "clusterFlavor").(string)
	state.ClusterVersion = getValueFromDriverOptions(driverOptions, types.StringType, "master-version", "masterVersion").(string)
	state.ClientID = getValueFromDriverOptions(driverOptions, types.StringType, "client-id", "accessKey").(string)
	state.ClientSecret = getValueFromDriverOptions(driverOptions, types.StringType, "client-secret", "secretKey").(string)
	state.ClusterBillingMode = getValueFromDriverOptions(driverOptions, types.IntType, "cluster-billing-mode", "clusterBillingMode").(int64)
	state.VpcID = getValueFromDriverOptions(driverOptions, types.StringType, "vpc-id", "vpcId").(string)
	state.SubnetID = getValueFromDriverOptions(driverOptions, types.StringType, "subnet-id", "subnetId").(string)
	state.ContainerNetworkMode = getValueFromDriverOptions(driverOptions, types.StringType, "container-network-mode", "containerNetworkMode").(string)
	state.ContainerNetworkCidr = getValueFromDriverOptions(driverOptions, types.StringType, "container-network-cidr", "containerNetworkCidr").(string)
	state.HighwaySubnet = getValueFromDriverOptions(driverOptions, types.StringType, "highway-subnet", "highwaySubnet").(string)
	state.NodeConfig.NodeFlavor = getValueFromDriverOptions(driverOptions, types.StringType, "node-flavor", "nodeFlavor").(string)
	state.NodeConfig.AvailableZone = getValueFromDriverOptions(driverOptions, types.StringType, "available-zone", "availableZone").(string)
	state.NodeConfig.SSHName = getValueFromDriverOptions(driverOptions, types.StringType, "ssh-key","sshKey").(string)
	state.NodeConfig.RootVolumeSize = getValueFromDriverOptions(driverOptions, types.IntType, "root-volume-size", "rootVolumeSize").(int64)
	state.NodeConfig.RootVolumeType = getValueFromDriverOptions(driverOptions, types.StringType, "root-volume-type", "rootVolumeType").(string)
	state.NodeConfig.DataVolumeSize = getValueFromDriverOptions(driverOptions, types.IntType, "data-volume-size", "dataVolumeSize").(int64)
	state.NodeConfig.DataVolumeType = getValueFromDriverOptions(driverOptions, types.StringType, "data-volume-type", "dataVolumeType").(string)
	state.NodeConfig.BillingMode = getValueFromDriverOptions(driverOptions, types.IntType, "billing-mode", "billingMode").(int64)
	state.NodeConfig.NodeCount = getValueFromDriverOptions(driverOptions, types.IntType, "node-count", "nodeCount").(int64)
	state.NodeConfig.PublicIP.Count = getValueFromDriverOptions(driverOptions, types.IntType, "eip-count", "eipCount").(int64)
	state.NodeConfig.PublicIP.Eip.Iptype = getValueFromDriverOptions(driverOptions, types.StringType, "eip-type", "eipType").(string)
	state.NodeConfig.PublicIP.Eip.Bandwidth.Size = getValueFromDriverOptions(driverOptions, types.IntType, "eip-bandwidth-size", "eipBandwidthSize").(int64)
	state.NodeConfig.PublicIP.Eip.Bandwidth.ShareType = getValueFromDriverOptions(driverOptions, types.StringType, "eip-share-type", "eipShareType").(string)
	state.NodeConfig.PublicIP.Eip.Bandwidth.ChargeMode = getValueFromDriverOptions(driverOptions, types.StringType, "eip-charge-mode", "eipChargeMode").(string)
	state.NodeConfig.NodeOperationSystem = getValueFromDriverOptions(driverOptions, types.StringType, "node-operation-system", "nodeOperationSystem").(string)
	state.NodeConfig.ExtendParam.BMSPeriodType = getValueFromDriverOptions(driverOptions, types.StringType, "bms-period-type", "bmsPeriodType").(string)
	state.NodeConfig.ExtendParam.BMSPeriodNum = getValueFromDriverOptions(driverOptions, types.IntType, "bms-period-num", "bmsPeriodNum").(int64)
	state.NodeConfig.ExtendParam.BMSIsAutoRenew = getValueFromDriverOptions(driverOptions, types.StringType, "bms-is-auto-renew", "bmsIsAutoRenew").(string)
	state.AuthenticatingProxyCa = getValueFromDriverOptions(driverOptions, types.StringType, "authenticating-proxy-ca", "authenticatingProxyCa").(string)
	state.ExternalServerEnabled = getValueFromDriverOptions(driverOptions, types.BoolType, "external-server-enabled", "externalServerEnabled").(bool)
	state.ClusterEipId = getValueFromDriverOptions(driverOptions, types.StringType, "cluster-eip-id", "clusterEipId").(string)

	if state.ClientID == "" || state.ClientSecret == "" {
		kubeconfig, err := rest.InClusterConfig()
		if err != nil {
			return state, err
		}

		clientset, err := kubernetes.NewForConfig(kubeconfig)
		if err != nil {
			return state, fmt.Errorf("error creating clientset: %v", err)
		}
		namespace := os.Getenv("NAMESPACE")
		if namespace == "" {
			namespace = defaultNamespace
		}

		accountSecret, err := clientset.CoreV1().Secrets(namespace).Get("huaweicloud-cce-account", v1.GetOptions{})
		if err != nil {
			return state, fmt.Errorf("error creating service account: %v", err)
		}
		state.ClientID = string(accountSecret.Data["accessKey"])
		state.ClientSecret = string(accountSecret.Data["secretKey"])
		isConsumerCloudMember = true
	}
	eipIds := getValueFromDriverOptions(driverOptions, types.StringSliceType, "eip-ids", "eipIds").(*types.StringSlice)
	logrus.Infof("eipIds: %v", eipIds)
	for _, eipId := range eipIds.Value {
		logrus.Infof("Eip: %s", eipId)
		state.NodeConfig.PublicIP.Ids = append(state.NodeConfig.PublicIP.Ids, eipId)
	}
	nodeLabels := getValueFromDriverOptions(driverOptions, types.StringSliceType, "node-labels", "nodeLabels").(*types.StringSlice)
	for _, nodeLabel := range nodeLabels.Value {
		kv := strings.Split(nodeLabel, "=")
		if len(kv) == 2 {
			state.NodeConfig.NodeLabels[kv[0]] = kv[1]
		}
	}
	clusterLabels := getValueFromDriverOptions(driverOptions, types.StringSliceType,"labels").(*types.StringSlice)
	for _, clusterLabel := range clusterLabels.Value {
		kv := strings.Split(clusterLabel, "=")
		if len(kv) ==  2 {
			state.ClusterLabels[kv[0]] = kv[1]
		}
	}

	return state, state.validate()
}

func getValueFromDriverOptions(driverOptions *types.DriverOptions, optionType string, keys ...string) interface{} {
	switch optionType {
	case types.IntType:
		for _, key := range keys {
			if value, ok := driverOptions.IntOptions[key]; ok {
				return value
			}
		}
		return int64(0)
	case types.StringType:
		for _, key := range keys {
			if value, ok := driverOptions.StringOptions[key]; ok {
				return value
			}
		}
		return ""
	case types.BoolType:
		for _, key := range keys {
			if value, ok := driverOptions.BoolOptions[key]; ok {
				return value
			}
		}
		return false
	case types.StringSliceType:
		for _, key := range keys {
			if value, ok := driverOptions.StringSliceOptions[key]; ok {
				return value
			}
		}
		return &types.StringSlice{}
	}
	return nil
}

func (state *state) validate() error {
	if state.ClusterName == "" {
		return fmt.Errorf("cluster name is required")
	}

	if state.ClientID == "" {
		return fmt.Errorf("client id is required")
	}

	if state.ClientSecret == "" {
		return fmt.Errorf("client secret is required")
	}

	if state.ProjectID == "" {
		return fmt.Errorf("project id is required")
	}

	if state.ClusterType == "" {
		return fmt.Errorf("cluster type is required")
	}

	if state.ClusterFlavor == "" {
		return fmt.Errorf("cluster flavor is required")
	}

	if state.ClusterVersion == "" {
		return fmt.Errorf("cluster version is required")
	}

	return nil
}

func (d *Driver) cceHTTPRequest(state state, uri, method, serviceType string, args interface{}) ([]byte, int, error) {
	var b []byte
	var err error
	var req *http.Request
	statusCode := 0

	//initial signer for request signature
	signer := signer.Signer{
		AccessKey: state.ClientID,
		SecretKey: state.ClientSecret,
		Region:    state.Zone,
		Service:   serviceType,
	}

	if args != nil {
		b, err = json.Marshal(args)
		if err != nil {
			return nil, statusCode, err
		}
	}

	//compose request URL
	requestURL := "https://" + serviceType + "." + state.Zone + "." + common.EndPoint + "/" + uri
	if args != nil {
		req, err = http.NewRequest(method, requestURL, bytes.NewReader(b))
	} else {
		req, err = http.NewRequest(method, requestURL, nil)
	}
	if err != nil {
		return nil, statusCode, err
	}

	req.Header.Set("Content-Type", "application/json")

	//signature http request
	//req.Header.Add("date", time.Now().Format(BasicDateFormat))
	if err := signer.Sign(req); err != nil {
		logrus.Errorf("sign error: %v", err)
		return nil, statusCode, err
	}

	//check request details for debuging
	logrus.Infof("request authorization header after signature: %v", req.Header.Get("authorization"))
	logrus.Infof("request is: %v", req)

	defaultTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	tlsConfig := tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	defaultTransport.TLSClientConfig = &tlsConfig
	httpClient := &http.Client{Transport: defaultTransport}

	//issue request
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("http request error: %v", err)
		return nil, statusCode, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error of response body after I/O reading is: %v", err)
		return nil, statusCode, fmt.Errorf("error reading body: %v", err)
	}

	logrus.Infof("response body after io read: %v", string(body))
	//err code handling
	statusCode = resp.StatusCode
	if statusCode >= 400 && statusCode <= 599 {
		logrus.Errorf("response status code for request is: %d", statusCode)
		return nil, statusCode, fmt.Errorf("Error code: %d", statusCode)
	}
	logrus.Infof("status code is: %d", statusCode)

	return body, statusCode, nil
}

func (d *Driver) addNode(ctx context.Context, state state, num int64)(error) {
	var nodeResp common.NodeInfo
	if isConsumerCloudMember == true && state.ClusterLabels["business"] != ""{
		err := d.preCheck(ctx, state)
		if err != nil {
			return fmt.Errorf("Quota check failed")
		}
	}

	uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID + "/nodes"
	nodesReq := &common.NodeInfo{
		Kind:       "Node",
		ApiVersion: "v3",
		MetaData: common.NodeMetaInfo{
			Name: state.ClusterName,
			Labels: state.NodeConfig.NodeLabels,
		},
		Spec: common.NodeSpecInfo{
			Flavor:        state.NodeConfig.NodeFlavor,
			AvailableZone: state.NodeConfig.AvailableZone,
			Login: common.NodeLogin{
				SSHKey: state.NodeConfig.SSHName,
			},
			RootVolume: common.NodeVolume{
				Size:       state.NodeConfig.RootVolumeSize,
				VolumeType: state.NodeConfig.RootVolumeType,
			},
			DataVolumes: []common.NodeVolume{
				{
					Size:       state.NodeConfig.DataVolumeSize,
					VolumeType: state.NodeConfig.DataVolumeType,
				},
			},
			PublicIP: common.PublicIP{
				Ids:   state.NodeConfig.PublicIP.Ids,
				Count: state.NodeConfig.PublicIP.Count,
				Eip: common.Eip{
					Iptype: state.NodeConfig.PublicIP.Eip.Iptype,
					Bandwidth: common.Bandwidth{
						ChargeMode: state.NodeConfig.PublicIP.Eip.Bandwidth.ChargeMode,
						Size:       state.NodeConfig.PublicIP.Eip.Bandwidth.Size,
						ShareType:  state.NodeConfig.PublicIP.Eip.Bandwidth.ShareType,
					},
				},
			},
			Count:       num,
			BillingMode: state.NodeConfig.BillingMode,
			OperationSystem: state.NodeConfig.NodeOperationSystem,
			ExtendParam: common.ExtendParam{
				BMSPeriodType: state.NodeConfig.ExtendParam.BMSPeriodType,
				BMSPeriodNum: state.NodeConfig.ExtendParam.BMSPeriodNum,
				BMSIsAutoRenew: state.NodeConfig.ExtendParam.BMSIsAutoRenew,
			},
		},
	}

	resp, _, err := d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceCCE, nodesReq)
	if err != nil {
		return fmt.Errorf("error creating node: %v", err)
	}

	if err := json.Unmarshal(resp, &nodeResp); err != nil {
		logrus.Debugf("creating node json unmarshal error is: %v", err)
		return err
	}

	_, err = d.waitForReady(state, "node", "create")
	if err != nil {
		return  err
	}

	logrus.Infof("Starting add node...")

	return nil
}

func (d *Driver) deleteNode(ctx context.Context, state state, num int64)(error) {

	logrus.Infof("Starting delete node...")

	nodes := &common.NodeListInfo{
	}

	uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID + "/nodes"

	resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
	if err != nil {
		return fmt.Errorf("error getting cluster info: %v", err)
	}

	err = json.Unmarshal(resp, nodes)
	if err != nil {
		logrus.Infof("error parsing cluster info: %v", err)
		return fmt.Errorf("error parsing cluster info: %v", err)
	}
	leftNodes := num
	for _, item := range nodes.Items {
		if leftNodes == 0 {
			break
		}

		uri = "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID + "/nodes/" + item.MetaData.Uid
		_, _, err := d.cceHTTPRequest(state, uri, http.MethodDelete, common.ServiceCCE, nil)
		if err != nil {
			return fmt.Errorf("error deleting node: %v", err)
		}
		leftNodes--
	}

	return nil
}

func (d *Driver) deleteResources(state state, source string) error {
	vpcUri := "/v1/" + state.ProjectID + "/vpcs/" + state.VpcID
	subnetUri := "/v1/" + state.ProjectID + "/vpcs/" + state.VpcID + "/subnets/" + state.SubnetID

	if source == "all" {
		_, _, err := d.cceHTTPRequest(state, subnetUri, http.MethodDelete, common.ServiceVPC, nil)
		if err != nil {
			logrus.Errorf("delete resource subnets err, %v", err)
			return err
		}
		_, _  = d.waitForReady(state, "subnet", "delete")

		_, _, err = d.cceHTTPRequest(state, vpcUri, http.MethodDelete, common.ServiceVPC, nil)
		if err != nil {
			logrus.Errorf("delete resource vpcs err, %v", err)
			return err
		}
		_, _  = d.waitForReady(state, "vpc", "delete")
	} else if source == "vpc" {
		_, _, err := d.cceHTTPRequest(state, vpcUri, http.MethodDelete, common.ServiceVPC, nil)
		if err != nil {
			logrus.Errorf("delete resource vpcs err, %v", err)
			return err
		}
		_, _  = d.waitForReady(state, "vpc", "delete")
	} else {
		_, _, err := d.cceHTTPRequest(state, subnetUri, http.MethodDelete, common.ServiceVPC, nil)
		if err != nil {
			logrus.Errorf("delete resource subnets err, %v", err)
			return err
		}
		_, _  = d.waitForReady(state, "subnet", "delete")
	}

	return nil
}

func (d *Driver) preCheck(ctx context.Context, state state) error {
	logrus.Infof("Starting preCheck")

	var b []byte
	e := map[string]interface{}{
		"businessName": state.ClusterLabels["business"],
		"nodeCount": state.NodeConfig.NodeCount,
	}

	b, err := json.Marshal(e)

	requestURL := "https://cattle-cce-service/quotacheckout"
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	defaultTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	tlsConfig := tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	defaultTransport.TLSClientConfig = &tlsConfig
	httpClient := &http.Client{Transport: defaultTransport}

	//issue request
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("http request error: %v", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error of response body after I/O reading is: %v", err)
		return fmt.Errorf("error reading body: %v", err)
	}

	logrus.Infof("response body after io read: %v", string(body))
	//err code handling
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode <= 599 {
		logrus.Errorf("response status code for request is: %d", statusCode)
		return fmt.Errorf("Error code: %d", statusCode)
	}
	logrus.Infof("status code is: %d", statusCode)
	return nil
}

func (d *Driver) Create(ctx context.Context, options *types.DriverOptions, _ *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Infof("Starting create")

	var resp []byte
	var err error
	var vpcResp common.VpcInfo
	var subnetResp common.SubnetInfo
	var clusterResp common.ClusterInfo
	var nodeResp common.NodeInfo
	state, err := getStateFromOptions(options)
	if err != nil {
		return nil, fmt.Errorf("error parsing state: %v", err)
	}
	if isConsumerCloudMember == true && state.ClusterLabels["business"] != "" {
		err = d.preCheck(ctx, state)
		if err != nil {
			return nil, fmt.Errorf("Quota check failed")
		}
	}

	logrus.Infof("Bringing up vpc...")
	if state.VpcID == "" {
		uri := "/v1/" + state.ProjectID + "/vpcs"
		vpcReq := &common.VpcRequest{
			Vpc: common.VpcSt{
				Name: state.ClusterName + "-vpc",
				Cidr: common.DefaultCidr,
			},
		}
		resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceVPC, vpcReq)
		if err != nil {
			return nil, fmt.Errorf("error creating vpc: %v", err)
		}

		if err := json.Unmarshal(resp, &vpcResp); err != nil {
			logrus.Errorf("creating vpc json unmarshal error is: %v", err)
			return nil, err
		}
		state.VpcID = vpcResp.Vpc.Id

		_, err = d.waitForReady(state, "vpc", "create")
		if err != nil {
			d.deleteResources(state, "vpc")
			return nil, err
		}

		subnetReq := &common.SubnetInfo{
			Subnet: common.Subnet{
				Name: state.ClusterName + "-subnet",
				Cidr: common.DefaultCidr,
				GatewayIp: common.DefaultGateway,
				VpcId: vpcResp.Vpc.Id,
				PrimaryDns: "114.114.114.114",
				SecondaryDns: "8.8.8.8",
				DhcpEnable: true,
			},
		}

		uri = "/v1/" + state.ProjectID + "/subnets"
		resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceVPC, subnetReq)
		if err != nil {
			d.deleteResources(state, "vpc")
			return nil, fmt.Errorf("error creating subnet: %v", err)
		}

		if err := json.Unmarshal(resp, &subnetResp); err != nil {
			d.deleteResources(state, "vpc")
			logrus.Errorf("creating subnet json unmarshal error is: %v", err)
			return nil, err
		}
		state.SubnetID = subnetResp.Subnet.Id
		_, err = d.waitForReady(state, "subnet", "create")
		if err != nil {
			d.deleteResources(state, "all")
			return nil, err
		}

		if state.ClusterType == common.BareMetal {
			//TODO: If the cluster type is BareMetal, then create highway Subnet and bare metal host
		}
	}


	logrus.Infof("Creating cce clusters...")
	clusterReq := &common.ClusterInfo{
		Kind: "cluster",
		ApiVersion: "v3",
		MetaData: common.MetaInfo{
			Name: state.ClusterName,
			Labels: state.ClusterLabels,
		},
		Spec: common.SpecInfo{
			ClusterType: state.ClusterType,
			Flavor: state.ClusterFlavor,
			K8sVersion: state.ClusterVersion,
			HostNetwork: &common.NetworkInfo{
				Vpc: state.VpcID,
				Subnet: state.SubnetID,
				HighwaySubnet: state.HighwaySubnet,
			},
			ContainerNetwork: &common.ContainerNetworkInfo{
				Mode: state.ContainerNetworkMode,
				Cidr: state.ContainerNetworkCidr,
			},
			BillingMode: state.ClusterBillingMode,
			Authentication: common.Authentication{
				Mode: "authenticating_proxy",
				AuthenticatingProxy: common.AuthenticatingProxy{
					Ca: state.AuthenticatingProxyCa,
				},
			},
		},
	}

	uri := "/api/v3/projects/" + state.ProjectID + "/clusters"
	resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceCCE, clusterReq)
	if err != nil {
		uri = "/v1/" + state.ProjectID + "/vpcs/" + state.VpcID
		return nil, fmt.Errorf("error creating cluster: %v", err)
	}

	if err := json.Unmarshal(resp, &clusterResp); err != nil {
		logrus.Errorf("creating cluster json unmarshar error is: %v", err)
		return nil, err
	}
	state.ClusterID = clusterResp.MetaData.Uid
	state.ClusterJobId = clusterResp.Status.JobID

	_, err = d.waitForReady(state, common.ServiceCCE, "create")
	if err != nil {
		return nil, err
	}

	logrus.Infof("Cluster provisioned successfully")

	if state.NodeConfig.NodeCount > 0 {
		logrus.Infof("Creating worker nodes")
		uri = "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID + "/nodes"
		nodesReq := &common.NodeInfo{
			Kind:       "Node",
			ApiVersion: "v3",
			MetaData: common.NodeMetaInfo{
				Name: state.ClusterName,
				Labels: state.NodeConfig.NodeLabels,
			},
			Spec: common.NodeSpecInfo{
				Flavor:        state.NodeConfig.NodeFlavor,
				AvailableZone: state.NodeConfig.AvailableZone,
				Login: common.NodeLogin{
					SSHKey: state.NodeConfig.SSHName,
				},
				RootVolume: common.NodeVolume{
					Size:       state.NodeConfig.RootVolumeSize,
					VolumeType: state.NodeConfig.RootVolumeType,
				},
				DataVolumes: []common.NodeVolume{
					{
						Size:       state.NodeConfig.DataVolumeSize,
						VolumeType: state.NodeConfig.DataVolumeType,
					},
				},
				PublicIP: common.PublicIP{
					Ids:   state.NodeConfig.PublicIP.Ids,
					Count: state.NodeConfig.PublicIP.Count,
					Eip: common.Eip{
						Iptype: state.NodeConfig.PublicIP.Eip.Iptype,
						Bandwidth: common.Bandwidth{
							ChargeMode: state.NodeConfig.PublicIP.Eip.Bandwidth.ChargeMode,
							Size:       state.NodeConfig.PublicIP.Eip.Bandwidth.Size,
							ShareType:  state.NodeConfig.PublicIP.Eip.Bandwidth.ShareType,
						},
					},
				},
				Count:       state.NodeConfig.NodeCount,
				BillingMode: state.NodeConfig.BillingMode,
				OperationSystem: state.NodeConfig.NodeOperationSystem,
				ExtendParam: common.ExtendParam{
					BMSPeriodType: state.NodeConfig.ExtendParam.BMSPeriodType,
					BMSPeriodNum: state.NodeConfig.ExtendParam.BMSPeriodNum,
					BMSIsAutoRenew: state.NodeConfig.ExtendParam.BMSIsAutoRenew,
				},
			},
		}

		resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceCCE, nodesReq)
		if err != nil {
			//try again
			resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceCCE, nodesReq)
			if err != nil {
				//delete the cluster
				uri = "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID
				_, _, err = d.cceHTTPRequest(state, uri, http.MethodDelete, common.ServiceCCE, nil)
				if err != nil {
					return nil, fmt.Errorf("error deleting cluster: %v", err)
				}
				//wait to be deleted
				_, err = d.waitForReady(state, "cce", "delete")
				if err != nil {
					return nil, err
				}
			}
			return nil, fmt.Errorf("error creating node: %v", err)
		}

		if err := json.Unmarshal(resp, &nodeResp); err != nil {
			logrus.Errorf("creating node json unmarshal error is: %v", err)
			return nil, err
		}
		_, err = d.waitForReady(state, "node", "create")
		if err != nil {
			return nil, err
		}
	}
	info := &types.ClusterInfo{}
	storeState(info, state)
	return info, nil
}

func (d *Driver) waitForReady(state state, serviceType string, opt string) (interface{}, error) {

	status := ""
	statusCode := 0
	switch serviceType {
	case "cce":
		cluster := &common.ClusterInfo{
			Status: &common.StatusInfo{},
		}
		jobInfo := &common.JobInfo{}
		if opt == "create" {
			for status != "Available" {
				time.Sleep(30 * time.Second)
				logrus.Infof("Waiting for cluster to finish provisioning")
				uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID

				resp, code, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
				if code == 404 {
					uri = "/api/v3/projects/" + state.ProjectID + "/jobs/" + state.ClusterJobId
					resp, _, err = d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
					if err != nil {
						return nil, fmt.Errorf("error getting job info: %v", err)
					}
					err = json.Unmarshal(resp, jobInfo)
					if err != nil {
						logrus.Infof("error parsing job info: %v", err)
						return nil, fmt.Errorf("error parsing job info: %v", err)
					}
					return nil, fmt.Errorf("error create cluster, message: %s", jobInfo.Status.Message)
				}
				if err != nil {
					return nil, fmt.Errorf("error getting cluster info: %v", err)
				}

				err = json.Unmarshal(resp, cluster)
				if err != nil {
					logrus.Infof("error parsing cluster info: %v", err)
					return nil, fmt.Errorf("error parsing cluster info: %v", err)
				}

				status = cluster.Status.Phase
			}
			return cluster, nil
		}
		if opt == "delete" {
			for statusCode != 404 {
				time.Sleep(5 * time.Second)
				logrus.Infof("Waiting for cluster to finish deleting")
				uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID

				_, code, _ := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
				if code == 404 {
					logrus.Info("Delete subnet successfully")
				}
				statusCode = code
			}
			return nil, nil
		}
	case "node":
		nodes := &common.NodeListInfo{
		}
		if opt == "create" {
			for status != "done" {
				time.Sleep(30 * time.Second)
				logrus.Infof("Waiting for nodes to finish provisioning")
				uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID + "/nodes"

				resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
				if err != nil {
					return nil, fmt.Errorf("error getting cluster info: %v", err)
				}

				err = json.Unmarshal(resp, nodes)
				if err != nil {
					logrus.Infof("error parsing cluster info: %v", err)
					return nil, fmt.Errorf("error parsing cluster info: %v", err)
				}
				for _, item := range nodes.Items {
					if item.Status.Phase != "Active" {
						break
					}
					status = "done"
				}
			}
		}

		return nodes, nil
	case "vpc":
		vpc := &common.VpcInfo{}
		if opt == "create" {
			for status != "OK" {
				time.Sleep(5 * time.Second)
				logrus.Infof("Waiting for vpc to finish provisioning")
				uri := "/v1/" + state.ProjectID + "/vpcs/" + state.VpcID

				resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
				if err != nil {
					return nil, fmt.Errorf("error getting cluster info: %v", err)
				}

				err = json.Unmarshal(resp, vpc)
				if err != nil {
					return nil, fmt.Errorf("error parsing vpc info: %v", err)
				}

				status = vpc.Vpc.Status
			}
			return vpc, nil
		} else if opt == "delete" {
			for statusCode != 404 {
				time.Sleep(5 * time.Second)
				logrus.Infof("Waiting for vpc to finish deleting")
				uri := "/v1/" + state.ProjectID + "/vpcs/" + state.VpcID

				_, statusCode, _ := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
				if statusCode == 404 {
					logrus.Info("Delete vpc successfully")
				}
			}
			return nil, nil
		}
	case "subnet":
		subnet := &common.SubnetInfo{}
		if opt == "create" {
			for status != "ACTIVE" {
				time.Sleep(5 * time.Second)
				logrus.Infof("Waiting for subnet to finish provisioning")
				uri := "/v1/" + state.ProjectID + "/subnets/" + state.SubnetID

				resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
				if err != nil {
					return nil, fmt.Errorf("error getting cluster info: %v", err)
				}

				err = json.Unmarshal(resp, subnet)
				if err != nil {
					return nil, fmt.Errorf("error parsing vpc info: %v", err)
				}

				status = subnet.Subnet.Status
			}
			return subnet, nil
		} else if opt == "delete" {
			for statusCode != 404 {
				time.Sleep(5 * time.Second)
				logrus.Infof("Waiting for subnet to finish deleting")
				uri := "/v1/" + state.ProjectID + "/subnets/" + state.SubnetID

				_, statusCode, _ := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
				if statusCode == 404 {
					logrus.Info("Delete subnet successfully")
				}
			}
			return nil, nil
		}
	case "eip":
		eip := &common.EipResp{}
		if opt == "create" {
			for status != "ACTIVE" {
				time.Sleep(10 * time.Second)
				logrus.Infof("Waiting for eip to finish associating")
				uri := "/v1/" + state.ProjectID + "/publicips/" + state.ClusterEipId

				resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
				if err != nil {
					return nil, fmt.Errorf("error getting eip info: %v", err)
				}

				err = json.Unmarshal(resp, eip)
				if err != nil {
					logrus.Infof("error parsing eip info: %v", err)
					return nil, fmt.Errorf("error parsing eip info: %v", err)
				}

				status = eip.Eip.Status
			}
		}

		return eip, nil
	}

	return nil, fmt.Errorf("have no this service type: %s", serviceType)
}

func storeState(info *types.ClusterInfo, state state) error {
	data, err := json.Marshal(state)

	if err != nil {
		return err
	}

	if info.Metadata == nil {
		info.Metadata = map[string]string{}
	}

	info.Metadata["state"] = string(data)

	return nil
}

func getState(info *types.ClusterInfo) (state, error) {
	state := state{}

	err := json.Unmarshal([]byte(info.Metadata["state"]), &state)
	if err != nil {
		logrus.Errorf("Error encountered while marshalling state: %v", err)
	}

	return state, err
}

func (d *Driver) Update(ctx context.Context, info *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {
	logrus.Infof("Starting update")

	state, err := getState(info)
	if err != nil {
		return nil, err
	}

	newState, err := getStateFromOptions(opts)
	if err != nil {
		return nil, err
	}
	newState.ClusterID = state.ClusterID

	if newState.NodeConfig.NodeCount > state.NodeConfig.NodeCount {
		addNodeCount := newState.NodeConfig.NodeCount - state.NodeConfig.NodeCount
		err = d.addNode(ctx, newState, addNodeCount)
		if err != nil {
			return nil, fmt.Errorf("error adding cluster: %v", err)
		}
	} else if newState.NodeConfig.NodeCount < state.NodeConfig.NodeCount {
		deleteNodeCount := state.NodeConfig.NodeCount - newState.NodeConfig.NodeCount
		err = d.deleteNode(ctx, newState, deleteNodeCount)
		if err != nil {
			return nil, fmt.Errorf("error deleting cluster: %v", err)
		}
	}

	state.NodeConfig.NodeCount = newState.NodeConfig.NodeCount
	logrus.Infof("Update complete")
	return info, storeState(info, state)
}

func getInternalIp(ep string) (string){
	if ep == "" {
		return ""
	}
	kv := strings.Split(ep, ":")
	if len(kv) > 1 {
		ip := strings.Split(kv[1], "//")
		if len(ip) ==2 {
			return ip[1]
		}
	}
	return ""
}

func (d *Driver) PostCheck(ctx context.Context, info *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Infof("Starting post-check")

	state, err := getState(info)
	if err != nil {
		return nil, err
	}

	uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID

	resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
	if err != nil {
		return nil, fmt.Errorf("error posting check: %v", err)
	}

	cluster := common.ClusterInfo{}
	err = json.Unmarshal(resp, &cluster)
	if err != nil {
		return nil, fmt.Errorf("error parsing cluster: %v", err)
	}

	if state.ExternalServerEnabled {
		if state.ClusterEipId == "" {
			logrus.Infof("Creating EIP ...")
			eipReq := &common.EipAllocArg{
				EipDesc: common.PubIp{
					Type: "5_bgp",
				},
				BandWidth: common.BandwidthDesc{
					Name:    state.ClusterName,
					Size:    10,
					ShrType: "PER",
					ChgMode: "traffic",
				},
			}
			uri = "/v1/" + state.ProjectID + "/publicips"
			resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPost, common.ServiceVPC, eipReq)
			if err != nil {
				return nil, fmt.Errorf("error creating eip")
			}

			eip := common.EipResp{}
			err = json.Unmarshal(resp, &eip)
			if err != nil {
				return nil, fmt.Errorf("error parsing public ip: %v", err)
			}
			state.ClusterEipId = eip.Eip.Id
		}
		uri := "/v1/" + state.ProjectID + "/subnets/" + state.SubnetID

		resp, _, err := d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
		if err != nil {
			return nil, fmt.Errorf("error getting cluster info: %v", err)
		}
		subnet := &common.SubnetInfo{}
		err = json.Unmarshal(resp, subnet)
		if err != nil {
			return nil, fmt.Errorf("error parsing vpc info: %v", err)
		}

		networkId := subnet.Subnet.NetworkId
		logrus.Infof("Get port info")
		uri = "/v1/ports?network_id=" + networkId
		resp, _, err = d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
		if err != nil {
			return nil, fmt.Errorf("error requring port info")
		}
		ports := common.Ports{}
		err = json.Unmarshal(resp, &ports)
		if err != nil {
			return nil, fmt.Errorf("error parsing port info")
		}
		internalEndpoint := ""
		for _, ep := range cluster.Status.Endpoints{
			if ep.Type == "Internal" {
				internalEndpoint = ep.Url
			}
		}
		internalIp := getInternalIp(internalEndpoint)
		portId := ""
		for _, port := range ports.Ports {
			for _, fixedIp := range port.FixedIps {
				if fixedIp.IpAddress == internalIp {
					portId = port.Id
					break
				}
			}
			if portId != "" {
				break
			}
		}
		logrus.Infof("Start Associating EIP, %s, %s, %s", internalEndpoint, internalIp, portId)
		uri = "/v1/" + state.ProjectID + "/publicips/" + state.ClusterEipId
		assocArg := common.EipAssocArg{
			Port: common.PortDesc{
				PortId: portId,
			},
		}
		resp, _, err = d.cceHTTPRequest(state, uri, http.MethodPut, common.ServiceVPC, &assocArg)
		if err != nil {
			return nil, fmt.Errorf("error associating eip")
		}
		eip := common.EipResp{}
		err = json.Unmarshal(resp, &eip)
		if err != nil {
			return nil, fmt.Errorf("error parsing public ip: %v", err)
		}
		_, err = d.waitForReady(state, "eip", "create")
		if err != nil {
			return nil, err
		}
	}

	logrus.Infof("Requiring cluster cert information")
	uri = "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID + "/clustercert"
	resp, _, err = d.cceHTTPRequest(state, uri, http.MethodGet, common.ServiceCCE, nil)
	if err != nil {
		return nil, fmt.Errorf("error requiring cluster cert: %v", err)
	}

	clusterCert := common.ClusterCert{}
	err = json.Unmarshal(resp, &clusterCert)
	if err != nil {
		return nil, fmt.Errorf("error parsing cluster cert: %v", err)
	}

	info.Version = state.ClusterVersion
	info.NodeCount = state.NodeConfig.NodeCount
	if len(clusterCert.Clusters) == 1 {
		info.Endpoint = clusterCert.Clusters[0].Cluster.Server
	} else {
		for _, cluster := range clusterCert.Clusters{
			if cluster.Name == "externalCluster" {
				info.Endpoint = cluster.Cluster.Server
				break
			}
		}
	}

	info.Status = cluster.Status.Phase
	info.ClientKey = clusterCert.Users[0].User.ClientKeyData
	info.ClientCertificate = clusterCert.Users[0].User.ClientCertificateData
	info.Username = clusterCert.Users[0].Name
	info.RootCaCertificate = clusterCert.Clusters[0].Cluster.CertificateAuthorityData

	capem, err := base64.StdEncoding.DecodeString(info.RootCaCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CA: %v", err)
	}

	key, err := base64.StdEncoding.DecodeString(info.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode client key: %v", err)
	}

	cert, err := base64.StdEncoding.DecodeString(info.ClientCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode client cert: %v", err)
	}

	host := info.Endpoint
	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}

	config := &rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   capem,
			KeyData:  key,
			CertData: cert,
		},
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating clientset: %v", err)
	}

	failureCount := 0

	for {
		info.ServiceAccountToken, err = util.GenerateServiceAccountToken(clientset)

		if err == nil {
			logrus.Info("service account token generated successfully")
			break
		} else {
			if failureCount < retries {
				logrus.Infof("service account token generation failed, retries left: %v", retries-failureCount)
				failureCount = failureCount + 1

				time.Sleep(pollInterval * time.Second)
			} else {
				logrus.Error("retries exceeded, failing post-check")
				return nil, err
			}
		}
	}

	logrus.Info("post-check completed successfully,info: %v", *info)

	return info, nil
}

func (d *Driver) Remove(ctx context.Context, info *types.ClusterInfo) error {
	logrus.Infof("Starting delete cluster")

	state, err := getState(info)
	if err != nil {
		return fmt.Errorf("error getting state: %v", err)
	}

	uri := "/api/v3/projects/" + state.ProjectID + "/clusters/" + state.ClusterID
	_, _, err = d.cceHTTPRequest(state, uri, http.MethodDelete, common.ServiceCCE,nil)
	if err != nil {
		return fmt.Errorf("error deleting cluster: %v", err)
	}

	return nil
}

func (d *Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return &d.driverCapabilities, nil
}

