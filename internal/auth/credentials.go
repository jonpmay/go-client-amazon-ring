package auth

import (
	"context"
	"fmt"
//	"os"

//	"github.com/BurntSushi/toml"
)

const credentialDirectory = ".ring"
const credentialFileName = "credentials"

// Saves API credentials to a local TOML file in ${HOME}/.ring/
func SaveCredentials(ctx context.Context, authInfo *AuthInfo) {
	fmt.Println("Saving credentials...")
}