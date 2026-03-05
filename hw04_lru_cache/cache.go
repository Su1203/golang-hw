package hw04lrucache

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if node, ok := c.items[key]; ok {
		node.Value.(*cacheItem).value = value
		c.queue.MoveToFront(node)
		return true
	}

	// проверим переполнение
	if c.queue.Len() == c.capacity {
		lastNode := c.queue.Back()
		if lastNode != nil {
			delete(c.items, lastNode.Value.(*cacheItem).key)
			c.queue.Remove(lastNode)
		}
	}

	item := &cacheItem{key: key, value: value}
	node := c.queue.PushFront(item)
	c.items[key] = node

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if node, ok := c.items[key]; ok {
		c.queue.MoveToFront(node)
		return node.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
