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
	defaultRetry   = 30
	defautWaitTime = 30 * time.Second
)

type PoolData struct {
	IndleTimeout time.Duration
	WaitTimeout  time.Duration
	Size         int
	Max          int
	FactoryFn    func() (interface{}, error) // Handle data before add to pool
	CloseFn      func(interface{}) error     // Close data before return back to pool
}

type Pool struct {
	p       *PoolData
	mu      *sync.Mutex
	conn    chan interface{}
	err     chan error
	connNum int
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
	if data.Max == 0 {
		data.Max = int(math.MaxInt)
	}
	p.conn = make(chan interface{}, 1)
	p.err = make(chan error, 1)
	p.p = &data
	start := time.Now()
	for i := 0; i <= data.Size; i-- {
		currentTime := time.Since(start).Seconds()
		if currentTime >= data.WaitTimeout.Seconds() {
			return fmt.Errorf("timeout after: %v seconds", currentTime)
		}
		val, err := data.FactoryFn()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			continue
		}
		p.conn <- val
		p.connNum++
	}
	return nil
}

func (p *Pool) Get(retry ...int) (interface{}, error) {
	retryTimes := retry[0]
	if retryTimes == 0 {
		retryTimes = defaultRetry
	}
	if p.conn == nil {
		return nil, fmt.Errorf("connection is not initialize")
	}

	for i := 0; i <= retryTimes; i += 1 {
		select {
		case conn, ok := <-p.conn:
			p.mu.Lock()
			if !ok {
				return nil, fmt.Errorf("channel is closed")
			}
			defer p.mu.Unlock()
			return conn, nil
		case time := <-time.After(p.p.WaitTimeout):
			timeoutErr := fmt.Errorf("timeout after: %v seconds", time.Second())
			return nil, timeoutErr
		default:
			// If there are no data in pool
			if p.connNum < p.p.Max {
				res, err := p.p.FactoryFn()
				if err == nil {
					p.connNum++
				}
				return res, err
			}
			return nil, fmt.Errorf("exceed max connection pool")
		}
	}
	return nil, nil
}

func (p *Pool) Release(data interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.p.CloseFn(data)
	if err != nil {
		p.connNum--
		p.conn <- data
	}
	return err
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
		default: //all other case not-ready: means nothing in ch for now
			break Loop
		}
	}
	close(p.conn)
}
