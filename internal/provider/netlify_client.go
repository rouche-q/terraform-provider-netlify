package provider

import (
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
	CustomDomain string `json:"custom_domain"`
	Name         string `json:"name"`
}

func (n NetlifyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+n.Token)
	return n.T.RoundTrip(req)
}

func (c *NetlifyClient) Get(path string) (*http.Response, error) {
	fullURL := c.BaseURL.String() + path

	res, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *NetlifyClient) ListSites() (*[]Site, error) {
	res, err := c.Get("/sites")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var sites []Site
	err = json.Unmarshal([]byte(body), &sites)
	if err != nil {
		return nil, err
	}

	return &sites, nil
}

func (c *NetlifyClient) GetSite(siteId string) (*Site, error) {
	res, err := c.Get("sites/" + siteId)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var site Site
	err = json.Unmarshal([]byte(body), &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

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
