package netlify

import (
	"bytes"
	"net/http"
)

type DeployKey struct {
	Id        string `json:"id"`
	Key       string `json:"public_key"`
	CreatedAt string `json:"created_at"`
}

func (c *NetlifyClient) CreateDeployKey() (*DeployKey, error) {
	var key DeployKey

	reqDo := Request{
		Method: http.MethodPost,
		Path:   "deploy_keys/",
		Body:   &bytes.Buffer{},
	}
	err := c.Do(reqDo, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *NetlifyClient) GetDeployKey(keyId string) (*DeployKey, error) {
	var key DeployKey

	reqDo := Request{
		Method: http.MethodGet,
		Path:   "deploy_keys/" + keyId,
		Body:   &bytes.Buffer{},
	}
	err := c.Do(reqDo, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *NetlifyClient) DeleteDeployKey(keyId string) error {
	reqDo := Request{
		Method: http.MethodDelete,
		Path:   "deploy_keys/" + keyId,
		Body:   &bytes.Buffer{},
	}
	return c.Do(reqDo, nil)
}
