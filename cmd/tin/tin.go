package main

import (
	"log"

	"github.com/Supercaly/tinlang/pkg/tin"
)

func main() {
	if err := tin.CompileFile("test/123.tin"); err != nil {
		log.Fatal(err)
	}
}
