package netlify

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type EnvVar struct {
	Key      string        `json:"key"`
	Scopes   []string      `json:"scopes,omitempty"`
	Values   []EnvVarValue `json:"values,omitempty"`
	IsSecret bool          `json:"is_secret"`
}

type EnvVarValue struct {
	Id               string `json:"id,omitempty"`
	Value            string `json:"value"`
	Context          string `json:"context,omitempty"`
	ContextParameter string `json:"context_parameter,omitempty"`
}

func (c *NetlifyClient) CreateEnvVar(accountSlug string, siteId string, envVar EnvVar) (*EnvVar, error) {
	jsonValue, err := json.Marshal([]EnvVar{envVar})
	if err != nil {
		return nil, err
	}

	reqDo := Request{
		Method: http.MethodPost,
		Path:   "accounts/" + accountSlug + "/env",
		Body:   bytes.NewBuffer(jsonValue),
		Query: map[string]string{
			"site_id": siteId,
		},
	}

	var resEnvVars []EnvVar
	err = c.Do(reqDo, &resEnvVars)
	if err != nil {
		return nil, err
	}

	return &resEnvVars[0], nil
}

func (c *NetlifyClient) GetEnvVar(accountSlug string, siteId string, key string) (*EnvVar, error) {
	reqDo := Request{
		Method: http.MethodGet,
		Path:   "accounts/" + accountSlug + "/env/" + key,
		Body:   &bytes.Buffer{},
		Query: map[string]string{
			"site_id": siteId,
		},
	}

	var resEnvVar EnvVar
	err := c.Do(reqDo, &resEnvVar)
	if err != nil {
		return nil, err
	}
	return &resEnvVar, nil
}

func (c *NetlifyClient) UpdateEnvVar(accountSlug string, siteId string, key string, envVar EnvVar) (*EnvVar, error) {
	jsonValue, err := json.Marshal(envVar)
	if err != nil {
		return nil, err
	}

	reqDo := Request{
		Method: http.MethodPut,
		Path:   "accounts/" + accountSlug + "/env/" + key,
		Body:   bytes.NewBuffer(jsonValue),
		Query: map[string]string{
			"site_id": siteId,
		},
	}

	var resEnvVars EnvVar
	err = c.Do(reqDo, &resEnvVars)
	if err != nil {
		return nil, err
	}

	return &resEnvVars, nil
}

func (c *NetlifyClient) DeleteEnvVar(accountSlug string, siteId string, key string) error {
	reqDo := Request{
		Method: http.MethodDelete,
		Path:   "accounts/" + accountSlug + "/env/" + key,
		Body:   &bytes.Buffer{},
		Query: map[string]string{
			"site_id": siteId,
		},
	}

	return c.Do(reqDo, nil)
}
