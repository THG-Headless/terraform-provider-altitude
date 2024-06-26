package provider

import "fmt"

type KVStoreitemError struct {
	detail       string
	shortMessage string
}

func (e *KVStoreitemError) Error() string {
	return fmt.Sprintf("%s\n%s", e.shortMessage, e.detail)
}

type AltitudeApiError struct {
	detail       string
	shortMessage string
}

func (e *AltitudeApiError) Error() string {
	return fmt.Sprintf("%s\n%s", e.shortMessage, e.detail)
}
