package main

import (
	"context"
	"fmt"
	"log"
	"main/app"
	"main/utils"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
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
		return utils.WrapError("startServer", startServer(ctx))
	})

	// Wait for server to close
	if err := serverPromise.Wait(); err != nil {
		return fmt.Errorf("serverPromise.Wait: %w", err)
	}

	return nil
}

func startServer(ctx context.Context) error {
	server := &http.Server{
		Addr:    app.ENV.HOST + ":" + strconv.Itoa(app.ENV.PORT),
		Handler: router(),
	}

	shutdownPromise := utils.NewPromise(func() error {
		return utils.WrapError("waitThenCloseServer", waitThenCloseServer(ctx, server))
	})

	log.Println("Starting server on http://" + server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}

	if err := shutdownPromise.Wait(); err != nil {
		return fmt.Errorf("shutdownPromise.Wait: %w", err)
	}
	return nil
}

func waitThenCloseServer(ctx context.Context, server *http.Server) error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	select {
	case <-ctx.Done():
	case <-sigint:
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server.Shutdown: %w", err)
	}
	return nil
}
