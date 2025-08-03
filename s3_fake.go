package main

import (
	"fmt"
	"io"
	"main/app"
	"os"
)

func FakeS3Upload(filename string, contents io.Reader) (string, error) {
	path := "public/audios/" + filename
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("os.Create: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, contents); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	return fmt.Sprintf("http://%s:%d/%s", app.ENV.HOST, app.ENV.PORT, path), nil
}
