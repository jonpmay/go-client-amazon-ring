package auth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	teardown := setup()
	defer teardown()

	username := "testuser@test.com"
	password := "notarealpassword"
	twoFactorAuthCode := "123456"

	Data := Token{
		AccessToken: "eyJhbGciOiJIUzUxMiIsImprdSI6Ii9vYXV0aC9pbnRlcm5hbC9qd2tzIiwia2lkIjoiYzEyODEwMGIiLCJ0eXAiOiJKV1QifQ.eyJhcHBfaWQiOiJyaW5nX29mZmljaWFsX2FuZHJvaWQiLCJjaWQiOiJyaW5nX29mZmljaWFsX2FuZHJvaWQiLCJleHAiOjE2NTk2NzQzMjQsImhhcmR3YXJlX2lkIjoiMzk2Y2E5MjEtNWU3ZS00YWY5LWJkYTMtYTYyZDM3YjVmMTViIiwiaWF0IjoxNjU5NjcwNzI0LCJpc3MiOiJSaW5nT2F1dGhTZXJ2aWNlLXByb2Q6dXMtZWFzdC0xOmYxMDJkZTliIiwicm5kIjoidi1saHZPZ2ZvdnlObGciLCJzY29wZXMiOlsiY2xpZW50Il0sInNlc3Npb25faWQiOiJmODBlZWZiYy02ZWZhLTQ4NjMtYjNjMy1jMzUzM2NhN2VmODgiLCJ1c2VyX2lkIjoxMjM0NTY3OH0.wibv9Ktr7fdC1M2Px_86JECjKWY3F5HT4x6c6augHicNxI9l9Yk3PGGPd297NM2UitQiUBSlcvj6mMGUW0mdCA",
		ExpiresIn: 3600,
		RefreshToken: "eyJhbGciOiJIUzUxMiIsImprdSI6Ii9vYXV0aC9pbnRlcm5hbC9qd2tzIiwia2lkIjoiYzEyODEwMGIiLCJ0eXAiOiJKV1QifQ.eyJpYXQiOjE2NTk2NzA3MjQsImlzcyI6IlJpbmdPYXV0aFNlcnZpY2UtcHJvZDp1cy1lYXN0LTE6ZjEwMmRlOWIiLCJyZWZyZXNoX2NpZCI6InJpbmdfb2ZmaWNpYWxfYW5kcm9pZCIsInJlZnJlc2hfc2NvcGVzIjpbImNsaWVudCJdLCJyZWZyZXNoX3VzZXJfaWQiOjEyMzQ1Njc4LCJybmQiOiJ2R3VJOXg4NktYdVN6ZyIsInNlc3Npb25faWQiOiJhNDhiZGU1NC0yZDE5LTQ3ZWItOWNiYi00MjFmNWZkMjM1ODUiLCJ0eXBlIjoicmVmcmVzaC10b2tlbiJ9.eVEVQVnZyPLL4nWg124OVGxe1Li69S9nERLVeikodQSJhZM9YZ-M6qmLSxYgMebTOhlgcvkWN_Gx4WdOKGXNVw",
		Scope: "client",
		TokenType: "Bearer",
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		var responseBody map[string]interface{}
		b, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(b, &responseBody)
		
		// Check test credentials
		if responseBody["username"] == username && responseBody["password"] == password && responseBody["2fa-code"] == twoFactorAuthCode {
			rb, _ := json.Marshal(Data)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(rb)
		} else {
			t.Fatal()
		}
	})

	authInfo := &AuthInfo{
		Username: username,
		Password: password,
		TwoFactorAuthCode: twoFactorAuthCode,
	}

	config.Auth(context.Background(), authInfo)
}

var (
	mux *http.ServeMux
	server *httptest.Server
	config *Config
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	config, _ = NewConfigWithOptions(nil, SetAuthURL(server.URL))

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/" + path)
	if err != nil {
		panic (err)
	}
	return string(b)
}
