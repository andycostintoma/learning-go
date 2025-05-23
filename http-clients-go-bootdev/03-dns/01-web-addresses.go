package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	DNSResponse struct {
		Status   int        `json:"Status"`
		Tc       bool       `json:"TC"`
		Rd       bool       `json:"RD"`
		Ra       bool       `json:"RA"`
		Ad       bool       `json:"AD"`
		Cd       bool       `json:"CD"`
		Question []Question `json:"Question"`
		Answer   []Answer   `json:"Answer"`
	}
	Question struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	}
	Answer struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	}
)

func getIPAddress(domain string) (string, error) {
	url := fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=%s&type=A", domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("accept", "application/dns-json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var dnsRes DNSResponse
	if err := json.Unmarshal(body, &dnsRes); err != nil {
		return "", fmt.Errorf("error unmarshalling json: %w", err)
	}

	if len(dnsRes.Answer) == 0 {
		return "", fmt.Errorf("no answer found")
	}

	return dnsRes.Answer[0].Data, nil
}

func main() {

	ip, err := getIPAddress("boot.dev")
	if err != nil {
		log.Fatalf("error getting IP address: %v", err)
	}
	fmt.Println(ip)
}
