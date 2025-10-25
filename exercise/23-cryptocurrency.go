package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type apiResponse struct {
	Status string                 `json:"status"`
	Stats  map[string]interface{} `json:"stats"`
}

func GetExchangeRate(source, destination string) (string, error) {
	source = strings.ToLower(strings.TrimSpace(source))
	destination = strings.ToLower(strings.TrimSpace(destination))

	if destination == "" {
		destination = "rls"
	}

	url := fmt.Sprintf("http://localhost:4001/rates?srcCurrency=%s&dstCurrency=%s", source, destination)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid request or server error")
	}

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != "OK" {
		return "", errors.New("response status not OK")
	}

	key := fmt.Sprintf("%s-%s", source, destination)
	stat, ok := result.Stats[key]
	if !ok {
		return "", errors.New("invalid response data")
	}

	m, ok := stat.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid stats format")
	}

	value, ok := m["latest"].(string)
	if !ok {
		return "", errors.New("missing latest field")
	}

	return value, nil
}
