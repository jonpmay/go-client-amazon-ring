package amazonring

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetInput retrieves one string input from the user via stdin
func GetInput(message string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(message)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.Trim(input, "\n")
}
