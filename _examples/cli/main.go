package main

import (
	"fmt"

	"github.com/pocke/gha"
)

func main() {
	key, err := gha.CLI("test-key", &gha.Request{
		Note: "test app for gha",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(key)
}
