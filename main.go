package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/abdul-rehman-d/orders-api/application"
)

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start server %v\n", err)
	}
}
