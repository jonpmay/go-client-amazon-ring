package amazonring

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// TO DO
func (c *Client) RefeshTokenGrant(t *Token) *Token {
	c.grantType = "refresh_token"
	return nil
}

func (c *Client) PasswordGrant() Token {
	c.grantType = "password"
	t := Token{}
	
	if c.Auth.Username == "" || c.Auth.Password == "" {
		c.Auth.Username = GetInput("Enter Ring email address: ")
		c.Auth.Password = GetInput("Enter Ring password: ")
	}
	b := map[string]interface{}{
		"client_id":  c.clientID,
		"scope": 			c.Auth.Scope,
		"username": 	c.Auth.Username,
		"password": 	c.Auth.Password,
		"grant_type": c.grantType,
	}

	reqBody, err := json.Marshal(b); Check(err)
	res, resBody, err := c.post(c.baseURL, reqBody); Check(err)

	if res.StatusCode == 200 {
		t = Token{
			AccessToken: resBody["access_token"].(string),
			RefreshToken: resBody["refresh_token"].(string),
			ExpiresIn: resBody["expires_in"].(float64),
			Expires: time.Now().Add(time.Second * time.Duration(resBody["expires_in"].(float64))),
			TokenType: resBody["token_type"].(string),
		}
	}	else if res.StatusCode == 412 || (res.StatusCode == 400 && strings.HasPrefix(resBody["error"].(string), "Verification Code")) {
		t = c.twoFactorAuth(context.WithValue(context.Background(), "phone", resBody["phone"].(string)))
	}
	
	return t
}

//Retrieves one string input via stdin
func GetInput(message string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(message)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.Trim(input, "\n")
}

func (c *Client) twoFactorAuth(ctx context.Context) Token {
	c.Auth.TwoFactorAuthCode = GetInput(fmt.Sprintf("Please enter the code sent to %s:", ctx.Value("phone").(string)))
	b := map[string]interface{}{
		"client_id":  c.clientID,
		"scope": 			c.Auth.Scope,
		"username": 	c.Auth.Username,
		"password": 	c.Auth.Password,
		"grant_type": c.grantType,
	}

	reqBody, _ := json.Marshal(b)
	res, resBody, err := c.post(c.baseURL, reqBody)
	if err != nil {
		panic(err)
	}
	
	if res.StatusCode == 200 {
		t := Token{
			AccessToken: resBody["access_token"].(string),
			RefreshToken: resBody["refresh_token"].(string),
			ExpiresIn: resBody["expires_in"].(float64),
			Expires: time.Now().Add(time.Second * time.Duration(resBody["expires_in"].(float64))),
			TokenType: resBody["token_type"].(string),
		}
		return t
	} else {
		panic(fmt.Sprintf("Failed to obtain token: %d %s", res.StatusCode, res.Body))
	}
}
