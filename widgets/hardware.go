package widgets

import (
	"context"
	"time"
	
	"github.com/harmony-one/harmony-tui/config"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/container/grid"
)

type fn func() int

func cpuUsage() int {
	progress, err := cpu.Percent(2000*time.Millisecond,false)
	if err!=nil {
		panic(err)
	}
	return int(progress[0])
}

func memoryUsage() int {
	usage, err := mem.VirtualMemory()
	if err!= nil {
		panic(err)
	}
	return int(usage.UsedPercent)
}

func diskUsage() int {
	usage, err := disk.Usage("/")
	if err!=nil {
		panic(err)
	}
	return int(usage.UsedPercent)
}

// TODO: can be added later
/*func netUsage() int {
	usage, err := net.IOCounters(false)
	fmt.Println(usage[0].BytesRecv)
	if err!=nil {
		panic(err)
	}
	return int(usage[0].BytesRecv)
}*/

func refresh(ctx context.Context, gauge *gauge.Gauge, f fn) {

	ticker := time.NewTicker(config.SystemStatsInterval)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			if err := gauge.Percent(f()); err != nil {
				panic(err)
			}

		case <- ctx.Done():
			return
		}
	}
}


func CpuLoadGrid(ctx context.Context) ([]grid.Element){

	// create cpu gauge
	cpuGauage, err := gauge.New(
		gauge.Height(1),
		gauge.Border(linestyle.Light),
		gauge.Color(cell.ColorWhite),
		gauge.BorderTitle("CPU Usage"),
	)
	if err!=nil {
		panic(err)
	}
	go refresh(ctx, cpuGauage, cpuUsage)


	// create memory gauge
	memGauge, err := gauge.New(
		gauge.Height(1),
		gauge.Border(linestyle.Light),
		gauge.Color(cell.ColorWhite),
		gauge.BorderTitle("Memory Usage"),
	)
	if err!= nil {
		panic(err)
	}
	go refresh(ctx, memGauge, memoryUsage)


	// create disk gauge
	diskGauge, err := gauge.New(
		gauge.Height(1),
		gauge.Border(linestyle.Light),
		gauge.Color(cell.ColorWhite),
		gauge.BorderTitle("Disk Usage"),
	)
	if err!=nil {
		panic(err)
	}
	go refresh(ctx, diskGauge, diskUsage)

	// create grid structure
	el1 := grid.RowHeightPerc(33, grid.Widget(cpuGauage))
	el2 := grid.RowHeightPerc(34, grid.Widget(memGauge))
	el3 := grid.RowHeightPerc(33, grid.Widget(diskGauge))
	
	return []grid.Element{el1,el2,el3}
}