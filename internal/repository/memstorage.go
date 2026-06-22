package repository

import "sync"

type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
	mu         sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMap:   make(map[string]float64),
		counterMap: make(map[string]int64),
	}
}
