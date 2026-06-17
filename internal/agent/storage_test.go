package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   int64
	}{
		{
			name:   "simple test",
			values: []int64{35, 35},
			want:   70,
		},
		{
			name:   "zero value",
			values: []int64{35, 0},
			want:   35,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewAgentStorage()

			for _, val := range test.values {
				storage.UpdateCounter("PollCount", val)
			}

			assert.Equal(t, test.want, storage.counterStorage["PollCount"])
		})
	}
}

func TestGetAll(t *testing.T) {
	storage := NewAgentStorage()

	storage.UpdateGauge("Alloc", 1)
	gauge, _ := storage.GetAll()
	gauge["Alloc"] = 5

	assert.Equal(t, float64(1), storage.gaugeStorage["Alloc"])
}
