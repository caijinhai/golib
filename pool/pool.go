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

/**
* idlelist 成员：池子中的连接
 */
type idle struct {
	p ConnPool
	c Conn
	t time.Time
}

// ConnPool manages the life cycle of connections
type ConnPool struct {
	sync.RWMutex

	// Dial is used to create a new connection when necessary.
	Dial         func() (Conn, error)
	TestOnBorrow func(Conn) error // 测试连接健康，比如检查角色是否master，ping是否正常

	MaxActive   int           // 同一时刻最多使用连接数 max active
	MaxIdle     int           // 池子最大保留连接 max idle
	IdleTimeout time.Duration // 池子中的连接过期时间

	// If Wait is true and the pool is at the MaxActive limit, then Get() Waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	active   int // 当前正在使用的连接数 active = idle + using
	idlelist list.List
	// mu protects fields defined below.
	mu   sync.Mutex
	cond *sync.Cond
}

func New(maxIdle int, maxActive int, idleTimeoutS int, dial func() (Conn, error), wait bool) *ConnPool {
	pool := &ConnPool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeoutS) * time.Second,
		Dial:        dial,
		Wait:        wait,
	}
	if pool.Wait {
		pool.cond = sync.NewCond(&pool.mu)
	}
	// pool.TestOnBorrow = func(c Conn) error {
	// 	_, err := c.Do("PING")
	// 	return err
	// }

	return pool
}

/**
* 设计理念：所有的对外接口，上层加锁，下层私有函数不加锁，防止同一个锁在不同层中使用导致重入，因为go没有可重入锁
 */
func (this *ConnPool) Get() (conn Conn, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.IdleTimeout > 0 {
		this.closeExipredIdle()
	}
	for {
		// 从连接池中取
		conn = this.getIdle()
		if conn != nil {
			if this.TestOnBorrow != nil {
				err = this.TestOnBorrow(conn)
				if err != nil {
					conn.Close()
					conn = this.getIdle()
					err = nil
				}
			}
			break
		}

		// 连接池不够，看是否超过活跃上限，超过则等待其它请求释放连接，否则dial一个新连接
		if this.MaxActive == 0 || !this.overMaxActive() {
			conn, err = this.Dial()
			if err == nil {
				this.increActive(1)
			}
			if this.TestOnBorrow != nil {
				err = this.TestOnBorrow(conn)
				if err != nil {
					conn.Close()
					conn = nil
				}
			}
			break
		}

		if !this.Wait {
			conn = nil
			err = ErrMaxConn
			break
		}
		this.cond.Wait()
	}
	return conn, err
}

/**
* 释放使用的连接
* 关掉连接 或 放入池中
 */
func (this *ConnPool) Release(conn Conn) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.overMaxIdle() {
		this.close(conn)
	} else {
		this.idlelist.PushFront(idle{t: time.Now(), c: conn})
	}
	if this.cond != nil {
		this.cond.Signal()
	}
}

/**
* 销毁关闭所有连接
 */
func (this *ConnPool) Destory() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.decreActive(this.len())
	idlelist := this.idlelist
	this.idlelist.Init()
	for e := idlelist.Front(); e != nil; e = e.Next() {
		e.Value.(idle).c.Close()
	}
	if this.cond != nil {
		this.cond.Broadcast()
	}
}

/**
* 尝试从池子中拿连接
 */
func (this *ConnPool) getIdle() Conn {
	if this.len() == 0 {
		return nil
	}

	e := this.idlelist.Front()
	ic := e.Value.(idle)
	this.idlelist.Remove(e)

	return ic.c
}

/**
* 获取连接池中的连接数
 */
func (this *ConnPool) len() int {
	return this.idlelist.Len()
}

/**
* 关闭连接
 */
func (this *ConnPool) close(conn Conn) {
	this.decreActive(1)
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
		ic := e.Value.(idle)
		if time.Now().Before(ic.t.Add(this.IdleTimeout)) {
			break
		}
		this.idlelist.Remove(e)
		this.close(ic.c)
	}
}

func (this *ConnPool) increActive(num int) {
	this.active += num
}

func (this *ConnPool) decreActive(num int) {
	this.active -= num
}

/**
* 超过活跃的连接数
 */
func (this *ConnPool) overMaxActive() bool {
	return this.active >= this.MaxActive
}

/**
* 超过池子中的连接数
 */
func (this *ConnPool) overMaxIdle() bool {
	return this.len() >= this.MaxIdle
}
