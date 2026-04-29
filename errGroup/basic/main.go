package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type ServiceHealth struct {
	Name string
	URL  string
}

func main() {
	services := []ServiceHealth{
		{"Google", "https://www.google.com"},
		{"BrokenAPI", "https://this.url.does.not.exist.reo"},
		{"Naver", "https://www.naver.com"},
		{"GitHubAPI", "https://api.github.com"},
	}

	g, ctx := errgroup.WithContext(context.Background())
	for _, svc := range services {
		svc := svc

		g.Go(func() error {
			reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(reqCtx, "GET", svc.URL, nil)
			if err != nil {
				return err
			}

			req.Header.Set("User-Agent", "Go-Backend-Handbook-Example")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("%s 접속 실패: %w", svc.Name, err)
			}
			defer resp.Body.Close()

			fmt.Printf("[%s] 상태: %s\n", svc.Name, resp.Status)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("system check error : %v\n", err)
		return
	}
	fmt.Println("all service is good")
}
