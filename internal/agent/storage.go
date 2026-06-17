package agent

import "sync"

type AgentStorage struct {
	gaugeStorage   map[string]float64
	counterStorage map[string]int64
	mu             sync.RWMutex
}

func NewAgentStorage() *AgentStorage {
	return &AgentStorage{
		gaugeStorage:   make(map[string]float64),
		counterStorage: make(map[string]int64),
		mu:             sync.RWMutex{},
	}
}

func (s *AgentStorage) UpdateGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gaugeStorage[name] = value
}

func (s *AgentStorage) UpdateCounter(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counterStorage[name] += value
}

func (s *AgentStorage) GetAll() (map[string]float64, map[string]int64) {

	s.mu.RLock()
	defer s.mu.RUnlock()
	tmpGauge := make(map[string]float64)
	for key, value := range s.gaugeStorage {
		tmpGauge[key] = value
	}

	tmpCounter := make(map[string]int64)
	for key, value := range s.counterStorage {
		tmpCounter[key] = value
	}

	return tmpGauge, tmpCounter
}
