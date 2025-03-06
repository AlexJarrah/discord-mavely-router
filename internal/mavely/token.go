package mavely

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal"
	"gitlab.com/AlexJarrah/discord-mavely-router/internal/filesystem"
)

// fetchToken fetches a new access token using the specified username and password
func fetchToken() (*TokenData, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "PItSrEo35MYmLjhY6wJp8sCQAQRRWYxr")
	data.Set("username", internal.Configuration.Mavely.Username)
	data.Set("password", internal.Configuration.Mavely.Password)
	data.Set("scope", "openid profile email offline_access")
	data.Set("audience", "https://auth.mave.ly/")

	req, err := http.NewRequest("POST", "https://auth.mave.ly/oauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-gpc", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var tokenData TokenData
	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token data: %v", err)
	}
	if tokenData.AccessToken == "" {
		return nil, fmt.Errorf("no access token received from API")
	}

	tokenData.FetchedAt = time.Now().UnixNano() / 1e6

	// Save token data to file
	file, err := os.Create(filepath.Join(filesystem.DataDirectory, ".token"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token file: %v", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to write token data to file: %v", err)
	}

	return &tokenData, nil
}

// refreshToken refreshes an existing token using the refresh token
func refreshToken(refreshToken string) (*TokenData, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", "PItSrEo35MYmLjhY6wJp8sCQAQRRWYxr")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://auth.mave.ly/oauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-gpc", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var tokenData TokenData
	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token data: %v", err)
	}
	if tokenData.AccessToken == "" {
		return nil, fmt.Errorf("no access token received from API")
	}

	tokenData.FetchedAt = time.Now().UnixNano() / 1e6

	// Save token data to file
	file, err := os.Create(filepath.Join(filesystem.DataDirectory, ".token"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token file: %v", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to write token data to file: %v", err)
	}

	return &tokenData, nil
}

// GetToken retrieves a valid token, by cache or generating a new one
func GetToken() (*TokenData, error) {
	if _, err := os.Stat(filepath.Join(filesystem.DataDirectory, ".token")); err == nil {
		// Token file exists, attempt to read it
		file, err := os.Open(filepath.Join(filesystem.DataDirectory, ".token"))
		if err != nil {
			log.Printf("Failed to open token file, fetching new token: %v", err)
			return fetchToken()
		}
		defer file.Close()

		var tokenData TokenData
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&tokenData)
		if err != nil {
			log.Printf("Failed to decode token data, fetching new token: %v", err)
			return fetchToken()
		}

		// Validate token
		now := time.Now().UnixNano() / 1e6
		timeElapsed := now - tokenData.FetchedAt
		expiresInMilliseconds := int64(tokenData.ExpiresIn) * 1000
		timeLeft := expiresInMilliseconds - timeElapsed

		if timeLeft > 60000 {
			return &tokenData, nil
		}

		log.Printf("Refreshing token...")
		newToken, err := refreshToken(tokenData.RefreshToken)
		if err == nil {
			return newToken, nil
		}
		log.Printf("Failed to refresh token, fetching a new one: %v", err)
		return fetchToken()
	}

	// No token file exists, fetch a new one
	log.Printf("No token found, fetching a new one...")
	return fetchToken()
}
