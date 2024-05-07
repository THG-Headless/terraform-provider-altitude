package client

import (
	"fmt"
	"io"
	"net/http"
)

type DeleteMteRulesMappingInput struct {
	Domain string
}

func (c *Client) DeleteMteRulesMapping(
	input DeleteMteRulesMappingInput,
) error {
	httpRes, err := c.initiateRequest(
		http.MethodDelete,
		fmt.Sprintf("/v1/mte/rules-mapping?domain=%s", input.Domain),
		nil,
	)

	if err != nil {
		return &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 404 {
		return &AltitudeClientError{
			shortMessage: "Domain not found",
			detail:       fmt.Sprintf("The Domain %s does not have an associated rule group.", input.Domain),
		}
	}

	if httpRes.StatusCode != 204 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		return &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-204 response of %s with body %s.", httpRes.Status, body),
		}
	}

	return nil
}
