package client

import (
	"fmt"
	"io"
	"net/http"
)

type ReadMteDomainMappingInput struct {
	Domain string
}

func (c *Client) ReadMteDomainMapping(
	input ReadMteDomainMappingInput,
) (string, error) {
	httpRes, err := c.initiateRequest(
		http.MethodGet,
		fmt.Sprintf("/v1/mte/domain-mapping?domain=%s", input.Domain),
		nil,
	)

	if err != nil {
		return "", &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 404 {
		return "", &AltitudeClientError{
			shortMessage: "Environment ID not found",
			detail:       fmt.Sprintf("The Environment %s does not have associated config.", input.Domain),
		}
	}

	if httpRes.StatusCode != 200 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		return "", &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s with body %s.", httpRes.Status, body),
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