package singleflight

import "sync"

//call 正在进行中，或已经结束的请求,使用sync.WaitGroup避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//Group 管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex //guard  m
	m  map[string]*call
}

//Do 针对函数fn,无论Do被调用多少次，fn都只执行一次，等待fn调用结束返回值或者错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok { //call is running
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	//call is not running
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	//execute call
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
