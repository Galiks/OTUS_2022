package main

import (
	"log"

	cache "github.com/Galiks/OTUS_2022/hw04_lru_cache/cache"
)

func main() {
	c, err := cache.NewCache(5)
	if err != nil {
		log.Fatal(err)
	}
	c.Set("aaa", 100)

	c.Get("aaa")
}
