package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/utils"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	defaultRetry     = 30
	defautWaitTime   = 30 * time.Second
	defaultIndleTime = 10 * time.Second
	defaultPoolSize  = 5
)

type PoolData struct {
	IndleTimeout time.Duration
	WaitTimeout  time.Duration
	Size         int
	Max          int
	FactoryFn    func() (interface{}, error)        // Handle data before add to pool
	CloseFn      func(interface{}) error            // Close data before return back to pool
	PingFn       func(interface{}) chan interface{} // Control signal data in pool
	ForceClose   bool
	ping         chan interface{}
}

type Pool struct {
	p          *PoolData
	mu         *sync.Mutex
	conn       chan interface{}
	availConn  int
	forceClose bool
}

func (p *Pool) FillPool(data PoolData) error {
	if data.CloseFn == nil {
		return fmt.Errorf("close function is undefined")
	}
	if data.FactoryFn == nil {
		return fmt.Errorf("get function is undefined")
	}
	if data.WaitTimeout == 0 {
		data.WaitTimeout = defautWaitTime
	}
	if data.IndleTimeout == 0 {
		data.IndleTimeout = defaultIndleTime
	}
	if data.Max == 0 {
		data.Max = int(math.MaxInt)
	}
	if data.Size == 0 {
		data.Size = defaultPoolSize
	}
	p.forceClose = data.ForceClose
	p.conn = make(chan interface{}, data.Size)
	p.mu = &sync.Mutex{}
	p.p = &data
	start := time.Now()
	for i := 0; i < data.Size; i++ {
		currentTime := time.Since(start).Seconds()
		if currentTime >= data.WaitTimeout.Seconds() {
			return fmt.Errorf("timeout after: %v seconds", currentTime)
		}
		val, err := data.FactoryFn()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			continue
		}

		if data.PingFn != nil {
			data.ping = data.PingFn(val)
		}
		p.conn <- val
		p.availConn++
	}
	return nil
}

func (p *Pool) Get(retry ...int) (interface{}, error) {
	retryTimes := 0
	if retry == nil {
		retryTimes = defaultRetry
	}
	if p.conn == nil {
		return nil, fmt.Errorf("connection is not initialize")
	}

	for i := 0; i < retryTimes; i += 1 {
		select {
		case conn, ok := <-p.conn:
			p.mu.Lock()
			if !ok {
				return nil, fmt.Errorf("channel is closed")
			}
			defer p.mu.Unlock()
			p.availConn--
			return conn, nil
		case time := <-time.After(p.p.WaitTimeout):
			timeoutErr := fmt.Errorf("timeout after: %v seconds", time.Second())
			return nil, timeoutErr
		default:
			// If there are no data in pool and haven't reach max
			// Then create new connection
			if p.availConn < p.p.Max {
				res, err := p.p.FactoryFn()
				if err == nil {
					p.availConn++
				}
				return res, err
			}
			time.Sleep(p.p.IndleTimeout)
		}
	}
	return nil, fmt.Errorf("exceed max connection pool")
}

func (p *Pool) Release(data interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.forceClose {
		err := p.p.CloseFn(data)
		if err == nil {
			p.availConn++
			p.conn <- data
		}
		return err
	} else {
		p.availConn++
		p.conn <- data
	}
	return nil
}

func (p *Pool) Close() {
Loop:
	for {
		select {
		case c, ok := <-p.conn:
			if !ok { //ch is closed
				break Loop
			}
			err := p.p.CloseFn(c)
			if err != nil {
				utils.Log(LogConstant.Error, err)
			}
			close(p.p.ping)
		default: //all other case not-ready: means nothing in ch for now
			break Loop
		}
	}
	close(p.conn)
}
