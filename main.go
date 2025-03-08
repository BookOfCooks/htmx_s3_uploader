package main

import (
	"context"
	"fmt"
	"log"
	"main/app"
	"main/utils"
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

	// Start server
	serverPromise := utils.NewPromise(func() error {
		return utils.WrapError("serve", serve(ctx))
	})

	// Wait for server to close
	if err := serverPromise.Wait(); err != nil {
		return fmt.Errorf("serverPromise.Wait: %w", err)
	}

	return nil
}
