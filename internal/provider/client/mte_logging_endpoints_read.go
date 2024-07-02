package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ReadMTELoggingEndpoints() (*MTELoggingEndpointsDto, error) {
	httpRes, err := c.initiateRequest(
		http.MethodGet,
		"/v1/admin/logging",
		nil)

	if err != nil {
		return nil, &AltitudeClientError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode != 200 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		return nil, &AltitudeClientError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s with body %s.", httpRes.Status, body),
		}
	}

	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)

	if err != nil {
		return nil, &AltitudeClientError{
			shortMessage: "Body Read Error",
			detail:       "Unable to read response body",
		}
	}

	var dto MTELoggingEndpointsDto
	err = json.Unmarshal(body, &dto)

	if err != nil {
		return nil, &AltitudeClientError{
			shortMessage: "Body Read Error",
			detail:       "Unable to parse JSON body from Altitude response: " + string(body),
		}
	}

	return &dto, nil
}
