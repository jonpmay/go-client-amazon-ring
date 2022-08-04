package amazonring

import (
	"net/http"
)

const BaseURL = "https://api.ring.com"
const ClientApiBaseURL = "https://api.ring.com/clients_api/"
const DeviceApiBaseURL = "https://api.ring.com/devices/v1/"
const AppApiBaseURL 	 = "https://app.ring.com/api/v1/"
const AuthBaseURL			 = "https://oauth.ring.com/oauth/token"

type Client struct {
	baseURL 	  string
	HttpClient 	*http.Client
}

// NewClient crreates a new HTTP client
func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{
		baseURL: BaseURL,
		HttpClient: httpClient,
	}

	return client, nil
}
