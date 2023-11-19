package netlify

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Site struct {
	Id           string `json:"id"`
	CustomDomain string `json:"custom_domain"`
	Name         string `json:"name"`
	Url          string `json:"url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	State        string `json:"state"`
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

type SiteRequest struct {
	Name         string     `json:"name"`
	CustomDomain string     `json:"custom_domain"`
	Repo         Repository `json:"repo"`
}

func (c *NetlifyClient) CreateSite(req SiteRequest) (*Site, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqDo := Request{
		Method: http.MethodPost,
		Path:   "sites/",
		Body:   bytes.NewBuffer(jsonValue),
	}

	var resSite Site
	err = c.Do(reqDo, &resSite)
	if err != nil {
		return nil, err
	}

	return &resSite, nil
}

func (c *NetlifyClient) GetSite(siteId string) (*Site, error) {
	var site Site

	reqDo := Request{
		Method: http.MethodGet,
		Path:   "sites/" + siteId,
		Body:   &bytes.Buffer{},
	}

	err := c.Do(reqDo, &site)
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

	reqDo := Request{
		Method: http.MethodPatch,
		Path:   "sites/" + siteId,
		Body:   bytes.NewBuffer(jsonValue),
	}

	var resSite Site
	err = c.Do(reqDo, &resSite)
	if err != nil {
		return nil, err
	}

	return &resSite, nil
}

func (c *NetlifyClient) DeleteSite(siteId string) error {
	reqDo := Request{
		Method: http.MethodDelete,
		Path:   "sites/" + siteId,
		Body:   &bytes.Buffer{},
	}
	return c.Do(reqDo, nil)
}
