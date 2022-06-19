package main

import (
	"fmt"
	"log"

	cache "github.com/Galiks/OTUS_2022/hw04_lru_cache/cache"
)

func main() {
	c, err := cache.NewCache(5)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		c.Set(cache.Key(fmt.Sprint(i)), i+1)
		fmt.Printf("\"1\": %v\n", i)
	}

	for i := 0; i < 6; i++ {
		if i == 0 {
			_, ok := c.Get(cache.Key(fmt.Sprint(i)))
			fmt.Printf("ok: %v\n", ok)
		} else {
			val, _ := c.Get(cache.Key(fmt.Sprint(i)))
			fmt.Printf("val: %v\n", val)
		}
	}
}
