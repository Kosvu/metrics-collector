package agent

import (
	"compress/gzip"
	"encoding/json"
	models "metrics/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeReader struct{}

func (f *fakeReader) GetAll() (map[string]float64, map[string]int64) {
	return map[string]float64{"Alloc": 1}, map[string]int64{"PollCount": 5}
}

func TestSend(t *testing.T) {
	got := make(map[string]models.Metrics)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var metrics []models.Metrics
		gz, err := gzip.NewReader(r.Body)
		require.NoError(t, err)
		defer gz.Close()
		err = json.NewDecoder(gz).Decode(&metrics)
		require.NoError(t, err)
		for _, m := range metrics {
			got[m.ID] = m
		}
	}))
	defer srv.Close()
	addr := strings.TrimPrefix(srv.URL, "http://")
	sender := NewSender(&fakeReader{}, addr, srv.Client())
	sender.Send()
	require.NotNil(t, got["Alloc"].Value)
	assert.Equal(t, float64(1), *got["Alloc"].Value)
	require.NotNil(t, got["PollCount"].Delta)
	assert.Equal(t, int64(5), *got["PollCount"].Delta)
}
