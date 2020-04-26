package bitflyer

import "net/http"

const baseUrl = "https://api.bitflyer.com/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

func (api *APIClient) header(method, endpoint string, body []byte) map[string]string {
	// timtestamp := strconv.Format
}
