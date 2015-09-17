package main

import (
	"fmt"

	"github.com/pocke/gha"
)

func main() {
	r := &gha.Request{
		Note:   "foo",
		Scopes: []string{"gist"},
	}

	key, err := gha.CLI("gha-test-scope", "gha-test-scope-key", r)
	if err != nil {
		panic(err)
	}
	fmt.Println(key)
}
