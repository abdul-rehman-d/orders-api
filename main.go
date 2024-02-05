package main

import (
	"context"
	"fmt"

	"github.com/abdul-rehman-d/orders-api/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Printf("Failed to start server %v\n", err)
	}
}
