package main

import (
	"net/url"
)

func _sample() {
	v := url.Values{}
	v.Set("id", "reo109")   // 기존에 "id"가 있다면 덮어쓰고, 없으면 새로 만듭니다.
	v.Add("tag", "go")      // 기존 "tag"를 유지하며 새로운 값을 리스트에 추가합니다.
	v.Add("tag", "backend") // tag=go&tag=backend
	v.Get("tag")            // 해당 키의 "첫 번째" 값만 반환합니다.
	v.Has("mode")           // 특정 키가 있는지 확인합니다.
	v.Del("id")             // 해당 키와 연결된 모든 값을 지웁니다.
	v.Encode()              // 전체 데이터를 "key=value&..." 형태의 문자열로 변환합니다.

	u, _ := url.Parse("https://example.com")
	u.JoinPath("api", "v1") // 슬래시(/) 중복을 방지하며 경로를 안전하게 붙입니다.
	u.String()              // 전체 URL 구조체를 하나의 문자열로 반환합니다.
	u.RequestURI()          // 호스트를 제외한 "/path?query" 부분만 가져옵니다.

	url.QueryEscape("a&b c") // 특수문자와 공백을 쿼리 규격에 맞게 변환합니다. (결과: a%26b+c)
	url.PathEscape("a&b c")  // URL 경로 규격에 맞게 변환합니다. (결과: a%26b%20c)
}
