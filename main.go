package main

import (
	"context"
	"fmt"
	"log"
	"main/app"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalln("run:", err)
	}
	log.Println("Application completed.")
}

func run(ctx context.Context) error {
	if err := app.Init(ctx); err != nil {
		return fmt.Errorf("app.Init: %w", err)
	}

	return nil
}
