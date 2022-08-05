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

	"github.com/google/uuid"
	ring "github.com/jonpmay/go-client-amazon-ring/internal/amazonring"
)

const authURL        = ring.AuthBaseURL
const oauthClientID = "ring_official_android"
const oauthScope    = "client"

type Oauth struct {
	ClientID					string		`json:"client_id"`
	Scope							string		`json:"scope"`
	GrantType					string		`json:"grant_type,omitempty"`
	Username					string 		`json:"username,omitempty"`
	Password  				string 		`json:"password,omitempty"`
	TwoFactorAuthCode string 		`json:"2fa-code,omitempty"`
	AccessToken				string		`json:"access_token,omitempty"`
	RefreshToken 			string 		`json:"refresh_token,omitempty"`
	HardwareId 				uuid.UUID `json:"hardware_id"`
}

type Config struct {
	client   	 *http.Client
	ClientID 	 string
	AuthURL    *url.URL
	Scope 	 	 string
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

func Auth(ctx context.Context, config *Config, oauth Oauth) {
	oauth.ClientID = oauthClientID
	oauth.Scope = oauthScope
	oauth.HardwareId = uuid.New()

	// If refresh token is not present, enter username/password
	if oauth.RefreshToken != "" {
		oauth.GrantType = "refresh_token"
	} else {
		oauth.GrantType = "password"
	}
	
	if oauth.TwoFactorAuthCode == "" {
		if oauth.Username == "" {
			oauth.Username = GetInput("Enter Ring mail address: ")
		}
		if oauth.Password == "" {
			oauth.Password = GetInput("Enter Ring password: ")
		}
	}

	// Build HTTP request
	data, err := json.Marshal(oauth);	Check(err)
	req, err := http.NewRequest("POST", config.AuthURL.String(), bytes.NewReader(data)); Check(err)

	req.Header.Set("2fa-support", "true")
	req.Header.Set("2fa-code", oauth.TwoFactorAuthCode)
	req.Header.Set("hardware_id", oauth.HardwareId.String())
	req.Header.Set("Content-Type", "application/json")

	resp, err := config.client.Do(req); Check(err)
	defer resp.Body.Close()

	//DEBUG
	fmt.Println(req)
	//DEBUG

	// Marshall response body
	var responseBody map[string]interface{}
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &responseBody)

	if resp.StatusCode == 412 || (resp.StatusCode == 400 && strings.HasPrefix(responseBody["error"].(string), "Verification Code")) {
		if val, ok := responseBody["tsv_state"]; ok {
			if val == "sms" {
				oauth.TwoFactorAuthCode = GetInput("Please enter the code sent to " + responseBody["phone"].(string) + ": ")
				Auth(ctx, config, oauth)
			}
		}
	}
	// DEBUG
	fmt.Println(responseBody)
	// DEBUG
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
