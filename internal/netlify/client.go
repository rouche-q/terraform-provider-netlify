package netlify

import (
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

func (n NetlifyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+n.Token)
	req.Header.Set("Content-Type", "application/json")
	return n.T.RoundTrip(req)
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
