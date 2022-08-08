package auth

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	ring "github.com/jonpmay/go-client-amazon-ring/internal/amazonring"
)

const authURL      = ring.AuthBaseURL
const authClientID = "ring_official_android"
const tokenScope   = "client"

type AuthInfo struct {
	ClientID					string		`json:"client_id"`
	Scope				 			string		`json:"scope"`
	GrantType					string		`json:"grant_type,omitempty"`
	Username					string 		`json:"username,omitempty"`
	Password  				string 		`json:"password,omitempty"`
	TwoFactorAuthCode string 		`json:"2fa-code,omitempty"`
	HardwareId 				uuid.UUID `json:"hardware_id"`
	AccessToken	 			string		`json:"access_token,omitempty"`
	RefreshToken 			string 	 	`json:"refresh_token,omitempty"`
	TokenType 	 			string 	 	`json:"token_type"`
	ExpiresIn 	 			float64 	`json:"expires_in"`
	Expires 		 			time.Time
}

type Token struct {
	AccessToken	 			string		`json:"access_token,omitempty"`
	Scope				 			string		`json:"scope"`
	RefreshToken 			string 	 	`json:"refresh_token,omitempty"`
	TokenType 	 			string 	 	`json:"token_type"`
	ExpiresIn 	 			float64 	`json:"expires_in"`
	Expires 		 			time.Time
}

type Config struct {
	client   	 *http.Client
	ClientID 	 string
	AuthURL    *url.URL
}

func NewConfig(httpClient *http.Client) *Config {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(authURL)
	c := &Config{
		client: httpClient,
		AuthURL: baseURL,
	}

	return c
}

type ConfigOpt func(*Config) error

func NewConfigWithOptions(httpClient *http.Client, opts ...ConfigOpt) (*Config, error) {
	c := NewConfig(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// SetAuthURL is a config option for setting the base URL.
func SetAuthURL(bu string) ConfigOpt {
	return func(c *Config) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.AuthURL = u
		return nil
	}
}

func (config *Config) Auth(ctx context.Context, authInfo *AuthInfo) *Token {
	authInfo.ClientID = authClientID
	authInfo.Scope = tokenScope
	authInfo.HardwareId = uuid.New()
	token := &Token{}

	// If refresh token is not present, enter username/password
	if authInfo.RefreshToken != "" {
		authInfo.GrantType = "refresh_token"
	} else {
		authInfo.GrantType = "password"
		if authInfo.TwoFactorAuthCode == "" {
			if authInfo.Username == "" {
				authInfo.Username = GetInput("Enter Ring email address: ")
			}
			if authInfo.Password == "" {
				authInfo.Password = GetInput("Enter Ring password: ")
			}
		}
	}
	
	count := 0
	for {
		req := &http.Request{}
		*req = config.createRequest(*authInfo)
		resp, jsonBody := config.makeRequest(req)

		switch {
		case resp.StatusCode == 412:
			if val, ok := jsonBody["tsv_state"]; ok {
				if val == "sms" {
					authInfo.TwoFactorAuthCode = GetInput("Please enter the code sent to " + jsonBody["phone"].(string) + ": ")
				}
			}
		case resp.StatusCode == 400 && strings.HasPrefix(jsonBody["error"].(string), "Verification Code"):
			authInfo.TwoFactorAuthCode = GetInput("Please enter the code sent to " + jsonBody["phone"].(string) + ": ")
		case resp.StatusCode == 429:
			panic("Received HTTP StatusCode 429: Too many requests")
		case resp.StatusCode == 200:
			token.AccessToken = jsonBody["access_token"].(string)
			token.RefreshToken = jsonBody["refresh_token"].(string)
			token.ExpiresIn = jsonBody["expires_in"].(float64)
			token.Expires = time.Now().Add(time.Second * time.Duration(authInfo.ExpiresIn))
			token.Scope = jsonBody["scope"].(string)
			token.TokenType = jsonBody["token_type"].(string)
		}
		
		if resp.StatusCode == 200 || count >= 3 {
			break
		} else {
			count++
		}
	}

	return token
}

func (config *Config) makeRequest(req *http.Request) (http.Response, map[string]interface{}) {
	var jsonBody map[string]interface{}

	resp, err := config.client.Do(req); Check(err)
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &jsonBody)
	resp.Body.Close()

	return *resp, jsonBody
}

func (config *Config) createRequest(authInfo AuthInfo) (http.Request) {
	data, err := json.Marshal(authInfo);	Check(err)
	req, err := http.NewRequest("POST", config.AuthURL.String(), bytes.NewReader(data)); Check(err)	
	req.Header.Set("2fa-support", "true")
	req.Header.Set("2fa-code", authInfo.TwoFactorAuthCode)
	req.Header.Set("hardware_id", authInfo.HardwareId.String())
	req.Header.Set("Content-Type", "application/json")

	return *req
}

func GetInput(message string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(message)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.Trim(input, "\n")
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
