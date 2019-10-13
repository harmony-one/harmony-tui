package widgets

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
)

func GetLineChart() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	)

	if err != nil {
		panic(err)
	}

	lc.Series("first", []float64{5, 20, 40},
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlue)),
		linechart.SeriesXLabels(map[int]string{
			0: "zero",
		}),
	)
	return lc
}
