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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")

		var responseBody map[string]interface{}
		b, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(b, &responseBody)
		
		// Check test credentials
		if responseBody["username"] == username && responseBody["password"] == password && responseBody["2fa-code"] == twoFactorAuthCode {
			w.WriteHeader(http.StatusOK)
		} else {
			t.Fatal()
		}
	})

	oauth := Oauth{
		Username: username,
		Password: password,
		TwoFactorAuthCode: twoFactorAuthCode,
	}

	Auth(context.Background(), config, oauth)
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