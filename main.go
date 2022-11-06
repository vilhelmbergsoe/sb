package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s, err := newServer()
	if err != nil {
		return err
	}

	fmt.Println("server running on :8080")

	if err := http.ListenAndServe(":8080", s.router); err != nil {
		return err
	}

	return nil
}
