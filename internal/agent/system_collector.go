package agent

import (
	"log"
	"strconv"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemCollector struct {
	updater MetricsWriter
}

func NewSystemCollector(updater MetricsWriter) *SystemCollector {
	return &SystemCollector{
		updater: updater,
	}
}

func (c *SystemCollector) Poll() {
	vms, err := mem.VirtualMemory()

	if err != nil {
		log.Println("virtual memory error", err)
		return
	}

	c.updater.UpdateGauge("TotalMemory", float64(vms.Total))
	c.updater.UpdateGauge("FreeMemory", float64(vms.Free))

	perArr, err := cpu.Percent(0, true)

	if err != nil {
		log.Println("percent error", err)
		return
	}

	for ind, elem := range perArr {
		num := ind + 1
		numStr := strconv.Itoa(num)
		name := "CPUutilization" + numStr

		c.updater.UpdateGauge(name, elem)
	}
}
