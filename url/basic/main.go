package main

import (
	"fmt"
	"net/url"
)

func main() {
	u := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/search",
	}

	q := url.Values{}
	q.Set("q", "go handbook")
	q.Set("page", "1")
	u.RawQuery = q.Encode()

	fmt.Println("built url:", u.String())

	parsed, _ := url.Parse(u.String())
	fmt.Println("parsed host:", parsed.Host)
	fmt.Println("query q:", parsed.Query().Get("q"))
}
