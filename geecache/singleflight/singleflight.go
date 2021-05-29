package singleflight

import "sync"

type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 如果请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	c.val, c.err = fn()
	c.wg.Done()

	delete(g.m, key)

	return c.val, c.err
}