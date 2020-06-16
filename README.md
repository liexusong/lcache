# Local cache for Golang

### Usage：

```go
package main

import(
    "github.com/liexusong/lcache"
)

func main() {
    cache := lcache.New(100000) // max 100000 objects
    
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

