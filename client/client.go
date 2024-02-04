package client

import (
	"github.com/go-resty/resty/v2"
)

func CreateRestyClient(accessToken string, baseUrl string) (*resty.Client, error) {
	client := resty.New()
	client.SetAuthScheme("Bearer")
	client.SetAuthToken(accessToken)
	client.SetBaseURL(baseUrl)

	return client, nil
}
