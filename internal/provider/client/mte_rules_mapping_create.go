package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateMteRulesMappingInput struct {
	Config MTERulesMappingDto
}

func (c *Client) CreateMteRulesMapping(
	input CreateMteRulesMappingInput,
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
		"/v1/mte/rules-mapping",
		bytes.NewBuffer(jsonBody))

	if err != nil {
		return &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 409 {
		return &AltitudeClientError{
			shortMessage: "Domain Conflict",
			detail:       "This environment already has an associated config block.",
		}
	}

	if httpRes.StatusCode != 201 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		return &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-201 response of %s with body %s.", httpRes.Status, body),
		}
	}

	return nil
}
