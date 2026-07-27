package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveCounter(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   int64
	}{
		{
			name:   "simple test",
			values: []int64{3, 5},
			want:   8,
		},
		{
			name:   "with negative number",
			values: []int64{10, -5},
			want:   5,
		},
		{
			name:   "with zero",
			values: []int64{5, 0},
			want:   5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemStorage()

			for _, v := range test.values {
				storage.SaveCounter(t.Context(), "PollCount", v)
			}

			assert.Equal(t, test.want, storage.counterMap["PollCount"])
		})
	}
}

func TestSaveGauge(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{
			name:   "simple test",
			values: []float64{10},
			want:   10,
		},
		{
			name:   "a few values",
			values: []float64{10, 15, 20},
			want:   20,
		},
		{
			name:   "with neagtive number",
			values: []float64{10, 15, -5},
			want:   -5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemStorage()

			for _, v := range test.values {
				storage.SaveGauge(t.Context(), "Alloc", v)
			}

			assert.Equal(t, test.want, storage.gaugeMap["Alloc"])
		})
	}
}
