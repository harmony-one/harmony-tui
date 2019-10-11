package widgets

import (
	"strconv"
	"context"
	"time"
	"io/ioutil"
	"log"

	"harmony-tui/src"
	"harmony-tui/src/data"
	"harmony-tui/config"

	"github.com/hpcloud/tail"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

func InstanceInfo() *text.Text {
	appVersion, err := src.Exec_cmd(config.HarmonyPath + "./harmony -version")
	if err!=nil {
		appVersion = "Error collecting data\n"
	}
	appVersion = "App version: " + appVersion

	// TODO: stop using wallet.sh to get balances
	balances, err := src.Exec_cmd(config.HarmonyPath + "./wallet.sh balances")
	if err!=nil {
		balances = "Error collecting data"
	}
	
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	if err := wrapped.Write(appVersion + "", text.WriteCellOpts(cell.FgColor(cell.ColorGreen))); err != nil {
		panic(err)
	}

	if err := wrapped.Write("ShardID: " + strconv.FormatFloat(data.ShardID, 'f', 0, 64) + "\n"); err != nil {
		panic(err)
	}

	if err:= wrapped.Write("\n" + balances); err != nil {
		panic(err)
	}

	return wrapped
}

func ChainInfo() *text.Text {

	widget, err1 := text.New(text.WrapAtRunes())
	if err1!=nil {
		panic(err1)
	}

	go refreshWidget( func(){
		
		widget.Reset()
		
		if err:= widget.Write("This node is connected to " + strconv.Itoa(int(data.PeerCount)) + " peers"); err != nil {
			panic(err)
		}

		if err:=widget.Write("\nLeader: " + data.Leader); err!=nil {
			panic(err)
		}

		if err:= widget.Write("\nEpoch: " + strconv.FormatFloat(data.Epoch, 'f', 0, 64)); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nAnnounce: " + data.Announce); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nOnAnnounce: " + data.OnAnnounce); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nOnPrepared: " + data.OnPrepared); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nOnCommitted: " + data.OnCommitted); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nBlock Reward: " + data.BlockReward); err != nil {
			panic(err)
		}
	})
	
	return widget
}

func BlockInfo() *text.Text {
	
	widget, err1 := text.New(text.WrapAtRunes())
	if err1!=nil {
		panic(err1)
	}

	go refreshWidget( func(){
		widget.Reset()
		if err:= widget.Write("Current BlockNumber: " + strconv.FormatFloat(data.BlockNumber, 'f', 0, 64) + ", size: " + strconv.FormatInt(data.SizeInt,10)); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nNum transactions in block: " + strconv.Itoa(data.NoOfTransaction)); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nCurrent BlockHash: " + data.BlockHash); err != nil {
			panic(err)
		}
		if err:= widget.Write("\nstateRoot: " + data.StateRoot); err != nil {
			panic(err)
		}

		if data.BlockData == nil {
			if err:= widget.Write("\nblockEpoch: no data"); err != nil {
				panic(err)
			}
	
			if err:= widget.Write("\nNumAccounts: no data"); err != nil {
				panic(err)
			}
	
			if err:= widget.Write("\nblockShard: no data"); err != nil {
				panic(err)
			}
		} else {
			if err:= widget.Write("\nblockEpoch: " + strconv.FormatFloat(data.BlockData["blockEpoch"].(float64), 'f', 0, 64)); err != nil {
				panic(err)
			}
	
			if err:= widget.Write("\nNumAccounts: " + data.BlockData["NumAccounts"].(string)); err != nil {
				panic(err)
			}
	
			if err:= widget.Write("\nblockShard: " + strconv.FormatFloat(data.BlockData["blockShard"].(float64), 'f', 0, 64)); err != nil {
				panic(err)
			}
		}
	})
	
	return widget
}

func LogInfo(ctx context.Context) *text.Text {
	widget, err := text.New(text.RollContent(), text.WrapAtWords())
	if err!=nil {
		panic(err)
	}
	go refreshLog(ctx, widget)
	return widget
}

func refreshLog(ctx context.Context, widget *text.Text) {

	ticker := time.NewTicker(3000*time.Millisecond)
	defer ticker.Stop()
	fname := src.GetLogFilePath("validator")

	for {
		select {
		case <- ticker.C:
			t, err := tail.TailFile(fname, tail.Config{Follow: true, MustExist: false, Logger: log.New(ioutil.Discard, "", 0), Location: &tail.SeekInfo{Offset: 1, Whence: 2}})
			
			for line := range t.Lines {
				if err = widget.Write(line.Text); err!=nil {
					panic(err)
				}
				if err= widget.Write("\n"); err!=nil {
					panic(err)
				}
			}
		case <- ctx.Done():
			return
		}
	}
}


func refreshWidget( f func()) {

	ticker := time.NewTicker(500*time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			f()
		}
	}
}
