package main

import (
	"log"
	"os"
	"testing"
)

func TestRunMain(t *testing.T) {
	if os.Getenv("FUNCTIONAL_TESTS") != "1" {
		t.Skip("skipping test; FUNCTIONAL_TESTS is not enabled")
	}
	log.Println("run example with coverage enabled")
	main()
}
