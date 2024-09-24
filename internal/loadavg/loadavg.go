package loadavg

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrUnsupported = errors.New("loadavg: unsupported platform")

type LoadAvg struct {
	mu    sync.RWMutex
	min1  float64
	min5  float64
	min15 float64
}

func New() *LoadAvg {
	return &LoadAvg{}
}

func (l *LoadAvg) WaitBelow(ctx context.Context, want float64, interval time.Duration) error {
	for {
		if err := l.Update(); err != nil {
			return err
		}

		if l.Get(Min1) < want {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
	return nil
}

type GetParam uint8

const (
	Min1 GetParam = iota
	Min5
	Min15
)

func (l *LoadAvg) Get(p GetParam) float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	switch p {
	case Min1:
		return l.min1
	case Min5:
		return l.min5
	case Min15:
		return l.min15
	default:
		panic(fmt.Sprintf("unknown LoadAvg.Get param: %d", p))
	}
}
