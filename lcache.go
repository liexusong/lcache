// Simple local cache implements using LRU and expired

package lcache

import (
	"container/heap"
	"container/list"
	"sync"
	"time"
)

type Item struct {
	idx int           // Heap index
	ttl int64         // Time to lift
	ele *list.Element // LRU element pointer
	key string        // key
	val interface{}   // value
}

type Heap []*Item

type Cache struct {
	mutex    sync.Mutex       // Cache locker
	items    map[string]*Item // Items table
	expire   Heap             // Expire time heap
	lru      *list.List       // Item LRU list
	counter  int64            // Current object numbers
	MaxSize  int64            // Max object numbers
	StopChan chan struct{}
}

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	return h[i].ttl > h[j].ttl
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *Heap) Push(x interface{}) {
	last := len(*h)

	item := x.(*Item)
	item.idx = last

	*h = append(*h, item)
}

func (h *Heap) Pop() interface{} {
	temp := *h

	last := len(temp)

	item := temp[last-1]
	item.idx = -1

	*h = temp[0 : last-1]

	return item
}

// Create new cache object
// MaxSize: set the max object numbers of cache
// If cache's objects above the MaxSize
// GCItemsCycle() routine is recycling objects
func New(maxSize int64) *Cache {
	cache := &Cache{
		items:    make(map[string]*Item),
		expire:   make(Heap, 0),
		lru:      list.New(),
		MaxSize:  maxSize,
		StopChan: make(chan struct{}),
	}

	heap.Init(&cache.expire)

	go cache.GCItemsCycle()

	return cache
}

func (c *Cache) GCItemsCycle() {
	ticker := time.NewTicker(10 * time.Second)

	for {

		exitFlag := false

		select {
		case <-ticker.C:
			current := time.Now().Unix()

			c.mutex.Lock()

			for {
				size := len(c.expire)

				if size > 0 {
					item := c.expire[size-1]

					if item.ttl < current { // Item was expired?
						c.removeItem(item)
						continue
					}
				}

				if c.counter > c.MaxSize { // Object numbers above the MaxSize?
					kick := int64(float64(c.counter) * 0.3)

					for i := int64(0); i < kick; i++ {
						elem := c.lru.Front()

						item := elem.Value.(*Item)
						if item == nil { // Bug?
							panic("Item in LRU list but is a nil object")
						}

						c.removeItem(item)
					}
				}

				break
			}

			c.mutex.Unlock()

		case <-c.StopChan:
			exitFlag = true
		}

		if exitFlag {
			break
		}
	}
}

func (c *Cache) Set(key string, val interface{}, expire int64) {
	ttl := int64(0)

	if expire > 0 {
		ttl = time.Now().Unix() + expire
	}

	item := &Item{
		ttl: ttl,
		key: key,
		val: val,
		idx: -1,
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if tmp, exists := c.items[key]; exists {
		c.removeItem(tmp)
	}

	c.pushItem(item)
}

func (c *Cache) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exists := c.items[key]
	if !exists {
		return nil
	}

	if item.ttl < time.Now().Unix() { // Item expired?
		c.removeItem(item)
		return nil
	}

	// Move item to the back of LRU list
	c.lru.Remove(item.ele)

	item.ele = c.lru.PushBack(item)

	val := item.val

	return val
}

func (c *Cache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exists := c.items[key]; exists {
		c.removeItem(item)
		return true
	}

	return false
}

func (c *Cache) Size() int64 {
	c.mutex.Lock()
	size := c.counter
	c.mutex.Unlock()
	return size
}

func (c *Cache) Free() {
	c.StopChan <- struct{}{} // Stop GC cycle

	c.mutex.Lock()
	for _, item := range c.items {
		c.removeItem(item)
	}
	c.mutex.Unlock()
}

func (c *Cache) pushItem(item *Item) {
	c.items[item.key] = item

	if item.ttl > 0 {
		heap.Push(&c.expire, item)
	}

	item.ele = c.lru.PushBack(item)

	c.counter++
}

func (c *Cache) removeItem(item *Item) {
	delete(c.items, item.key)

	if item.idx >= 0 {
		heap.Remove(&c.expire, item.idx)
	}

	c.lru.Remove(item.ele)

	c.counter--
}
