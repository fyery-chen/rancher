package harbor

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/registry/common"
	"github.com/sirupsen/logrus"
)

const (
	httpHeaderJSON        = "application/json"
	httpHeaderContentType = "Content-Type"
	httpHeaderAccept      = "Accept"
)

type APIClient struct {
	client *http.Client
	config common.APIClientConfig
}

func NewAPIClient(config common.APIClientConfig) (*APIClient, error) {
	if config.RegistryServer == "" {
		return nil, errors.New("missing registry address")
	}
	url, err := url.Parse(config.RegistryServer)
	if err != nil {
		return nil, errors.Wrapf(err, "parse url %s failed", config.RegistryServer)
	}
	transport := &http.Transport{}
	if url.Scheme == "https" {
		tlsConfig := &tls.Config{}
		if config.RootCA == "" || config.ClientCert == "" || config.ClientKey == "" {
			tlsConfig.InsecureSkipVerify = true
		} else {
			var decodeClientKeyBytes = []byte(config.ClientKey)
			cert, err := tls.X509KeyPair([]byte(config.ClientCert), decodeClientKeyBytes)
			if err != nil {
				return nil, errors.Wrap(err, "load client cert and key failed")
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(config.RootCA))

			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
			tlsConfig.RootCAs = caCertPool
		}

		transport.TLSClientConfig = tlsConfig
	}

	if len(strings.TrimSpace(config.Proxy)) > 0 {
		if proxyURL, err := url.Parse(config.Proxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client := &http.Client{
		Transport: transport,
	}

	return &APIClient{
		client: client,
		config: config,
	}, nil

}

func (ac *APIClient) sendRequest(url string, method string, in []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(string(in)))
	if err != nil {
		return nil, err
	}

	req.Header.Set(httpHeaderAccept, httpHeaderJSON)
	if ac.config.Username != "" && ac.config.Password != "" {
		req.SetBasicAuth(ac.config.Username, ac.config.Password)
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ac *APIClient) Get() ([]byte, error) {
	if strings.TrimSpace(ac.config.RegistryServer) == "" {
		return nil, errors.New("empty url")
	}
	url := ""
	condition := "?page=" + strconv.Itoa(ac.config.Page) + "&page_size=" + strconv.Itoa(ac.config.PageSize)
	switch ac.config.RequestType {
	case common.Project:
		url = ac.config.RegistryServer + "/api/projects" + condition
		return ac.sendRequest(url, http.MethodGet, nil)
	case common.Repository:
		condition = condition + "&project_id=" + strconv.Itoa(ac.config.ProjectID)
		url = ac.config.RegistryServer + "/api/repositories" + condition
		return ac.sendRequest(url, http.MethodGet, nil)
	case common.Tag:
		url = ac.config.RegistryServer + "/api/repositories/" + ac.config.RepositoryName + "/tags"
		return ac.sendRequest(url, http.MethodGet, nil)
	case common.All:
		var projectResp []common.HarborProject
		var repositoryRet []byte
		url = ac.config.RegistryServer + "/api/projects"
		resp, err := ac.sendRequest(url, http.MethodGet, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "send request error, send type: %s", common.All)
		}
		if err = json.Unmarshal(resp, &projectResp); err != nil {
			return nil, errors.Wrap(err, "getting projects json unmarshall error")
		}
		for _, project := range projectResp {
			url = ac.config.RegistryServer + "/api/repositories?project_id=" + strconv.Itoa(project.ProjectID)
			resp, err := ac.sendRequest(url, http.MethodGet, nil)
			if err != nil {
				logrus.Errorf("send request error, send type: %s, error: %s", common.All, err.Error())
				continue
			}
			repositoryRet = append(repositoryRet, resp...)
		}
		return repositoryRet, nil
	}

	return nil, errors.New("cannot to anything")
}

func (ac *APIClient) Post(url string, data []byte) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("Empty url")
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	req.Header.Set(httpHeaderContentType, httpHeaderJSON)
	if ac.config.Username != "" && ac.config.Password != "" {
		req.SetBasicAuth(ac.config.Username, ac.config.Password)
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

func (ac *APIClient) Delete(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("Empty url")
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set(httpHeaderAccept, httpHeaderJSON)
	if ac.config.Username != "" && ac.config.Password != "" {
		req.SetBasicAuth(ac.config.Username, ac.config.Password)
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

func (ac *APIClient) SwitchAccount(username, password string) {
	if len(strings.TrimSpace(username)) == 0 ||
		len(strings.TrimSpace(password)) == 0 {
		return
	}

	ac.config.Username = username
	ac.config.Password = password
}
