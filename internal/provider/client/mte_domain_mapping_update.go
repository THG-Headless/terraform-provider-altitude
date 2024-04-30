package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UpdateMteDomainMappingInput struct {
	Config MTEDomainMappingDto
}

func (c *Client) UpdateMteDomainMapping(
	input UpdateMteDomainMappingInput,
) (string, error) {
	jsonBody, err := json.Marshal(input.Config)
	if err != nil {
		return "", &AltitudeClientError{
			"Input Error.",
			"Input unable to be JSON encoded.",
		}
	}

	httpRes, err := c.initiateRequest(
		http.MethodPut,
		"/v1/mte/domain-mapping",
		bytes.NewBuffer([]byte(jsonBody)))

	if err != nil {
		return "", &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode != 201 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		return "", &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-201 response of %s with body %s.", httpRes.Status, body),
		}
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", &AltitudeClientError{
			shortMessage: "Body Read Error",
			detail:       "Unable to read response body",
		}
	}

	return string(body[:]), nil
}