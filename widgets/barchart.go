package widgets

import (
	"strconv"
	"time"

	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/barchart"
)

// TODO: This widget is not being used as of now. May be we can use or modify this file to get bar chart

const (
	barCount = 15
	barWidth = 5
)

//GetEarningRate returns barchart of earning rate per minute
func GetEarningRate() *barchart.BarChart {

	var (
		barColor   []cell.Color
		valueColor []cell.Color
		labelColor []cell.Color
		labels     []string
	)

	for i := 0; i < barCount; i++ {
		barColor = append(barColor, cell.ColorWhite)
		valueColor = append(valueColor, cell.ColorBlack)
		labelColor = append(labelColor, cell.ColorWhite)
		labels = append(labels, strconv.Itoa(i+1))
	}

	bc, err := barchart.New(
		barchart.BarColors(barColor),
		barchart.ValueColors(valueColor),
		barchart.ShowValues(),
		barchart.BarWidth(barWidth),
		barchart.BarGap(1),
		barchart.LabelColors(labelColor),
		barchart.Labels(labels),
	)
	if err != nil {
		panic(err)
	}

	go playBarChart(bc, config.EarningRateInterval)

	return bc
}

func playBarChart(bc *barchart.BarChart, delay time.Duration) {
	initialBalance := 0

	var values = []int{}
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			if data.TotalBalance != 0 && initialBalance == 0 {
				initialBalance = int(data.TotalBalance)
				continue
			}

			values = append([]int{int(data.TotalBalance) - initialBalance}, values...)
			if len(values) >= bc.ValueCapacity()-1 {
				values = values[:len(values)-1]
			}

			max := 0
			for _, val := range values {
				if val > max {
					max = val
				}
			}

			if err := bc.Values(values, max+1); err != nil {
				panic(err)
			}

			initialBalance = int(data.TotalBalance)

		}
	}
}
