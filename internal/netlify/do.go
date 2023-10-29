package netlify

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method string
	Path   string
	Query  map[string]string
	Body   *bytes.Buffer
}

func (c *NetlifyClient) Do(req Request, dest any) error {
	reqURL, err := url.Parse(c.BaseURL.String() + req.Path)
	if err != nil {
		return err
	}

	if len(req.Query) > 0 {
		for key, value := range req.Query {
			reqURL.Query().Add(key, value)
		}
	}

	httpReq, err := http.NewRequest(req.Method, reqURL.String(), req.Body)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(httpReq)
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
