package utils

import "fmt"

// Wraps the error with fmt.Errorf if err is not nil
func WrapError(label string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", label, err)
	} else {
		return nil
	}
}
