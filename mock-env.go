package main

import (
	"fmt"
	"os"
)

func WithEnv(env map[string]string, work func()) {

	originalEnv := make(map[string]string)

	for key, val := range env {
		originalEnv[key] = os.Getenv(key)
		os.Setenv(key, val)
	}

	defer func(env map[string]string) {
		for key, val := range env {
			os.Setenv(key, val)
		}
	}(originalEnv)
	work()
}

func main() {
	os.Setenv("foo", "123")
	WithEnv(map[string]string{
		"foo": "999",
	}, func() {
		fmt.Println("Inside the worker function")
		fmt.Println(os.Getenv("foo"))
	})

	fmt.Println("Back in main")
	fmt.Println(os.Getenv("foo"))
}
