package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type NetlifyClient struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

type NetlifyTransport struct {
	T     http.RoundTripper
	Token string
}

type Site struct {
	Id           string `json:"id"`
	CustomDomain string `json:"custom_domain"`
	Name         string `json:"name"`
	Url          string `json:"url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	State        string `json:"state"`
}

type SiteRequest struct {
	Name         string     `json:"name"`
	CustomDomain string     `json:"custom_domain"`
	Repo         Repository `json:"repo"`
}

type Repository struct {
	Provider    string `json:"provider"`
	Path        string `json:"repo"`
	Branch      string `json:"branch"`
	DeployKeyId string `json:"deploy_key_id"`
	Cmd         string `json:"cmd"`
	Dir         string `json:"dir"`
	Url         string `json:"repo_url"`
}

type DeployKey struct {
	Id        string `json:"id"`
	Key       string `json:"public_key"`
	CreatedAt string `json:"created_at"`
}

func (n NetlifyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+n.Token)
	req.Header.Set("Content-Type", "application/json")
	return n.T.RoundTrip(req)
}

func (c *NetlifyClient) Do(method string, path string, body *bytes.Buffer, dest any) error {
	fullURL := c.BaseURL.String() + path

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if dest != nil {
		err = json.Unmarshal(resBody, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

// Sites related function

func (c *NetlifyClient) ListSites() (*[]Site, error) {
	var sites []Site

	err := c.Do(http.MethodGet, "sites/", &bytes.Buffer{}, &sites)
	if err != nil {
		return nil, err
	}

	return &sites, nil
}

func (c *NetlifyClient) CreateSite(req SiteRequest) (*Site, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resSite Site
	err = c.Do(http.MethodPost, "sites/", bytes.NewBuffer(jsonValue), &resSite)
	if err != nil {
		return nil, err
	}

	return &resSite, nil
}

func (c *NetlifyClient) GetSite(siteId string) (*Site, error) {
	var site Site

	err := c.Do(http.MethodGet, "sites/"+siteId, &bytes.Buffer{}, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

func (c *NetlifyClient) UpdateSite(siteId string, req SiteRequest) (*Site, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resSite Site
	err = c.Do(http.MethodPatch, "sites/"+siteId, bytes.NewBuffer(jsonValue), &resSite)
	if err != nil {
		return nil, err
	}

	return &resSite, nil
}

func (c *NetlifyClient) DeleteSite(siteId string) error {
	return c.Do(http.MethodDelete, "sites/"+siteId, &bytes.Buffer{}, nil)
}

// Deploy key related functions

func (c *NetlifyClient) CreateDeployKey() (*DeployKey, error) {
	var key DeployKey
	err := c.Do(http.MethodPost, "deploy_keys/", &bytes.Buffer{}, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *NetlifyClient) GetDeployKey(keyId string) (*DeployKey, error) {
	var key DeployKey
	err := c.Do(http.MethodGet, "deploy_keys/"+keyId, &bytes.Buffer{}, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *NetlifyClient) DeleteDeployKey(keyId string) error {
	return c.Do(http.MethodDelete, "deploy_keys/"+keyId, &bytes.Buffer{}, nil)
}

// Create NetlifyClient

func NewNetlifyClient(baseUrl string, personalToken string) (*NetlifyClient, error) {
	parsedURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	tr := &NetlifyTransport{Token: personalToken, T: http.DefaultTransport}
	client := &http.Client{
		Transport: tr,
	}

	return &NetlifyClient{
		BaseURL:    parsedURL,
		HTTPClient: client,
	}, nil
}
