package amazonring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/google/uuid"
)

const BaseURL 				 = "https://api.ring.com"
const ClientApiBaseURL = "https://api.ring.com/clients_api/"
const DeviceApiBaseURL = "https://api.ring.com/devices/v1/"
const AppApiBaseURL 	 = "https://app.ring.com/api/v1/"
const AuthURL      		 = "https://oauth.ring.com/oauth/token"
const authClientID 		 = "ring_official_android"
const TokenScope   		 = "client"

type Auth struct {
	Scope				 			string		`json:"scope"`
	Username					string 		`json:"username,omitempty"`
	Password  				string 		`json:"password,omitempty"`
	TwoFactorAuthCode string 		`json:"2fa-code,omitempty"`
}

type Token struct {
	AccessToken	 			string		`json:"access_token,omitempty"`
	RefreshToken 			string 	 	`json:"refresh_token,omitempty"`
	TokenType 	 			string 	 	`json:"token_type"`
	ExpiresIn 	 			float64 	`json:"expires_in"`
	Expires 		 			time.Time
}

type Client struct {
	Auth				*Auth
	HttpClient 	*http.Client
	Token				*Token
	baseURL 	   string
	clientID 	 	 string
	grantType 	 string
	hardwareID 	 uuid.UUID
}

type ClientOpt func(*Client) error

// Creates a new HTTP client
func NewClient(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	auth := &Auth{
		Scope: TokenScope,
	}

	client := &Client{
		baseURL: BaseURL,
		HttpClient: httpClient,
		clientID: authClientID,
		Auth: auth,
		Token: nil,
	}

	return client, nil
}

func NewClientWithOptions(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c, err := NewClient(httpClient)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// SetBaseURL is a config option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		c.baseURL = bu
		return nil
	}
}

//Sets HTTP headers for all requests
func (c *Client) setHeaders(req *http.Request) {
	if c.Token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.Token.TokenType, c.Token.AccessToken))
	} else {
		req.Header.Set("2fa-support", "true")
		req.Header.Set("2fa-code", c.Auth.TwoFactorAuthCode)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hardware_id", c.hardwareID.String())
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Client) get(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	c.setHeaders(req)
	resp, err := c.HttpClient.Do(req); Check(err);
	rb, _ := ioutil.ReadAll(resp.Body)

	return rb, nil
}

func (c Client) post(url string, body []byte) (*http.Response, map[string]interface{}, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	c.setHeaders(req)
	res, err := c.HttpClient.Do(req)
	if err != err {
		panic(err)
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	var resBody map[string]interface{}
	json.Unmarshal(b, &resBody)

	return res, resBody, nil
}

func (c Client) loadToken() {
	
}