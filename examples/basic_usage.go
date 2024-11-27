package main

import (
	"fmt"

	"github.com/nurmuh-alhakim18/gocache/cache"
)

func main() {
	c := cache.NewCache(10)
	c.Set("key1", "value1", 0)
	value, found := c.Get("key1")
	if found {
		fmt.Println(value) // Outputs "value1"
	}
}
