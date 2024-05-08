package vwap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDecodedResponseFromRPC(t *testing.T) {
	decodedResponse, err := getDecodedResponseFromRPC()
	if err != nil {
		t.Fatalf("Failed to get decoded response from RPC API: %v", err)
	}
	assert.Contains(t, decodedResponse, "stat")
	assert.Contains(t, decodedResponse, "height")
	assert.Contains(t, decodedResponse, "timestamp")
	assert.Contains(t, decodedResponse, "response")
	assert.Contains(t, decodedResponse, "token")
	assert.Contains(t, decodedResponse, "ratio")
}

func TestExtractJSON(t *testing.T) {
	log := `I[2023-05-08|18:12:00.123] Received response                           response="{\"stat\":{\"height\":123,\"timestamp\":1683540720},\"response\":[{\"token\":\"ABC\",\"ratio\":\"1.23\"},{\"token\":\"XYZ\",\"ratio\":\"4.56\"}]}"
`
	expected := `{"stat":{"height":123,"timestamp":1683540720},"response":[{"token":"ABC","ratio":"1.23"},{"token":"XYZ","ratio":"4.56"}]}`

	jsonStr := extractJSON(log)
	assert.Equal(t, expected, jsonStr)
}
