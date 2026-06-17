package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeReader struct{}

func (f *fakeReader) GetAll() (map[string]float64, map[string]int64) {
	return map[string]float64{"Alloc": 1}, map[string]int64{"PollCount": 5}
}

func TestSend(t *testing.T) {
	var gotPaths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
	}))
	defer srv.Close()
	addr := strings.TrimPrefix(srv.URL, "http://")
	sender := NewSender(&fakeReader{}, addr, srv.Client())
	sender.Send()
	assert.Equal(t, 2, len(gotPaths))
	assert.Contains(t, gotPaths, "/update/gauge/Alloc/1")
	assert.Contains(t, gotPaths, "/update/counter/PollCount/5")
}
