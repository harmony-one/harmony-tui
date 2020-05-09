package widgets

import (
	"math/big"
	"time"

	"github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony/numeric"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/spf13/viper"
)

var (
	zeroInt = big.NewInt(0)
)

// GetLineChart retunrs linechart of total balance in one account
func GetLineChart() *linechart.LineChart {
	lc, _ := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.XAxisUnscaled(),
		linechart.YAxisCustomScale(0.00, 0.80),
	)
	go playLineChart(lc)
	return lc
}

func playLineChart(lc *linechart.LineChart) {
	initialRewards := big.NewInt(0)
	ticker := time.NewTicker(viper.GetDuration("EarningRateInterval"))
	defer ticker.Stop()
	values := []float64{}
	for {
		select {
		case <-ticker.C:
			if data.ValidatorInfo.Lifetime.BlockReward.Cmp(zeroInt) > 0 && initialRewards.Cmp(zeroInt) == 0 {
				initialRewards = data.ValidatorInfo.Lifetime.BlockReward
				continue
			}
			data.EarningRate = numeric.NewDecFromBigInt(data.ValidatorInfo.Lifetime.BlockReward).Sub(numeric.NewDecFromBigInt(initialRewards))
			data.EarningRate = data.EarningRate.Quo(oneAsDec)
			earningRate, _, _ := new(big.Float).Parse(data.EarningRate.String(), 10)
			floatRate, _ := earningRate.Float64()
			values = append(values, floatRate)
			if len(values) > 15 {
				values[0] = 0
				values = values[1:]
			}

			lc.Series("amount", values,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorWhite)),
				linechart.SeriesXLabels(map[int]string{
					0: "time",
				}),
			)
			initialRewards = data.ValidatorInfo.Lifetime.BlockReward
		}
	}
}
