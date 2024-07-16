package geecache

import (
	"fmt"
	"sync"
)

/**
Group 真正与用户交互：在cache的基础上又新增了没命中的回掉函数功能。
	一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。比如可以创建三个 Group，缓存学生的成绩命名为 scores，缓存学生信息的命名为 info，缓存学生课程的命名为 courses。
	第二个属性是 getter Getter，即缓存未命中时获取源数据的回调(callback)。
	第三个属性是 mainCache cache，即一开始实现的并发缓存。
	构建函数 NewGroup 用来实例化 Group，并且将 group 存储在全局变量 groups 中。
	GetGroup 用来特定名称的 Group，这里使用了只读锁 RLock()，因为不涉及任何冲突变量的写操作。
*/

type Group struct {
	name   string
	cache  Cache
	getter Getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	mu.Lock()
	defer mu.Unlock()
	if group, ok := groups[name]; ok {
		return group
	}
	group := &Group{
		name:   name,
		cache:  Cache{cacheBytes: cacheBytes},
		getter: getter,
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if group, ok := groups[name]; ok {
		return group
	}
	return nil
}

func (g *Group) Get(key string) (ByteView, error) {
	if len(key) <= 0 {
		return ByteView{}, fmt.Errorf("key is empty")
	}
	value, ok := g.cache.get(key)
	if ok {
		return value, nil
	}
	return g.load(key)

}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{data: bytes}
	g.cache.add(key, value)
	return value, nil
}
