package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GotifyResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"errorCode"`
}

func AlertServiceUnealthy(serviceName, message, gotifyToken, gotifyHost string) error {
	var gotifyResponse GotifyResponse
	var client http.Client
	requestBody := formatRequestBody(serviceName, message)
	r, err := http.NewRequest("POST", gotifyHost+"/message", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("Creating request to gotify failed: %v", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+gotifyToken)

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("Sending alert to gotify failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Reading gotify's response failed: %v", err)
	}

	err = json.Unmarshal(body, &gotifyResponse)
	if err != nil {
		return fmt.Errorf("Unmarshaling gotify's response failed: %v", err)
	}

	if gotifyResponse.Error != "" {
		return fmt.Errorf("Alert could not be created: %v, status code %d", gotifyResponse.Error, gotifyResponse.ErrorCode)
	}

	return err
}

func formatRequestBody(lxcName, message string) []byte {
	return []byte(fmt.Sprintf(`{
        "title": "Storage of %s almost full",
        "message": "%s"
    }`, lxcName, message))
}
