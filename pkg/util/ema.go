package util

import "sync"

type EMA interface {
	Add(float64)
	AddAndReturn(float64) float64
	Value() float64
}

type ema struct {
	alpha  float64
	mu     sync.RWMutex
	value  float64
	warmed bool
}

func NewEMA(alpha float64) *ema {
	return &ema{alpha: alpha}
}

func (e *ema) Add(value float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.warmed {
		e.value = (value * e.alpha) + (e.value * (1 - e.alpha))
	} else {
		e.value = value
		e.warmed = true
	}
}

func (e *ema) AddAndReturn(value float64) float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.warmed {
		e.value = (value * e.alpha) + (e.value * (1 - e.alpha))
	} else {
		e.value = value
		e.warmed = true
	}
	return e.value
}

func (e *ema) Value() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.value
}
