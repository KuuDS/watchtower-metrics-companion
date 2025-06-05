package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	apiAddress string
	apiToken   string
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Construct target URL
	targetURL := apiAddress
	if !strings.HasSuffix(targetURL, "/v1/metrics") {
		if strings.HasSuffix(targetURL, "/") {
			targetURL = targetURL + "v1/metrics"
		} else {
			targetURL = targetURL + "/v1/metrics"
		}
	}

	// Create request
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %s", err), http.StatusInternalServerError)
		return
	}

	// Add Authorization header
	req.Header.Set("Authorization", apiToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching metrics: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Forward response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, fmt.Sprintf("Error copying response body: %s", err), http.StatusInternalServerError)
	}
}

func main() {
	// Get API address and token from environment variables
	apiAddress = os.Getenv("WATCHTOWER_API_ADDRESS")
	if apiAddress == "" {
		fmt.Println("Error: WATCHTOWER_API_ADDRESS not set")
		os.Exit(1)
	}

	apiToken = os.Getenv("WATCHTOWER_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Error: WATCHTOWER_API_TOKEN not set")
		os.Exit(1)
	}

	http.HandleFunc("/v1/metrics", metricsHandler)
	http.HandleFunc("/health", healthHandler)

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	fmt.Printf("Using API address: %s\n", apiAddress)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
