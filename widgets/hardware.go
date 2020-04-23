package widgets

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type fn func() int

func CpuUsage() int {
	progress, err := cpu.Percent(2000*time.Millisecond, false)
	if err != nil {
		return 0
	}
	return int(progress[0])
}

func MemoryUsage() int {
	usage, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return int(usage.UsedPercent)
}

func DiskUsage() int {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0
	}
	return int(usage.UsedPercent)
}

func refresh(ctx context.Context, gauge *gauge.Gauge, f fn) {

	ticker := time.NewTicker(viper.GetDuration("SystemStatsInterval"))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gauge.Percent(f())

		case <-ctx.Done():
			return
		}
	}
}

func CpuLoadGrid(ctx context.Context) []grid.Element {

	// create cpu gauge
	cpuGauage, _ := gauge.New(
		gauge.Height(1),
		gauge.Border(linestyle.Light),
		gauge.Color(cell.ColorWhite),
		gauge.BorderTitle(" CPU Usage "),
	)
	go refresh(ctx, cpuGauage, CpuUsage)

	// create memory gauge
	memGauge, _ := gauge.New(
		gauge.Height(1),
		gauge.Border(linestyle.Light),
		gauge.Color(cell.ColorWhite),
		gauge.BorderTitle(" Memory Usage "),
	)
	go refresh(ctx, memGauge, MemoryUsage)

	// create disk gauge
	diskGauge, _ := gauge.New(
		gauge.Height(1),
		gauge.Border(linestyle.Light),
		gauge.Color(cell.ColorWhite),
		gauge.BorderTitle(" Disk Usage "),
	)
	go refresh(ctx, diskGauge, DiskUsage)

	// create grid structure
	el1 := grid.RowHeightPerc(33, grid.Widget(cpuGauage))
	el2 := grid.RowHeightPerc(34, grid.Widget(memGauge))
	el3 := grid.RowHeightPerc(33, grid.Widget(diskGauge))

	return []grid.Element{el1, el2, el3}
}
