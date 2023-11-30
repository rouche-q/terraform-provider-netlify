package netlify

import (
	"bytes"
	"net/http"
)

type CurrentUser struct {
	Id          string `json:"id"`
	Uid         string `json:"uid"`
	Slug        string `json:"slug"`
	FullName    string `json:"full_name"`
	AvatarUrl   string `json:"avatar_url"`
	Email       string `json:"email"`
	AffiliateId string `json:"affiliate_id"`
	SiteCount   int    `json:"site_count"`
	CreatedAt   string `json:"created_at"`
	LastLogin   string `json:"last_login"`
}

func (c *NetlifyClient) GetCurrentUser() (*CurrentUser, error) {
	var currentUser CurrentUser

	reqDo := Request{
		Method: http.MethodGet,
		Path:   "user/",
		Body:   &bytes.Buffer{},
	}
	err := c.Do(reqDo, &currentUser)
	if err != nil {
		return nil, err
	}

	return &currentUser, nil
}
