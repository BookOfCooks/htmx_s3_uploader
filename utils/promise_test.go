package utils

import (
	"testing"
)

func TestPromiseBasic(t *testing.T) {
	p := NewPromise(func() string {
		return "Hello"
	})

	if str := p.Wait(); str != "Hello" {
		t.Errorf(`p.Wait() = %q, want "Hello"`, str)
	}
}

func TestPromisePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	p := NewPromise(func() string {
		panic("Hello")
	})

	if str := p.Wait(); str != "" {
		t.Errorf(`p.Wait() = %q, want ""`, str)
	}
}
