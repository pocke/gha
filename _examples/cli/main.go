package main

import (
	"fmt"

	"github.com/pocke/gha"
)

func main() {
	key, err := gha.CLI("test app for gha", "test-key")
	if err != nil {
		panic(err)
	}
	fmt.Println(key)
}
