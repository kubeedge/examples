package utils

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

// ReadSecretKey is a helper function that will return the value of
// the key passed in as an argument. The keys in this case are secrets
// containing the Twitter credentials for authentication.
func ReadSecretKey(key string) (string, error) {
	secretKey := fmt.Sprintf("/etc/secret/%s", key)
	keyData, err := readInternal(secretKey)
	if err != nil {
		return "", err
	}
	return keyData, nil
}

func readInternal(key string) (string, error) {
	cmd := exec.Command("cat", key)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(out), nil
}