package connection

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	defaultRetry    = 30
	defaultWaitTime = 30 * time.Second // time to get from channel
	defaultIdleTime = 5 * time.Minute  // time limit for one idle channel
	defaultPoolSize = 5
)

type PoolData struct {
	IdleTimeout time.Duration
	WaitTimeout time.Duration
	Size        int
	Max         int
	FactoryFn   func() (interface{}, error) // Handle data before add to pool
	CloseFn     func(interface{}) error     // Close data before return back to pool
	PingFn      func(interface{}) error     // Control signal data in pool
	ForceClose  bool
	ping        chan interface{}
}

type Pool struct {
	p          *PoolData
	mu         *sync.Mutex
	conn       chan PoolDetail
	availConn  int
	forceClose bool
}

type PoolDetail struct {
	data  interface{}
	ctx   context.Context
	timer *time.Timer // timer for idle data
}

func (p *Pool) FillPool(data PoolData) error {
	if data.CloseFn == nil {
		return fmt.Errorf("close function is undefined")
	}
	if data.FactoryFn == nil {
		return fmt.Errorf("get function is undefined")
	}
	if data.WaitTimeout == 0 {
		data.WaitTimeout = defaultWaitTime
	}
	if data.IdleTimeout == 0 {
		data.IdleTimeout = defaultIdleTime
	}
	if data.Max == 0 {
		data.Max = 100
	}
	if data.Size == 0 {
		data.Size = defaultPoolSize
	}
	p.forceClose = data.ForceClose
	p.conn = make(chan PoolDetail, data.Max)
	p.mu = &sync.Mutex{}
	p.p = &data
	start := time.Now()
	for i := 0; i < data.Size; i++ {
		currentTime := time.Since(start).Seconds()
		if currentTime >= data.WaitTimeout.Seconds() {
			close(p.conn)
			return fmt.Errorf("timeout after: %v seconds", currentTime)
		}
		val, err := data.FactoryFn()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return err
		}

		if data.PingFn != nil {
			err := data.PingFn(val)
			if err != nil {
				close(p.conn)
				return err
			}
		}
		p.conn <- PoolDetail{data: val, ctx: context.Background()}
		p.availConn++
	}
	return nil
}

func (p *Pool) Get(retry ...int) (interface{}, context.Context, error) {
	retryTimes := 0
	if retry == nil {
		retryTimes = defaultRetry
	}
	if p.conn == nil {
		return nil, nil, fmt.Errorf("connection is not initialize")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < retryTimes; i += 1 {
		select {
		case poolDetail, ok := <-p.conn:
			if !ok {
				return nil, nil, fmt.Errorf("channel is closed")
			}
			p.availConn--
			if poolDetail.ctx.Err() != nil {
				// ctx is closed, skip this data
				continue
			}
			if poolDetail.timer != nil {
				if !poolDetail.timer.Stop() {
					<-poolDetail.timer.C
				}
				poolDetail.timer.Reset(p.p.IdleTimeout)
			}
			return poolDetail.data, poolDetail.ctx, nil
		case time := <-time.After(p.p.WaitTimeout):
			// in case of channel is nil
			timeoutErr := fmt.Errorf("timeout after: %v seconds", time.Second())
			return nil, nil, timeoutErr
		default:
			// If there are no data in pool and haven't reach max
			// Then create new connection
			if len(p.conn) < p.p.Max {
				newCtx, cancel := context.WithCancel(context.Background())

				// Create timer to cancel this context after idle time
				timer := time.AfterFunc(p.p.IdleTimeout, cancel)
				res, err := p.p.FactoryFn()
				data := PoolDetail{ctx: newCtx, data: res, timer: timer}
				p.conn <- data
				// The newly created connection should not be opened forever
				// wait for idle time out before running close function
				go p.handleTimeoutCtx(data.data, newCtx)
				p.availConn++
				return res, newCtx, err
			}
			time.Sleep(1 * time.Second)
		}
	}
	return nil, nil, fmt.Errorf("exceed max connection pool")
}

func (p *Pool) handleTimeoutCtx(data interface{}, ctx context.Context) error {
	<-ctx.Done()
	utils.Log(LogConstant.Info, "Exceed Idle time:", p.p.IdleTimeout.Seconds(), ", begin close this resource: ")
	p.p.CloseFn(data)
	return ctx.Err()
}

// return back to Pool
func (p *Pool) Release(data interface{}, ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.forceClose {
		err := p.p.CloseFn(data)
		if err == nil {
			p.availConn++
			p.conn <- PoolDetail{data: data, ctx: ctx}
		}
		return err
	} else {
		p.availConn++
		p.conn <- PoolDetail{data: data, ctx: ctx}
	}
	return nil
}

func (p *Pool) Close() {
Loop:
	for {
		select {
		case c, ok := <-p.conn:
			if !ok { //ch is closed
				utils.Log(LogConstant.Info, "Channel is closed")
				break Loop
			}
			utils.Log(LogConstant.Info, "Close all child channel")
			err := p.p.CloseFn(c.data)
			if err != nil {
				utils.Log(LogConstant.Error, err)
			}
			if p.p.ping != nil {
				close(p.p.ping)
			}
		default: //all other case not-ready: means nothing in ch for now
			utils.Log(LogConstant.Info, "Nothing is in channel")
			break Loop
		}
	}
	close(p.conn)
}
