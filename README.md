# Simple Cache for Go

A lightweight, thread-safe cache implementation in Go using LRU method

## Features

- Implemented using LRU cache.
- Thread-safe for concurrent operations.
- Easy-to-use.

## Installation

To add to your project, use:

```
go get github.com/nurmuh-alhakim18/gocache
```

## Usage

Here’s how you can use the library:

```
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
```
