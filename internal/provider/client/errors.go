package client

import "fmt"

type AltitudeClientError struct {
	shortMessage string
	detail       string
}

func (e *AltitudeClientError) Error() string {
	return fmt.Sprintf("%s\n%s", e.shortMessage, e.detail)
}

type InternalServerError struct{}
