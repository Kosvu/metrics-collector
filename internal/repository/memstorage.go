package repository

type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMap:   make(map[string]float64),
		counterMap: make(map[string]int64),
	}
}
