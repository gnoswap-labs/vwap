package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

type Response struct {
	Result struct {
		Response struct {
			ResponseBase struct {
				Data string `json:"data"`
			} `json:"responseBase"`
		} `json:"response"`
	} `json:"result"`
}

func GetDecodedStringFromRPC() string {
	baseURL := "https://dev.rpc.gnoswap.io/abci_query"
	query := url.QueryEscape(`"vm/qeval"`)
	data := url.QueryEscape(`"gno.land/r/demo/router\nApiGetRatiosFromBase()"`)
	url := fmt.Sprintf("%s?path=%s&data=%s", baseURL, query, data)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("JSON unmarshal failed: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(result.Result.Response.ResponseBase.Data)
	if err != nil {
		log.Fatalf("Base64 decode failed: %v", err)
	}

	return string(decoded)
}

type Data struct {
	Stat struct {
		Height    int `json:"height"`
		Timestamp int `json:"timestamp"`
	} `json:"stat"`
	Response []struct {
		Token string `json:"token"`
		Ratio string `json:"ratio"`
	} `json:"response"`
}

func main() {
	log := GetDecodedStringFromRPC()

	start := strings.Index(log, "{")
	end := strings.LastIndex(log, "}") + 1
	jsonStr := log[start:end]
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

	type TokenResponse struct {
		Token string `json:"token"`
		Ratio string `json:"ratio"`
	}
	type ResponseData struct {
		Stat struct {
			Height    int `json:"height"`
			Timestamp int `json:"timestamp"`
		} `json:"stat"`
		Response []TokenResponse `json:"response"`
	}

	var data ResponseData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println("Error Parsing JSON:", err)
		return
	}

	fmt.Printf("Height: %d, Timestamp: %d\n", data.Stat.Height, data.Stat.Timestamp)
	for _, r := range data.Response {
		ratio, ok := new(big.Int).SetString(r.Ratio, 10)
		if !ok {
			fmt.Println("Error converting ratio to big.Int")
			return
		}

		divisor := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
		result := new(big.Float).Quo(new(big.Float).SetInt(ratio), new(big.Float).SetInt(divisor))

		fmt.Printf("Token: %s, Ratio: %s\n", r.Token, result.String())
	}
}
