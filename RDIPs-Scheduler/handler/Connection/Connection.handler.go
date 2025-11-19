package connection

import (
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	"RDIPs-Scheduler/utils"
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
	deadline := time.Now().Add(data.WaitTimeout)
	for i := 0; i < data.Size; i++ {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %.2fs while creating pool connections", data.WaitTimeout.Seconds())
		}

		val, err := data.FactoryFn()
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return fmt.Errorf("factory failed: %w", err)
		}

		if data.PingFn != nil {
			if err := data.PingFn(val); err != nil {
				utils.Log(LogConstant.Error, err)
				if data.CloseFn != nil {
					_ = data.CloseFn(val)
				}
				return fmt.Errorf("ping failed: %w", err)
			}
		}

		select {
		case p.conn <- PoolDetail{data: val, ctx: context.Background()}:
			p.mu.Lock()
			p.availConn++
			p.mu.Unlock()
		case <-time.After(data.WaitTimeout):
			// Avoid blocking indefinitely if channel is full
			if data.CloseFn != nil {
				_ = data.CloseFn(val)
			}
			return fmt.Errorf("timeout adding connection to pool")
		}
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

	for i := 0; i < retryTimes; i += 1 {
		select {
		case poolDetail, ok := <-p.conn:
			if !ok {
				return nil, nil, fmt.Errorf("channel is closed")
			}
			p.mu.Lock()
			p.availConn--
			p.mu.Unlock()
			if poolDetail.ctx.Err() != nil {
				// ctx is closed, skip this data
				continue
			}

			if poolDetail.timer != nil {
				if !poolDetail.timer.Stop() {
					<-poolDetail.timer.C
				}
			}
			if p.p.PingFn != nil {
				err := p.p.PingFn(poolDetail.data)
				if err != nil {
					// Connection is dead, close it
					utils.Log(LogConstant.Warning, "Connection validation failed:", err)
					p.p.CloseFn(poolDetail.data)

					// Create new connection
					newData, err := p.p.FactoryFn()
					if err != nil {
						utils.Log(LogConstant.Error, "Failed to create new connection:", err)
						continue
					}
					poolDetail.data = newData
					utils.Log(LogConstant.Info, "Created new connection to replace dead one")
				}
			}
			return poolDetail.data, poolDetail.ctx, nil
		case time := <-time.After(p.p.WaitTimeout):
			// in case of channel is nil
			timeoutErr := fmt.Errorf("timeout after: %v seconds", time.Second())
			return nil, nil, timeoutErr
		default:
			p.mu.Lock()
			currentSize := len(p.conn)
			canCreate := currentSize < p.p.Max
			p.mu.Unlock()
			// If there are no data in pool and haven't reach max
			// Then create new connection
			if canCreate {
				newCtx, cancel := context.WithCancel(context.Background())

				// Create timer to cancel this context after idle time
				timer := time.AfterFunc(p.p.IdleTimeout, cancel)
				res, err := p.p.FactoryFn()
				if err != nil {
					timer.Stop()
					return nil, nil, err
				}
				data := PoolDetail{ctx: newCtx, data: res, timer: timer}
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
func (p *Pool) Release(data interface{}, need_validate bool, ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if need_validate {
		err := p.p.PingFn(data)
		if err != nil {
			// Connection is dead, just close it and don't return to pool
			utils.Log(LogConstant.Warning, "Dead connection on release, closing:", err)
			p.p.CloseFn(data)
			return err
		}
	}

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

			if c.timer != nil {
				c.timer.Stop()
			}
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
