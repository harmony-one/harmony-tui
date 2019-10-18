package widgets

import (
	"time"

	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
)

// TODO: This widget is not being used as of now. May be we can use or modify this file to get line chart


// GetLineChart retunrs linechart of total balance in one account
func GetLineChart() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorBlack)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XAxisUnscaled(),
		linechart.YAxisAdaptive(),
	)

	if err != nil {
		panic(err)
	}
	go playLineCharrt(lc)
	return lc
}

func playLineCharrt(lc *linechart.LineChart) {
	ticker := time.NewTicker(config.LinechartInterval)
	defer ticker.Stop()
	values := []float64{}
	for {
		select {
		case <-ticker.C:
			values = append(values, data.TotalBalance)
			if len(values) > 10 {
				values = values[1:]
			}

			lc.Series("amount", values,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlue)),
				linechart.SeriesXLabels(map[int]string{
					0: "time",
				}),
			)
		}
	}
}
