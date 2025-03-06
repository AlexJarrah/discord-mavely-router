package mavely

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateLink creates an affiliate link using the provided token and URL
func CreateLink(token, url string) (string, map[string]any, error) {
	requestBody, err := json.Marshal(map[string]any{
		"query":     "mutation ($v1:String!){createAffiliateLink(url:$v1){id,link,metaDescription,metaTitle,metaImage,metaUrl,metaLogo,metaSiteName,metaVideo,brand{id,name,slug},originalUrl,canonicalLink,attributionUrl}}",
		"variables": map[string]any{"v1": url},
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://mavely.live/", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("client-name", "@mavely/creator-app")
	req.Header.Set("client-revision", "9cac5acc")
	req.Header.Set("client-version", "1.0.3")
	req.Header.Set("content-type", "application/json")
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
		return "", nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse response: %v", err)
	}

	data, ok := result["data"].(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("invalid response structure")
	}

	createAffiliateLink, ok := data["createAffiliateLink"].(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("missing createAffiliateLink field")
	}

	link, ok := createAffiliateLink["link"].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing link field")
	}

	return link, result, nil
}
