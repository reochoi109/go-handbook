package main

import (
	"fmt"
	"net/url"
)

func sample() {
	v := url.Values{}

	// --- 쿼리 값 조작 ---
	v.Set("id", "reo109")   // 값 설정: 기존에 "id"가 있다면 덮어쓰고, 없으면 새로 만듭니다.
	v.Add("tag", "go")      // 값 추가: 기존 "tag"를 유지하며 새로운 값을 리스트에 추가합니다.
	v.Add("tag", "backend") // 결과: tag=go&tag=backend (하나의 키에 여러 값)
	v.Get("tag")            // 값 가져오기: 해당 키의 "첫 번째" 값만 반환합니다. (결과: "go")
	v.Has("mode")           // 존재 확인: 특정 키가 있는지 확인합니다. (Go 1.17+, 결과: false)
	v.Del("id")             // 값 삭제: 해당 키와 연결된 모든 값을 지웁니다.
	v.Encode()              // 인코딩: 전체 데이터를 "key=value&..." 형태의 문자열로 변환합니다.

	// --- URL 구조체 및 경로 조작 ---
	u, _ := url.Parse("https://example.com")
	u.JoinPath("api", "v1") // 경로 결합: 슬래시(/) 중복을 방지하며 경로를 안전하게 붙입니다. (Go 1.19+)
	u.String()              // 문자열 변환: 전체 URL 구조체를 하나의 문자열로 반환합니다.
	u.RequestURI()          // URI 추출: 호스트를 제외한 "/path?query" 부분만 가져옵니다.

	// --- 인코딩 유틸리티 (Escape) ---
	url.QueryEscape("a&b c") // 쿼리용 인코딩: 특수문자와 공백을 쿼리 규격에 맞게 변환합니다. (결과: a%26b+c)
	url.PathEscape("a&b c")  // 경로용 인코딩: URL 경로 규격에 맞게 변환합니다. (결과: a%26b%20c)

	// 확인용 출력
	fmt.Println("Final Query:", v.Encode())
	fmt.Println("Full URL:", u.String())
}
