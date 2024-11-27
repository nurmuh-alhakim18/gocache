package cache

type nodeList struct {
	key        string
	value      interface{}
	prev, next *nodeList
}

type LRU struct {
	cache       map[string]*nodeList
	capacity    int
	front, back *nodeList
}

func NewLRU(capacity int) *LRU {
	lru := &LRU{
		cache:    make(map[string]*nodeList),
		capacity: capacity,
		front:    &nodeList{},
		back:     &nodeList{},
	}

	lru.front.next, lru.back.prev = lru.back, lru.front

	return lru
}

func (l *LRU) remove(node *nodeList) {
	prev, next := node.prev, node.next
	prev.next, next.prev = next, prev
}

func (l *LRU) insert(node *nodeList) {
	// insert before tail
	prev, next := l.back.prev, l.back
	prev.next = node
	next.prev = node
	node.next = next
	node.prev = prev
}

func (l *LRU) moveToBack(node *nodeList) {
	l.remove(node)
	l.insert(node)
}

func (l *LRU) Get(key string) (interface{}, bool) {
	if node, ok := l.cache[key]; ok {
		l.moveToBack(node)
		return node.value, true
	}

	return nil, false
}

func (l *LRU) Set(key string, value interface{}) {
	if node, ok := l.cache[key]; ok {
		l.remove(node)
		delete(l.cache, key)
	}

	newNode := &nodeList{
		key:   key,
		value: value,
	}

	l.insert(newNode)
	l.cache[key] = newNode

	if len(l.cache) > l.capacity {
		lru := l.front.next
		l.remove(lru)
		delete(l.cache, lru.key)
	}
}

func (l *LRU) Delete(key string) {
	if node, ok := l.cache[key]; ok {
		l.remove(node)
		delete(l.cache, key)
	}
}
