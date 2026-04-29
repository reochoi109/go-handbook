package main

import (
	"context"
	"fmt"
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
	time.Sleep(100 * time.Millisecond)
	return "User: Name", nil
}

func fetchOrders(ctx context.Context) (string, error) {
	time.Sleep(150 * time.Millisecond)
	return "Orders: 3 items", nil
}

func fetchHistory(ctx context.Context) (string, error) {
	return "History: Login at 10:00", nil
}
