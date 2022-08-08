package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestSaveCredentials(t *testing.T) {
	teardown := setup()
	defer teardown()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(fixture("authResult.json"))
		w.Write(b)
	})

	data := []byte{}
	req, err := http.NewRequest("POST", config.AuthURL.String(), bytes.NewReader(data)); Check(err)

	resp, _ := config.client.Do(req)
	fmt.Println(resp)
}