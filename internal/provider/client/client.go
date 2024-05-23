package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Mode string

const (
	Production Mode = "Production"
	UAT        Mode = "UAT"
	Local      Mode = "Local"
)

func (m *Mode) IsValid() bool {
	return string(*m) == "Production" || string(*m) == "UAT" || string(*m) == "Local"
}

type AltitudeClientVariables struct {
	baseUrl  string
	audience string
	issuer   string
}

type Client struct {
	clientVariables AltitudeClientVariables
	httpClient      *http.Client
	token           string
}

func New(
	clientId string,
	clientSecret string,
	mode Mode) (*Client, error) {
	c := new(Client)
	c.httpClient = http.DefaultClient
	switch mode {
	case Production:
		c.clientVariables = AltitudeClientVariables{
			baseUrl:  "https://api.platform.thgaltitude.com",
			audience: "https://api.platform.thgaltitude.com/",
			issuer:   "https://thgaltitude.eu.auth0.com",
		}
	case UAT:
		c.clientVariables = AltitudeClientVariables{
			baseUrl:  "https://uat-api.platform.thgaltitude.com",
			audience: "https://platform.thgaltitude.co.uk/api/",
			issuer:   "https://dev-thgaltitude.eu.auth0.com",
		}
	case Local:
		c.clientVariables = AltitudeClientVariables{
			baseUrl:  "http://localhost:8080",
			audience: "http://localhost:8080/",
			issuer:   "https://vega-local-dev.uk.auth0.com",
		}
	}
	err := c.generateAuthToken(clientId, clientSecret)
	if err != nil {
		return nil, &AltitudeClientError{
			"The Altitude Client is unable to generate an auth token",
			err.Error(),
		}
	}
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
		fmt.Sprintf("%s%s", c.clientVariables.baseUrl, path),
		body,
	)
	if err != nil {
		return nil, &AltitudeClientError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodGet {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	c.addAuthenticationToRequest(httpReq)
	return c.httpClient.Do(httpReq)
}

type AuthDto struct {
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthResBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *Client) generateAuthToken(
	clientId string,
	clientSecret string,
) error {
	authDto := AuthDto{
		Audience:     c.clientVariables.audience,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	reqBody, err := json.Marshal(authDto)
	if err != nil {
		return err
	}

	if c.httpClient == nil {
		return &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       "Default HTTP Client is undefined",
		}
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/oauth/token", c.clientVariables.issuer),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s with body %s.", resp.Status, body),
		}
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var body AuthResBody
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return err
	}

	c.token = body.AccessToken
	return nil
}
