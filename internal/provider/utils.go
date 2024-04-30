package provider

import "net/http"

func AddAuthenticationToRequest(req *http.Request, apiKey string) {
	bearer := "Bearer " + apiKey
	req.Header.Add("Authorization", bearer)
}
