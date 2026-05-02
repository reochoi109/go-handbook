package main

import (
	"fmt"
	"net/url"
)

func sample() {
	// --- 1. url.Values: 쿼리 스트링 조작 (map[string][]string 기반) ---
	v := url.Values{}

	// Set: 값을 설정 (기존 값은 덮어씀)
	v.Set("id", "reo109")
	v.Set("mode", "dark")

	// Add: 동일한 키에 값을 추가 (배열 형태로 저장됨)
	v.Add("tag", "go")
	v.Add("tag", "backend")
	v.Add("tag", "handbook")

	// Get: 첫 번째 값만 가져옴
	fmt.Println("Tag Get:", v.Get("tag")) // 출력: go

	// Has: 키 존재 여부 확인 (Go 1.17+)
	fmt.Println("Has mode:", v.Has("mode")) // 출력: true

	// Del: 해당 키의 모든 값 삭제
	v.Del("mode")

	// Encode: 알파벳 순서로 정렬하여 인코딩된 문자열 반환
	fmt.Println("Encoded Query:", v.Encode())
	// 출력: id=reo109&tag=go&tag=backend&tag=handbook

	// --- 2. url.URL: 경로 및 URI 조작 ---
	u, _ := url.Parse("https://example.com/api/v1/users")

	// JoinPath: 경로를 안전하게 결합 (Go 1.19+, 슬래시 중복 해결)
	newUrl := u.JoinPath("profile", "settings")
	fmt.Println("Joined Path:", newUrl.String())
	// 출력: https://example.com/api/v1/users/profile/settings

	// RequestURI: 스키마/호스트를 제외한 경로+쿼리 반환 (HTTP 요청 시 유용)
	u.RawQuery = v.Encode()
	fmt.Println("RequestURI:", u.RequestURI())
	// 출력: /api/v1/users?id=reo109&tag=go&tag=backend&tag=handbook

	// --- 3. Escape 유틸리티: 독립적인 인코딩/디코딩 ---
	// QueryEscape: 쿼리 스트링용 인코딩 (공백을 +로 변환)
	text := "hello world & go"
	escaped := url.QueryEscape(text)
	fmt.Println("QueryEscape:", escaped)
	// 출력: hello+world+%26+go

	// PathEscape: URL 경로용 인코딩 (공백을 %20으로 변환)
	pathEscaped := url.PathEscape(text)
	fmt.Println("PathEscape:", pathEscaped)
	// 출력: hello%20world%20%26%20go
}
