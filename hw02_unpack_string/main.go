package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Galiks/OTUS_2022/hw02_unpack_string/unpack"
)

func main() {
	var input string

	flag.StringVar(&input, "input", "", "input value")
	flag.Parse()
	if input == "" {
		err := fmt.Errorf("invalid input. Please, use flag `-input`")
		log.Fatal(err)
	}
	log.Printf("Your input: %v\n", input)
	output, err := unpack.Unpack(input)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Result: %v\n", output)
}
