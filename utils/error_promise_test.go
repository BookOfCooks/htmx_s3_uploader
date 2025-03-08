package utils

import (
	"errors"
	"testing"
)

func TestErrorPromiseBasic(t *testing.T) {
	p := NewErrorPromise(func() (string, error) {
		return "Hello", errors.ErrUnsupported
	})

	if v, err := p.Wait(); err != errors.ErrUnsupported || v != "Hello" {
		t.Errorf(`p.Wait() = %v, %v, want "Hello", %v`, v, err, errors.ErrUnsupported)
	}
}

func TestErrorPromisePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	p := NewErrorPromise(func() (string, error) {
		panic("Hello")
	})

	if _, err := p.Wait(); err == nil {
		t.Errorf(`p.Wait() = %v, want error`, err)
	}
}
