package client

import (
	"fmt"
	"io"
	"net/http"
)

type DeleteMTEConfigInput struct {
	EnvironmentId string
}

func (c *Client) DeleteMTEConfig(
	input DeleteMTEConfigInput,
) error {
	httpRes, err := c.initiateRequest(
		http.MethodDelete,
		fmt.Sprintf("/v2/environment/%s/mte/altitude-config", input.EnvironmentId),
		nil)

	if err != nil {
		return &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}
	if httpRes.StatusCode == 404 {
		return &AltitudeClientError{
			shortMessage: "Environment ID not found",
			detail:       fmt.Sprintf("The Environment %s does not have associated config.", input.EnvironmentId),
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
