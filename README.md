# Local cache for Golang

Implements by LRU algorithm

## Methods:

### 1. Create cache object

```go
// Create local cache object
// @objectMaxSize: is the max object numbers
// @gcRate: is the objects GC rate (seconds)
cache := lcache.New(objectMaxSize int64, gcRate int64)
```

### 2. Set a key and value peer

```go
// Set a key and value peer into cache
// @key: the key
// @value: the value
// @expire: the expire time by seconds(0 is away not expire)
cache.Set(key string, value interface{}, expire int64)
```

### 3. Get value by a key

```go
// Get value by a key
// @key: the key
cache.Get(key string) interface{}
```

### 4. Delete a key

```go
// Delete a key
// @key: the key
cache.Delete(key string) bool
```

### 5. Free the cache

```go
// Free the cache
cache.Free()
```

## Example：

```go
package main

import(
    "github.com/liexusong/lcache"
)

func main() {
    cache := lcache.New(100000, 5) // max 100000 objects
    
    for i := 0; i < 100000; i++ {
        key := fmt.Sprintf("key_%d", i)
		cache.Set(key, key, 5)
    }
    
    for i := 0; i < 100000; i++ {
		key := fmt.Sprintf("key_%d", i)
		val := cache.Get(key)
		fmt.Println(val)
	}

	fmt.Println(cache.Size())

	time.Sleep(20 * time.Second)

	for i := 0; i < 100000; i++ {
		key := fmt.Sprintf("key_%d", i)
		val := cache.Get(key)
		fmt.Println(val)
	}

	fmt.Println(cache.Size())
}
```

