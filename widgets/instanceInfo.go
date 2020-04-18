package widgets

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony-tui/src"

	"github.com/hpcloud/tail"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

func InstanceInfo() *text.Text {

	showEarningRate := false
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	go refreshWidget(func() {
		wrapped.Reset()
		if err := wrapped.Write(" Harmony Version: " + data.AppVersion, text.WriteCellOpts(cell.FgColor(cell.ColorGreen))); err != nil {
			panic(err)
		}

		if err := wrapped.Write("\n ShardID    : " + strconv.FormatFloat(data.ShardID, 'f', 0, 64) + "\n"); err != nil {
			panic(err)
		}

		if data.Bingo != "" {
			t, parseErr := time.Parse(viper.GetString("TimestampLayout"), data.Bingo)
			if parseErr == nil {
				if err := wrapped.Write(" BINGO      : " + time.Since(t).Round(time.Second).String() + " ago\n"); err != nil {
					panic(err)
				}
				if time.Since(t).Minutes() > viper.GetFloat64("OutOfSyncTimeInMin") {
					if err := wrapped.Write(" "); err != nil {
						panic(err)
					}
					if err := wrapped.Write(" Node out of sync ", text.WriteCellOpts(cell.BgColor(cell.ColorRGB24(255, 127, 80)))); err != nil {
						panic(err)
					}
				}
			}
		}

		if showEarningRate || data.EarningRate != 0 {
			showEarningRate = true
			if err := wrapped.Write(fmt.Sprintf("\n Earning rate : %.4f/%.0fs", data.EarningRate, viper.GetDuration("EarningRateInterval").Seconds())); err != nil {
				panic(err)
			}
		}

		if err := wrapped.Write("\n\n " + data.Balance); err != nil {
			panic(err)
		}

	})

	return wrapped
}

func ChainInfo() *text.Text {

	widget, err1 := text.New(text.WrapAtRunes())
	if err1 != nil {
		panic(err1)
	}

	go refreshWidget(func() {

		widget.Reset()

		if err := widget.Write(" This node is connected to " + strconv.Itoa(int(data.PeerCount)) + " peers"); err != nil {
			panic(err)
		}

		if err := widget.Write("\n NetworkID: " + data.NetworkID); err != nil {
			panic(err)
		}

		if err := widget.Write("\n IsArchival: " + strconv.FormatBool(data.IsArchival)); err != nil {
			panic(err)
		}

		if err := widget.Write("\n BLS Keys: " + fmt.Sprintf("%v", data.BLSKeys)); err != nil {
			panic(err)
		}

		if err := widget.Write("\n Leader: " + data.Leader); err != nil {
			panic(err)
		}

		if err := widget.Write("\n Epoch: " + strconv.FormatFloat(data.Epoch, 'f', 0, 64)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n\n Announce    : " + data.Announce); err != nil {
			panic(err)
		}
		if err := widget.Write("\n OnAnnounce  : " + data.OnAnnounce); err != nil {
			panic(err)
		}
		if err := widget.Write("\n OnPrepared  : " + data.OnPrepared); err != nil {
			panic(err)
		}
		if err := widget.Write("\n OnCommitted : " + data.OnCommitted); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Block Reward: " + data.BlockReward); err != nil {
			panic(err)
		}
	})

	return widget
}

func BlockInfo() *text.Text {

	widget, err1 := text.New(text.WrapAtRunes())
	if err1 != nil {
		panic(err1)
	}

	go refreshWidget(func() {
		widget.Reset()
		if err := widget.Write(" BlockNumber: " + strconv.FormatFloat(data.BlockNumber, 'f', 0, 64) + ", BlockSize: " + strconv.FormatInt(data.SizeInt, 10)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Num transactions in block: " + strconv.Itoa(data.NoOfTransaction)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Num staking transactions in block: " + strconv.Itoa(data.NoOfStakingTransaction)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Num staking transactions in block: " + strconv.Itoa(data.NumKeys)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n BlockHash: " + data.BlockHash); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Block Timestamp: " + data.BlockTimestamp); err != nil {
			panic(err)
		}
		if err := widget.Write("\n StateRoot: " + data.StateRoot); err != nil {
			panic(err)
		}

		if data.BlockData == nil {
			if err := widget.Write("\n BlockEpoch: no data"); err != nil {
				panic(err)
			}

			if err := widget.Write("\n Number of signers: no data"); err != nil {
				panic(err)
			}

			if err := widget.Write("\n BlockShard: no data"); err != nil {
				panic(err)
			}
		} else {
			if blockEpoch := data.BlockData["blockEpoch"]; blockEpoch != nil {
				if err := widget.Write("\n BlockEpoch: " + strconv.FormatFloat(blockEpoch.(float64), 'f', 0, 64)); err != nil {
					panic(err)
				}
			}

			if numAccounts := data.BlockData["NumAccounts"]; numAccounts != nil {
				if err := widget.Write("\n Number of signers: " + numAccounts.(string)); err != nil {
					panic(err)
				}
			}

			if blockShard := data.BlockData["blockShard"]; blockShard != nil {
				if err := widget.Write("\n BlockShard: " + strconv.FormatFloat(blockShard.(float64), 'f', 0, 64)); err != nil {
					panic(err)
				}
			}
		}
	})

	return widget
}

func ValidatorInfo() *text.Text {
	widget, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	go refreshWidget(func() {
		widget.Reset()
		if err := widget.Write(" Address: " + viper.GetString("OneAddress")); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Elected: " + strconv.FormatBool(data.Elected)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n EPOS Status: " + data.EposStatus); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Booted Status: " + data.BootedStatus ); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Total Delegation: " + data.TotalDelegation.String()); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Lifetime Rewards: " + data.LifetimeRewards.String()); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Lifetime Uptime: " + data.LifetimeAvalibility); err != nil {
			panic(err)
		}
		if err := widget.Write("\n APR: " + data.APR ); err != nil {
			panic(err)
		}
		if err := widget.Write(fmt.Sprintf("\n Validator Keys: %v", data.ValidatorBLSKeys)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Current Uptime: " + data.CurrentAvailibility); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Effective Stake: " + data.EffectiveStake); err != nil {
			panic(err)
		}
	})

	return widget
}

func LogInfo(ctx context.Context) *text.Text {
	widget, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}
	go refreshLog(ctx, widget)
	return widget
}

func refreshLog(ctx context.Context, widget *text.Text) {

	fname, err := src.GetLogFilePath("validator")
	if err != nil {
		if err = widget.Write(err.Error()); err != nil {
			panic(err)
		}
		return
	}

	t, err := tail.TailFile(fname, tail.Config{ReOpen: true, Follow: true, MustExist: false, Logger: log.New(ioutil.Discard, "", 0), Location: &tail.SeekInfo{Offset: 1, Whence: 2}})
	defer t.Cleanup()
	for line := range t.Lines {
		if err = widget.Write(line.Text); err != nil {
			panic(err)
		}
		if err = widget.Write("\n"); err != nil {
			panic(err)
		}
	}
}

func refreshWidget(f func()) {

	ticker := time.NewTicker(viper.GetDuration("WidgetInterval"))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}

func GetAppVersion() string {
	appVersion, err := src.Exec_cmd(viper.GetString("HarmonyPath") + "./harmony -version")
	if err != nil {
		return "Error collecting data"
	}
	appVersion = " App version: " + appVersion
	return appVersion
}
