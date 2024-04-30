package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateMTEConfigInput struct {
	Config        MTEConfigDto
	EnvironmentId string
}

func (c *Client) CreateMTEConfig(
	input CreateMTEConfigInput,
) error {
	jsonBody, err := json.Marshal(input.Config)
	if err != nil {
		return &AltitudeClientError{
			"Input Error.",
			"Input unable to be JSON encoded.",
		}
	}

	httpRes, err := c.initiateRequest(
		http.MethodPost,
		fmt.Sprintf("/v1/environment/%s/mte/altitude-config", input.EnvironmentId),
		bytes.NewBuffer(jsonBody))

	if err != nil {
		return &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 409 {
		return &AltitudeClientError{
			shortMessage: "Environment ID Conflict",
			detail:       "This environment already has an associated config block.",
		}
	}

	if httpRes.StatusCode != 201 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		return &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s with body %s.", httpRes.Status, body),
		}
	}
	return nil
}
