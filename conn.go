package vwap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Response represents the response structure from th RPC API.
type Response struct {
	Result struct {
		Response struct {
			ResponseBase struct {
				Data string `json:"data"`
			} `json:"responseBase"`
		} `json:"response"`
	} `json:"result"`
}

// TokenResponse represents the token and ratio data.
type TokenResponse struct {
	Token string `json:"token"`
	Ratio string `json:"ratio"`
}

// ResponseData represents the decoded response data.
type ResponseData struct {
	Stat struct {
		Height    int `json:"height"`
		Timestamp int `json:"timestamp"`
	} `json:"stat"`
	Response []TokenResponse `json:"response"`
}

// getDecodedResponseFromRPC retrieves the decoded response from the RPC API.
func getDecodedResponseFromRPC() (string, error) {
	baseURL := "https://dev.rpc.gnoswap.io/abci_query"
	query := url.QueryEscape(`"vm/qeval"`)
	data := url.QueryEscape(`"gno.land/r/demo/router\nApiGetRatiosFromBase()"`)
	url := fmt.Sprintf("%s?path=%s&data=%s", baseURL, query, data)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get response: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(result.Result.Response.ResponseBase.Data)
	if err != nil {
		log.Fatalf("Base64 decode failed: %v", err)
	}

	return string(decoded), nil
}

// extractJSON extracts the JSON string from the given log string.
func extractJSON(log string) string {
	start := strings.Index(log, "{")
	end := strings.LastIndex(log, "}") + 1

	jsonStr := log[start:end]
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

	return jsonStr
}

// TokenData represents the token data structure.
type TokenData struct {
	Status struct {
		Height    int `json:"height"`
		Timestamp int `json:"timestamp"`
	} `json:"stat"`
	TokenRatio []struct {
		TokenName string `json:"token"`
		Ratio     string `json:"ratio"`
	} `json:"response"`
}

// unmarshalResponseData unmarshals the JSON string into a ResponseData struct.
func unmarshalResponseData(jsonStr string) (*ResponseData, error) {
	var data ResponseData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return &data, nil
}
