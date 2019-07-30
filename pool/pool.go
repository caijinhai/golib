package pool

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

var ErrMaxConn = fmt.Errorf("maximum connections reached")

type Conn interface {
	Close() error
}

type IdleConn struct {
	c Conn
	t time.Time
}

// ConnPool manages the life cycle of connections
type ConnPool struct {
	sync.RWMutex

	// Dial is used to create a new connection when necessary.
	Dial         func() (Conn, error)
	TestOnBorrow bool // 测试连接健康

	// Ping is use to check the conn fetched from pool
	Ping func(Conn) error

	MaxActive   int           // 同一时刻最多使用连接数 max active
	MaxIdle     int           // 池子最大保留连接 max idle
	IdleTimeout time.Duration // 池子中的连接过期时间
	active      int           // 当前正在使用的连接数 active = idle + using
	idlelist    list.List

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	wait bool

	// mu protects fields defined below.
	mu   sync.Mutex
	cond *sync.Cond
}

func New(maxIdle int, maxActive int, idleTimeout int, dial func() (Conn, error)) *ConnPool {
	return &ConnPool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial:        dial,
	}
}

func (pool *ConnPool) SetWait(wait bool) {
	pool.wait = wait
}

func (this *ConnPool) Get() (conn Conn, err error) {
	if this.IdleTimeout > 0 {
		this.closeExipredIdle()
	}
	conn = this.get()
	if conn != nil {
		if this.TestOnBorrow {
			err = this.Ping(conn)
			if err != nil {
				conn.Close()
				conn = this.get()
				err = nil
			}
		}
		return
	}

	if this.overMaxActive() {
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
* 释放使用的连接
* 关掉连接 或 放入池中
 */
func (this *ConnPool) Release(conn Conn) {
	if this.overMaxIdle() {
		this.close(conn)
	} else {
		this.Lock()
		defer this.Unlock()
		this.idlelist.PushFront(IdleConn{t: time.Now(), c: conn})
	}
}

/**
* 销毁关闭所有连接
 */
func (this *ConnPool) Destory() {
	this.Lock()
	defer this.Unlock()

	this.resetActive()
	for e := this.idlelist.Front(); e != nil; e = e.Next() {
		e.Value.(IdleConn).c.Close()
	}
}

/**
* 尝试从池子中拿连接
 */
func (this *ConnPool) get() Conn {
	this.Lock()
	defer this.Unlock()

	if this.idlelist.Len() == 0 {
		return nil
	}

	e := this.idlelist.Front()
	ic := e.Value.(IdleConn)
	this.idlelist.Remove(e)

	return ic.c
}

/**
* 关闭连接
 */
func (this *ConnPool) close(conn Conn) {
	this.decreActive()
	if conn != nil {
		conn.Close()
	}
}

/**
* 关闭过期的连接
 */
func (this *ConnPool) closeExipredIdle() {
	for {
		// 从最后往前
		e := this.idlelist.Back()
		if e == nil {
			break
		}
		ic := e.Value.(IdleConn)
		if time.Now().Before(ic.t.Add(this.IdleTimeout)) {
			break
		}
		this.idlelist.Remove(e)
		this.close(ic.c)
	}
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

func (this *ConnPool) resetActive() {
	this.Lock()
	defer this.Unlock()
	this.active = 0
}

/**
* 超过活跃的连接数
 */
func (this *ConnPool) overMaxActive() bool {
	this.RLock()
	defer this.RUnlock()
	return this.active >= this.MaxActive
}

/**
* 超过池子中的连接数
 */
func (this *ConnPool) overMaxIdle() bool {
	this.RLock()
	defer this.RUnlock()
	return this.idlelist.Len() >= this.MaxIdle
}
