package widgets

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"

	"github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony-tui/src"

	"github.com/hpcloud/tail"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

var (
	oneAsDec   = numeric.NewDec(denominations.One)
	PercentDec = numeric.NewDec(100)
)

func InstanceInfo() *text.Text {

	showEarningRate := false
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	go refreshWidget(func() {
		wrapped.Reset()
		wrapped.Write(" Harmony Version: " + data.Metadata.Version, text.WriteCellOpts(cell.FgColor(cell.ColorGreen)))
		wrapped.Write("\n ShardID    : " + strconv.FormatFloat(data.Metadata.ShardID, 'f', 0, 64) + "\n")

		if data.Bingo != "" {
			t, parseErr := time.Parse(viper.GetString("TimestampLayout"), data.Bingo)
			if parseErr == nil {
				wrapped.Write(" BINGO      : " + time.Since(t).Round(time.Second).String() + " ago\n")
				if time.Since(t).Minutes() > viper.GetFloat64("OutOfSyncTimeInMin") {
					wrapped.Write(" ")
					wrapped.Write(" Node out of sync ", text.WriteCellOpts(cell.BgColor(cell.ColorRGB24(255, 127, 80))))
				}
			}
		}

		if showEarningRate || data.EarningRate.Cmp(zeroInt) > 0 {
			showEarningRate = true
		  wrapped.Write(fmt.Sprintf("\n Earning rate : %s/%.0fs", data.EarningRate.String(), viper.GetDuration("EarningRateInterval").Seconds()))
		}

		wrapped.Write("\n\n " + data.Balance)

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

		widget.Write(" This node is connected to " + strconv.Itoa(int(data.PeerCount)) + " peers")
		widget.Write("\n NetworkID: " + data.Metadata.NetworkType)
		widget.Write("\n IsArchival: " + strconv.FormatBool(data.Metadata.ArchivalNode))
		widget.Write("\n BLS Keys: " + fmt.Sprintf("%v", data.Metadata.BLSKeys))
		widget.Write("\n Beacon Endpoint: " + data.BeaconChainEndpoint)
		widget.Write("\n Leader: " + data.LatestHeader.Leader)
		widget.Write("\n Epoch: " + strconv.Itoa(data.LatestHeader.Epoch))

		widget.Write("\n\n Announce    : " + data.Announce)
		widget.Write("\n OnAnnounce  : " + data.OnAnnounce)
		widget.Write("\n OnPrepared  : " + data.OnPrepared)
		widget.Write("\n OnCommitted : " + data.OnCommitted)
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
		widget.Write(" Block Number: " + strconv.FormatFloat(data.LatestHeader.BlockNumber, 'f', 0, 64))
		widget.Write("\n Block Size: " + strconv.Itoa(data.LatestBlock.BlockSizeInt))
		widget.Write("\n Num transactions in block: " + strconv.Itoa(data.LatestBlock.NumTransactions))
		widget.Write("\n Num staking transactions in block: " + strconv.Itoa(data.LatestBlock.NumStakingTransactions))
		widget.Write("\n Block Hash: " + data.LatestHeader.BlockHash)
		widget.Write("\n Block Epoch: " + strconv.Itoa(data.LatestHeader.Epoch))
		widget.Write("\n Block Shard: " + strconv.Itoa(int(data.LatestHeader.ShardID)))
		widget.Write("\n Block Timestamp: " + data.LatestHeader.Timestamp)
		widget.Write("\n State Root: " + data.LatestBlock.StateRoot)
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
		widget.Write(" Address: " + viper.GetString("OneAddress"))
		widget.Write("\n Elected: " + strconv.FormatBool(data.ValidatorInfo.CurrentlyInCommittee))
		widget.Write("\n EPOS Status: " + data.ValidatorInfo.EPoSStatus)
		if bootedStatus := data.ValidatorInfo.BootedStatus; bootedStatus != nil {
			widget.Write("\n Booted Status: " + *bootedStatus)
		} else {
			widget.Write("\n Booted Status: N/A")
		}

		if totalDelegated := data.ValidatorInfo.TotalDelegated; totalDelegated != nil {
			totalDelegationAsOne := numeric.NewDecFromBigInt(totalDelegated).Quo(oneAsDec)
			widget.Write("\n Total Delegation: " + totalDelegationAsOne.String())
		}

		if lifetime := data.ValidatorInfo.Lifetime; lifetime != nil {
			lifetimeRewardAsOne := numeric.NewDecFromBigInt(lifetime.BlockReward).Quo(oneAsDec)
			widget.Write("\n Lifetime Rewards: " + lifetimeRewardAsOne.String())
			widget.Write("\n Lifetime Uptime: " + data.LifetimeAvail.Mul(PercentDec).TruncateDec().String())
			widget.Write("\n APR: " + lifetime.APR.TruncateDec().String() + "%")
		} else {
			widget.Write("\n Lifetime Rewards: N/A")
			widget.Write("\n Lifetime Uptime: N/A")
			widget.Write("\n APR: N/A")
		}
		if performance := data.ValidatorInfo.Performance; performance != nil {
			widget.Write("\n Current Uptime: " + performance.CurrentSigningPercentage.Percentage.Mul(PercentDec).TruncateDec().String())
		} else {
			widget.Write("\n Current Uptime: N/A")
		}
		if winningStake := data.ValidatorInfo.EPoSWinningStake; winningStake != nil {
			widget.Write("\n Effective Stake: " + winningStake.Quo(oneAsDec).String())
		} else {
			widget.Write("\n Effective Stake: N/A")
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

	fname, err := src.GetLogFilePath("zerolog")
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
