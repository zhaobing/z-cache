package singleflight

import "sync"

//call 正在进行中，或已经结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//Group 管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex //guard  m
	m  map[string]int
}
