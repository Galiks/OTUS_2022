package main

import (
	"fmt"

	. "github.com/Galiks/OTUS_2022/hw04_lru_cache/list"
)

func main() {
	l := NewList()

	l.PushBack("Hello")
	l.PushBack(",")
	l.PushBack("World")
	l.PushBack(123)
	l.PushBack([]byte{77})

	fmt.Printf("l.Print(): %v\n", l.Print())
}
