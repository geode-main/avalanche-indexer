package client

import "sync/atomic"

type ClientPool struct {
	clients []*Client
	num     uint64
	idx     uint64
}

func NewClientPool(endpoints ...string) ClientPool {
	pool := ClientPool{
		clients: make([]*Client, len(endpoints)),
		num:     uint64(len(endpoints)),
		idx:     0,
	}

	for idx, endpoint := range endpoints {
		pool.clients[idx] = New(endpoint)
	}

	return pool
}

func (pool *ClientPool) Get() *Client {
	if pool.num == 0 {
		return nil
	}

	atomic.AddUint64(&pool.idx, 1)
	return pool.clients[pool.idx%pool.num]
}
