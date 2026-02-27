package driver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SetDeviceValue sets the value of a device by sending an HTTP PUT request.
func SetDeviceValue(ip string, port int, urlPath string, value string) error {
	url := fmt.Sprintf("http://%s:%d%s", ip, port, urlPath)
	requestBody := bytes.NewBufferString(value)

	req, err := http.NewRequest(http.MethodPut, url, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetDeviceValue retrieves the value of a device by sending an HTTP GET request.
func GetDeviceValue(ip string, port int, urlPath string) (string, error) {
	url := fmt.Sprintf("http://%s:%d%s", ip, port, urlPath)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
