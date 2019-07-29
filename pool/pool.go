package pool

import (
	"fmt"
	"sync"
)

var ErrMaxConn = fmt.Errorf("maximum connections reached")

type Conn interface {
	Close() error
}

// ConnPool manages the life cycle of connections
type ConnPool struct {
	sync.RWMutex

	// Dial is used to create a new connection when necessary.
	Dial         func() (Conn, error)
	TestOnBorrow bool // 测试连接健康

	// Ping is use to check the conn fetched from pool
	Ping func(Conn) error

	MaxActive int // 同一时刻最多使用连接数 max active
	MaxIdle   int // 池子最大保留连接 max idle
	active    int // 当前正在使用的连接数 active = idle + using
	idlelist  []Conn

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	wait bool

	// mu protects fields defined below.
	mu   sync.Mutex
	cond *sync.Cond
}

func NewPool(maxIdle int, maxActive int, dial func() (Conn, error)) *ConnPool {
	return &ConnPool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial:      dial,
	}
}

func (pool *ConnPool) SetWait(wait bool) {
	pool.wait = wait
}

func (this *ConnPool) Get() (conn Conn, err error) {
	conn = this.getFromIdle()
	if conn != nil {
		if this.TestOnBorrow {
			err = this.Ping(conn)
			if err != nil {
				conn.Close()
				conn = this.getFromIdle()
				err = nil
			}
		}
		return
	}

	if this.reachedMax() {
		return nil, ErrMaxConn
		// TODO
		// wait 等其它请求释放
	}

	conn, err = this.Dial()
	if err != nil {
		return
	}

	if this.TestOnBorrow {
		err = this.Ping(conn)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	this.increActive()
	return
}

/**
* Close conn in pool
* 关掉连接 或 放入池中
 */
func (this *ConnPool) Close(conn Conn) {
	if this.overMaxIdle() {
		this.decreActive()
		if conn != nil {
			conn.Close()
		}
	} else {
		this.Lock()
		defer this.Unlock()
		this.idlelist = append(this.idlelist, conn)
	}
}

func (this *ConnPool) ForceClose(conn Conn) {
	this.decreActive()
	if conn != nil {
		conn.Close()
	}
}

func (this *ConnPool) Destroy() {
	this.Lock()
	defer this.Unlock()

	for _, conn := range this.idlelist {
		if conn != nil {
			conn.Close()
		}
	}
}

/**
* 尝试从池子中拿连接
 */
func (this *ConnPool) getFromIdle() Conn {
	this.Lock()
	defer this.Unlock()

	if len(this.idlelist) == 0 {
		return nil
	}

	conn := this.idlelist[0]
	this.idlelist = this.idlelist[1:]

	return conn
}

func (this *ConnPool) increActive() {
	this.Lock()
	defer this.Unlock()
	this.active += 1
}

func (this *ConnPool) decreActive() {
	this.Lock()
	defer this.Unlock()
	this.active -= 1
}

func (this *ConnPool) reachedMax() bool {
	this.RLock()
	defer this.RUnlock()
	return this.active >= this.MaxActive
}

func (this *ConnPool) overMaxIdle() bool {
	this.RLock()
	defer this.RUnlock()
	return len(this.idlelist) >= this.MaxIdle
}
