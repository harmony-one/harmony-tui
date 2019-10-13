package widgets

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony-tui/src"

	"github.com/hpcloud/tail"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

func InstanceInfo() *text.Text {

	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	data.AppVersion, err = src.Exec_cmd(config.HarmonyPath + "./harmony -version")
	if err != nil {
		data.AppVersion = "Error collecting data\n"
	}
	data.AppVersion = " App version: " + data.AppVersion

	go refreshWidget(func() {
		wrapped.Reset()
		if err := wrapped.Write(data.AppVersion+"", text.WriteCellOpts(cell.FgColor(cell.ColorGreen))); err != nil {
			panic(err)
		}

		if err := wrapped.Write(" ShardID: " + strconv.FormatFloat(data.ShardID, 'f', 0, 64) + "\n"); err != nil {
			panic(err)
		}

		if err := wrapped.Write("\n " + data.Balance); err != nil {
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
		if err := widget.Write("\n BlockHash: " + data.BlockHash); err != nil {
			panic(err)
		}
		if err := widget.Write("\n StateRoot: " + data.StateRoot); err != nil {
			panic(err)
		}

		if data.BlockData == nil {
			if err := widget.Write("\n BlockEpoch: no data"); err != nil {
				panic(err)
			}

			if err := widget.Write("\n Number if signers: no data"); err != nil {
				panic(err)
			}

			if err := widget.Write("\n BlockShard: no data"); err != nil {
				panic(err)
			}
		} else {
			if blockEpoch := data.BlockData["blockEpoch"]; blockEpoch!=nil{
				if err := widget.Write("\n BlockEpoch: " + strconv.FormatFloat(blockEpoch.(float64), 'f', 0, 64)); err != nil {
					panic(err)
				}
			}

			if numAccounts := data.BlockData["NumAccounts"]; numAccounts!=nil{
				if err := widget.Write("\n Number of signers: " + numAccounts.(string)); err != nil {
					panic(err)
				}
			}

			if blockShard := data.BlockData["blockShard"]; blockShard!=nil{
				if err := widget.Write("\n BlockShard: " + strconv.FormatFloat(blockShard.(float64), 'f', 0, 64)); err != nil {
					panic(err)
				}
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
	
	t, err := tail.TailFile(fname, tail.Config{Follow: true, MustExist: false, Logger: log.New(ioutil.Discard, "", 0), Location: &tail.SeekInfo{Offset: 1, Whence: 2}})

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

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}
