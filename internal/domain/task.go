package domain

import (
	"time"
)

type Task struct {
	RequestsPerSec int
	PackPerSec     int
	Summary        *TaskSummary
}
type EMA interface {
	Add(float64)
	AddAndReturn(float64) float64
	Value() float64
}

type DurationSummary struct {
	Max       time.Duration
	Min       time.Duration
	EMA       EMA
	Histogram map[float64]int
}

type TaskSummary struct {
	Total    int
	Duration *DurationSummary
}
