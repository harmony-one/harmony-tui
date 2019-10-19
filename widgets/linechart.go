package widgets

import (
	"time"

	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
)

// GetLineChart retunrs linechart of total balance in one account
func GetLineChart() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.XAxisUnscaled(),
		linechart.YAxisCustomScale(0.00, 0.80),
	)

	if err != nil {
		panic(err)
	}
	go playLineChart(lc)
	return lc
}

func playLineChart(lc *linechart.LineChart) {
	initialBalance := 0.00
	ticker := time.NewTicker(config.EarningRateInterval)
	defer ticker.Stop()
	values := []float64{}
	for {
		select {
		case <-ticker.C:
			if data.TotalBalance != 0 && initialBalance == 0 {
				initialBalance = data.TotalBalance
				continue
			}
			data.EarningRate = data.TotalBalance - initialBalance
			values = append(values, data.EarningRate)
			if len(values) > 15 {
				values = values[1:]
			}

			lc.Series("amount", values,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorWhite)),
				linechart.SeriesXLabels(map[int]string{
					0: "time",
				}),
			)
			initialBalance = data.TotalBalance
		}
	}
}
