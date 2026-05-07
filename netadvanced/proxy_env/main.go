package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	setupEnv()
	defer cleanupEnv()

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	targets := []string{
		"http://example.com",
		"http://localhost:8080",
		"http://openai.com",
		"http://google.com",
	}

	for _, target := range targets {
		printProxyDecision(transport, target)
	}
}

func setupEnv() {
	mustSetenv("HTTP_PROXY", "http://proxy.local:3128")
	mustSetenv("NO_PROXY", "example.com,localhost")
}

func cleanupEnv() {
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("NO_PROXY")
}

func mustSetenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
}

func printProxyDecision(transport *http.Transport, target string) {
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		fmt.Printf("%-25s -> request error: %v\n", target, err)
		return
	}

	proxyURL, err := transport.Proxy(req)
	if err != nil {
		fmt.Printf("%-25s -> proxy error: %v\n", target, err)
		return
	}

	fmt.Printf("%-25s -> %s\n", target, formatProxy(proxyURL))
}

func formatProxy(proxyURL *url.URL) string {
	if proxyURL == nil {
		return "DIRECT"
	}
	return "PROXY: " + proxyURL.String()
}
