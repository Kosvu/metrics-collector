package agent

import (
	"math/rand/v2"
	"runtime"
)

type Collector struct {
	updater MetricsWriter
}

type MetricsWriter interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
}

func NewCollector(updater MetricsWriter) *Collector {
	return &Collector{
		updater: updater,
	}
}

func (c *Collector) Poll() {
	var r runtime.MemStats
	runtime.ReadMemStats(&r)

	c.updater.UpdateGauge("Alloc", float64(r.Alloc))
	c.updater.UpdateGauge("BuckHashSys", float64(r.BuckHashSys))
	c.updater.UpdateGauge("Frees", float64(r.Frees))
	c.updater.UpdateGauge("GCCPUFraction", float64(r.GCCPUFraction))
	c.updater.UpdateGauge("GCSys", float64(r.GCSys))
	c.updater.UpdateGauge("HeapAlloc", float64(r.HeapAlloc))
	c.updater.UpdateGauge("HeapIdle", float64(r.HeapIdle))
	c.updater.UpdateGauge("HeapInuse", float64(r.HeapInuse))
	c.updater.UpdateGauge("HeapObjects", float64(r.HeapObjects))
	c.updater.UpdateGauge("HeapReleased", float64(r.HeapReleased))
	c.updater.UpdateGauge("HeapSys", float64(r.HeapSys))
	c.updater.UpdateGauge("LastGC", float64(r.LastGC))
	c.updater.UpdateGauge("Lookups", float64(r.Lookups))
	c.updater.UpdateGauge("MCacheInuse", float64(r.MCacheInuse))
	c.updater.UpdateGauge("MCacheSys", float64(r.MCacheSys))
	c.updater.UpdateGauge("MSpanInuse", float64(r.MSpanInuse))
	c.updater.UpdateGauge("MSpanSys", float64(r.MSpanSys))
	c.updater.UpdateGauge("Mallocs", float64(r.Mallocs))
	c.updater.UpdateGauge("NextGC", float64(r.NextGC))
	c.updater.UpdateGauge("NumForcedGC", float64(r.NumForcedGC))
	c.updater.UpdateGauge("NumGC", float64(r.NumGC))
	c.updater.UpdateGauge("OtherSys", float64(r.OtherSys))
	c.updater.UpdateGauge("PauseTotalNs", float64(r.PauseTotalNs))
	c.updater.UpdateGauge("StackInuse", float64(r.StackInuse))
	c.updater.UpdateGauge("StackSys", float64(r.StackSys))
	c.updater.UpdateGauge("Sys", float64(r.Sys))
	c.updater.UpdateGauge("TotalAlloc", float64(r.TotalAlloc))

	c.updater.UpdateCounter("PollCount", int64(1))
	c.updater.UpdateGauge("RandomValue", rand.Float64())
}
