package main

import (
	"fmt"
	"io"
	"main/app"
	"os"

	"github.com/google/uuid"
)

func FakeS3Upload(filename string, ext string, contents io.Reader) (string, error) {
	path := fmt.Sprintf("public/audios/%s-%s.%s", filename, uuid.NewString(), ext)
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("os.Create: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, contents); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	return fmt.Sprintf("http://%s:%d/%s", app.ENV.HOST, app.ENV.PORT, path), nil
}
