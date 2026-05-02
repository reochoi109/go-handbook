package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

func main() {
	tests := []string{
		"https://example.com/path",
		"http://example.com/path",          // scheme 제한 예시
		"https://127.0.0.1/admin",          // loopback 차단
		"https://localhost/admin",          // loopback 차단
		"https://169.254.169.254/latest/",  // metadata IP(예시) 차단
		"https://evil.com@127.0.0.1/admin", // userinfo trick
	}

	for _, raw := range tests {
		err := validateOutboundURL(raw, []string{"example.com"})
		fmt.Println(raw, "->", err)
	}
}

func validateOutboundURL(raw string, allowHosts []string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}

	// scheme 제한
	if u.Scheme != "https" {
		return errors.New("scheme not allowed")
	}

	if u.User != nil {
		return errors.New("userinfo not allowed")
	}

	host := u.Hostname()
	if host == "" {
		return errors.New("missing host")
	}

	// allowlist
	allowed := false
	for _, h := range allowHosts {
		if strings.EqualFold(host, h) {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("host not allowed")
	}

	if ip := net.ParseIP(host); ip != nil {
		if isBadIP(ip) {
			return errors.New("ip not allowed")
		}
	}

	return nil
}

func isBadIP(ip net.IP) bool {
	// loopback, link-local, private 대역 등은 환경/정책에 따라 차단 대상
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	// RFC1918 private (IPv4)
	if v4 := ip.To4(); v4 != nil {
		switch {
		case v4[0] == 10:
			return true
		case v4[0] == 172 && v4[1] >= 16 && v4[1] <= 31:
			return true
		case v4[0] == 192 && v4[1] == 168:
			return true
		case v4[0] == 169 && v4[1] == 254: // link-local
			return true
		}
	}
	return false
}
