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
		if err := wrapped.Write(" Harmony Version: " + data.Metadata.Version, text.WriteCellOpts(cell.FgColor(cell.ColorGreen))); err != nil {
			panic(err)
		}

		if err := wrapped.Write("\n ShardID    : " + strconv.FormatFloat(data.Metadata.ShardID, 'f', 0, 64) + "\n"); err != nil {
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

		if err := widget.Write("\n NetworkID: " + data.Metadata.NetworkType); err != nil {
			panic(err)
		}

		if err := widget.Write("\n IsArchival: " + strconv.FormatBool(data.Metadata.ArchivalNode)); err != nil {
			panic(err)
		}

		if err := widget.Write("\n BLS Keys: " + fmt.Sprintf("%v", data.Metadata.BLSKeys)); err != nil {
			panic(err)
		}

		if err := widget.Write("\n Beacon Endpoint: " + data.BeaconChainEndpoint); err != nil {
			panic(err)
		}

		if err := widget.Write("\n Leader: " + data.LatestHeader.Leader); err != nil {
			panic(err)
		}

		if err := widget.Write("\n Epoch: " + strconv.Itoa(data.LatestHeader.Epoch)); err != nil {
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
		if err := widget.Write(" BlockNumber: " + strconv.FormatFloat(data.LatestHeader.BlockNumber, 'f', 0, 64) + ", BlockSize: " + strconv.Itoa(data.LatestBlock.BlockSizeInt)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Num transactions in block: " + strconv.Itoa(data.LatestBlock.NumTransactions)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Num staking transactions in block: " + strconv.Itoa(data.LatestBlock.NumStakingTransactions)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n BlockHash: " + data.LatestHeader.BlockHash); err != nil {
			panic(err)
		}
		if err := widget.Write("\n Block Timestamp: " + data.LatestHeader.Timestamp); err != nil {
			panic(err)
		}
		if err := widget.Write("\n StateRoot: " + data.LatestBlock.StateRoot); err != nil {
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
		if err := widget.Write("\n Elected: " + strconv.FormatBool(data.ValidatorInfo.CurrentlyInCommittee)); err != nil {
			panic(err)
		}
		if err := widget.Write("\n EPOS Status: " + data.ValidatorInfo.EPoSStatus); err != nil {
			panic(err)
		}
		if bootedStatus := data.ValidatorInfo.BootedStatus; bootedStatus != nil {
			if err := widget.Write("\n Booted Status: " + *bootedStatus); err != nil {
				panic(err)
			}
		} else {
			if err := widget.Write("\n Booted Status: N/A"); err != nil {
				panic(err)
			}
		}
		if err := widget.Write("\n Total Delegation: " + data.ValidatorInfo.TotalDelegated.String()); err != nil {
			panic(err)
		}
		if lifetime := data.ValidatorInfo.Lifetime; lifetime != nil {
			if err := widget.Write("\n Lifetime Rewards: " + lifetime.BlockReward.String()); err != nil {
				panic(err)
			}
			if err := widget.Write("\n Lifetime Uptime: " + data.LifetimeAvail); err != nil {
				panic(err)
			}
			if err := widget.Write("\n APR: " + lifetime.APR.String()); err != nil {
				panic(err)
			}
		} else {
			if err := widget.Write("\n Lifetime Rewards: N/A"); err != nil {
				panic(err)
			}
			if err := widget.Write("\n Lifetime Uptime: N/A"); err != nil {
				panic(err)
			}
			if err := widget.Write("\n APR: N/A"); err != nil {
				panic(err)
			}
		}
		if performance := data.ValidatorInfo.Performance; performance != nil {
			if err := widget.Write("\n Current Uptime: " + performance.CurrentSigningPercentage.Percentage.String()); err != nil {
				panic(err)
			}
		} else {
			if err := widget.Write("\n Current Uptime: N/A"); err != nil {
				panic(err)
			}
		}
		if winningStake := data.ValidatorInfo.EPoSWinningStake; winningStake != nil {
			if err := widget.Write("\n Effective Stake: " + data.ValidatorInfo.EPoSWinningStake.TruncateInt().String()); err != nil {
				panic(err)
			}
		} else {
			if err := widget.Write("\n Effective Stake: N/A"); err != nil {
				panic(err)
			}
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
