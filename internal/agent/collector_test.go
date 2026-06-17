package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	gauge map[string]float64
	count map[string]int64
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		gauge: make(map[string]float64),
		count: make(map[string]int64),
	}
}

func (m *mockWriter) UpdateGauge(name string, value float64) {
	m.gauge[name] = value
}
func (m *mockWriter) UpdateCounter(name string, value int64) {
	m.count[name] += value
}

func TestCollector(t *testing.T) {
	mock := newMockWriter()
	collector := NewCollector(mock)

	collector.Poll()
	collector.Poll()

	assert.Equal(t, int64(2), mock.count["PollCount"])
	assert.Equal(t, 28, len(mock.gauge))
	assert.Contains(t, mock.gauge, "Alloc")
	assert.Contains(t, mock.gauge, "RandomValue")
}
