package lcache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	cache := New(10000)

	for i := 0; i < 100000; i++ {
		key := fmt.Sprintf("key_%d", i)
		cache.Set(key, key, 2)
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
