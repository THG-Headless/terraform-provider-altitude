package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseUrl string
	httpClient *http.Client
	token string
}

func New(
	baseUrl string,
	clientId string,
	clientSecret string,
	audience string) (*Client, error) {
	c, err := new(Client).generateAuthToken(clientId, clientSecret, audience)
	if err != nil {
		return nil, &AltitudeClientError {
			"The Altitude Client is unable to generate an auth token",
			err.Error(),
		}
	}
	c.baseUrl = baseUrl
	return c, nil
}

func (c *Client) addAuthenticationToRequest(req *http.Request) {
	bearer := "Bearer " + c.token
	req.Header.Add("Authorization", bearer)
}

func (c *Client) initiateRequest(
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, &AltitudeClientError{
			shortMessage: "Incorrect Path Format",
			detail:       fmt.Sprintf("The path %s should be specified with a prefixed slash.", path),
		}
	}
	httpReq, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", c.baseUrl, path),
		body,
	)
	if err != nil {
		return nil, &AltitudeClientError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}
	c.addAuthenticationToRequest(httpReq)
	return c.httpClient.Do(httpReq)
}

type AuthDto struct {
	Audience string `json:"audience"`
	GrantType string `json:"grant_type"`
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthResBody struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func (c *Client) generateAuthToken(
	clientId string,
	clientSecret string,
	audience string,
) (*Client, error) {
	authDto := AuthDto{
		Audience: audience,
		ClientId: clientId,
		ClientSecret: clientSecret,
		GrantType: "client_credentials",
	}

	reqBody, err := json.Marshal(authDto)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("https://%s/oauth/token", "thgaltitude.eu.auth0.com"),
		"application/json",
		bytes.NewBuffer([]byte(reqBody)),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s with body %s.", resp.Status, body),
		}
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var body AuthResBody
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return nil, err
	}

	c.token = body.AccessToken
	return c, nil
}