package pool

import (
	"context"
	"sync"
	"time"
)

type Connection interface {
	Close() error
	IsAlive() bool
	ID() string
}

type Pool struct {
	connections chan Connection
	factory     func() (Connection, error)
	maxSize     int
	minSize     int
	mu          sync.RWMutex
	closed      bool
}

type PoolConfig struct {
	MaxSize     int
	MinSize     int
	MaxIdleTime time.Duration
	MaxLifetime time.Duration
}

func NewPool(factory func() (Connection, error), config PoolConfig) *Pool {
	pool := &Pool{
		connections: make(chan Connection, config.MaxSize),
		factory:     factory,
		maxSize:     config.MaxSize,
		minSize:     config.MinSize,
	}
	
	// Initialize minimum connections
	for i := 0; i < config.MinSize; i++ {
		conn, err := factory()
		if err == nil {
			pool.connections <- conn
		}
	}
	
	return pool
}

func (p *Pool) Get(ctx context.Context) (Connection, error) {
	select {
	case conn := <-p.connections:
		if conn.IsAlive() {
			return conn, nil
		}
		// Connection is dead, create new one
		return p.factory()
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// No available connections, create new one
		return p.factory()
	}
}

func (p *Pool) Put(conn Connection) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		conn.Close()
		return
	}
	p.mu.RUnlock()
	
	select {
	case p.connections <- conn:
	default:
		// Pool is full, close the connection
		conn.Close()
	}
}

func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.closed {
		return nil
	}
	
	p.closed = true
	close(p.connections)
	
	// Close all connections
	for conn := range p.connections {
		conn.Close()
	}
	
	return nil
}

func (p *Pool) Size() int {
	return len(p.connections)
}

func (p *Pool) IsClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.closed
}
