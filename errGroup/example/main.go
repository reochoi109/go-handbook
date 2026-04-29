package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

type UserDashboard struct {
	Profile string
	Orders  string
	History string
}

func main() {
	ctx := context.Background()
	g, gCtx := errgroup.WithContext(ctx)

	var profile, orders, history string
	g.Go(func() error {
		var err error
		profile, err = fetchProfile(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		orders, err = fetchOrders(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		history, err = fetchHistory(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("failed to create dashboard: %v\n", err)
		return
	}
	dashboard := UserDashboard{
		Profile: profile,
		Orders:  orders,
		History: history,
	}
	fmt.Printf("result data : %+v\n", dashboard)
}

func fetchProfile(ctx context.Context) (string, error) {
	fmt.Println("fetchProfile: start")
	defer fmt.Println("fetchProfile: end")

	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return "", ctx.Err()
	}
	return "User: Name", nil
}

func fetchOrders(ctx context.Context) (string, error) {
	fmt.Println("fetchOrders: start")
	defer fmt.Println("fetchOrders: end")

	select {
	case <-time.After(150 * time.Millisecond):
	case <-ctx.Done():
		return "", ctx.Err()
	}

	// FAIL_ORDERS=1 이면 실패 → errgroup 컨텍스트 cancel 전파 확인
	if os.Getenv("FAIL_ORDERS") == "1" {
		return "", fmt.Errorf("orders service unavailable")
	}
	return "Orders: 3 items", nil
}

func fetchHistory(ctx context.Context) (string, error) {
	fmt.Println("fetchHistory: start")
	defer fmt.Println("fetchHistory: end")

	select {
	case <-time.After(250 * time.Millisecond):
	case <-ctx.Done():
		return "", ctx.Err()
	}
	return "History: Login at 10:00", nil
}
