package auth

import (
	"bytes"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const credentialDirectory = ".ring"
const credentialFileName = "credentials"

type Credentials struct {
	Username string 
	Password string
	AccessToken string
	RefreshToken string
	ExpiresIn float64
	Expires time.Time
}

// Creates a new set of credentials
func NewCredentials(ai AuthInfo, t Token) Credentials {
	c := &Credentials{
		Username: ai.Username,
		Password: ai.Password,
		AccessToken: t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresIn: t.ExpiresIn,
		Expires: t.Expires,
	}

	return *c
}

// Encodes the passed credentials as a TOML configuration file
func (c Credentials) EncodeTOMLFile() []byte {
	var b bytes.Buffer
	toml.NewEncoder(&b).Encode(c)
	return b.Bytes()
}

// Saves the specified credentials to a local file
func (c Credentials) SaveCredentials(d []byte) {
	uhd, err := os.UserHomeDir(); Check(err)
	os.MkdirAll(uhd + "/" + credentialDirectory, 0700)
	os.WriteFile(uhd + "/" + credentialDirectory + "/" + credentialFileName, d, 0600)
}