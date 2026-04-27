package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// .env 파일을 로드합니다.
	err := godotenv.Load()
	if err != nil {
		log.Println("경고: .env 파일을 찾을 수 없습니다. OS 환경변수를 사용합니다.")
	}

	// 환경변수를 불러옵니다.
	dbHost := os.Getenv("DB_HOST")
	apiKey := os.Getenv("API_KEY")

	fmt.Printf("DB Host: %s\n", dbHost)
	fmt.Printf("API Key: %s\n", apiKey)
}
