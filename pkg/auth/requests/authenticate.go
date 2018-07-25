package requests

import (
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	"fmt"
	"strings"

	"github.com/rancher/rancher/pkg/auth/tokens"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	businessv3 "github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"github.com/rancher/norman/httperror"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Authenticator interface {
	Authenticate(req *http.Request) (authed bool, user string, groups []string, err error)
	TokenFromRequest(req *http.Request) (*v3.Token, error)
	Checkout(req *http.Request) (error)
}

func NewAuthenticator(ctx context.Context, mgmtCtx *config.ScaledContext) Authenticator {
	tokenInformer := mgmtCtx.Management.Tokens("").Controller().Informer()
	tokenInformer.AddIndexers(map[string]cache.IndexFunc{tokenKeyIndex: tokenKeyIndexer})

	return &tokenAuthenticator{
		ctx:          ctx,
		tokenIndexer: tokenInformer.GetIndexer(),
		tokenClient:  mgmtCtx.Management.Tokens(""),
		clusterClient: mgmtCtx.Management.Clusters(""),
		nodeClient: mgmtCtx.Management.Nodes(""),
		businessClient: mgmtCtx.Business.Businesses(""),
	}
}

type tokenAuthenticator struct {
	ctx          context.Context
	tokenIndexer cache.Indexer
	tokenClient  v3.TokenInterface
	clusterClient v3.ClusterInterface
	businessClient businessv3.BusinessInterface
	nodeClient v3.NodeInterface
}

const (
	tokenKeyIndex = "authn.management.cattle.io/token-key-index"
)

func tokenKeyIndexer(obj interface{}) ([]string, error) {
	token, ok := obj.(*v3.Token)
	if !ok {
		return []string{}, nil
	}

	return []string{token.Token}, nil
}

func (a *tokenAuthenticator) Authenticate(req *http.Request) (bool, string, []string, error) {
	token, err := a.TokenFromRequest(req)
	if err != nil {
		return false, "", []string{}, err
	}

	var groups []string
	for _, principal := range token.GroupPrincipals {
		// TODO This is a short cut for now. Will actually need to lookup groups in future
		name := strings.TrimPrefix(principal.Name, "local://")
		groups = append(groups, name)
	}

	return true, token.UserID, groups, nil
}

func (a *tokenAuthenticator) TokenFromRequest(req *http.Request) (*v3.Token, error) {
	tokenAuthValue := tokens.GetTokenAuthFromRequest(req)
	if tokenAuthValue == "" {
		return nil, fmt.Errorf("must authenticate")
	}

	tokenName, tokenKey := tokens.SplitTokenParts(tokenAuthValue)
	if tokenName == "" || tokenKey == "" {
		return nil, fmt.Errorf("must authenticate")
	}

	lookupUsingClient := false
	objs, err := a.tokenIndexer.ByIndex(tokenKeyIndex, tokenKey)
	if err != nil {
		if apierrors.IsNotFound(err) {
			lookupUsingClient = true
		} else {
			return nil, fmt.Errorf("failed to retrieve auth token from cache, error: %v", err)
		}
	} else if len(objs) == 0 {
		lookupUsingClient = true
	}

	storedToken := &v3.Token{}
	if lookupUsingClient {
		storedToken, err = a.tokenClient.Get(tokenName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return nil, fmt.Errorf("must authenticate")
			}
			return nil, fmt.Errorf("failed to retrieve auth token, error: %#v", err)
		}
	} else {
		storedToken = objs[0].(*v3.Token)
	}

	if storedToken.Token != tokenKey || storedToken.ObjectMeta.Name != tokenName {
		return nil, fmt.Errorf("must authenticate")
	}

	if tokens.IsExpired(*storedToken) {
		return nil, fmt.Errorf("must authenticate")
	}

	return storedToken, nil
}

func (a *tokenAuthenticator)Checkout(req *http.Request) (error) {
	logrus.Debugf("Checkout the business quota")

	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("checkout failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	input := businessv3.BusinessQuotaCheck{}

	err = json.Unmarshal(bytes, &input)
	if err != nil {
		logrus.Errorf("unmarshal failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	//check business quota if it is cce provider
	businessName := input.BusinessName
	set := labels.Set{}
	set["businessName"] = businessName
	business, err := a.businessClient.Get(businessName, v1.GetOptions{})
	if err != nil {
		return err
	}
	logrus.Infof("Retrive businesses: %v", business)
	clusters, err := a.clusterClient.List(v1.ListOptions{LabelSelector: set.String()})
	if err != nil {
		return err
	}
	field := fields.Set{}
	requestedHosts := 0
	for _, cluster := range clusters.Items {
		field["namespace"] = cluster.Name
		nodes, err := a.nodeClient.List(v1.ListOptions{FieldSelector: field.String()})
		if err != nil {
			return err
		}
		requestedHosts += len(nodes.Items)
	}

	logrus.Infof("Business name: %s input name: %s", business.Name, input.BusinessName)
	requestedHosts += input.NodeCount
	if requestedHosts > business.Spec.NodeCount {
		return fmt.Errorf("there is no enough quota")
	}

	return nil
}
