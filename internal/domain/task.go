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

type TaskRequest struct {
	GwgrId      int64  `json:"gwgr_id"`
	OrigAnumber uint64 `json:"original_anumber"`
	OrigBnumber uint64 `json:"original_bnumber"`
	Anumber     uint64 `json:"anumber"`
	Bnumber     uint64 `json:"bnumber"`
}
