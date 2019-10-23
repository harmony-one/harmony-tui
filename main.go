package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/harmony-one/harmony-tui/config"
	data "github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony-tui/src"
	"github.com/harmony-one/harmony-tui/widgets"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/spf13/viper"
)

var (
	version string
	commit  string
	builtAt string
	builtBy string
)

func main() {
	showVersion := flag.Bool("version", false, "version of the binary")

	// setting up config
	config.SetConfig()

	if *showVersion {
		fmt.Fprintf(os.Stderr,
			"Harmony (C) 2019. %v, version %v-%v (%v %v)\n",
			path.Base(os.Args[0]), version, commit, builtBy, builtAt)
		os.Exit(0)
	}

	// start go routine to tail the log file
	go src.TailZeroLogFile()

	t, err := termbox.New()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	builder := grid.New()

	// Placing widgets in grids in TUI
	builder.Add(
		grid.RowHeightPerc(30,
			grid.ColWidthPerc(50,
				grid.Widget(widgets.ChainInfo(), container.Border(linestyle.Round), container.BorderTitle(" Harmony Blockchain ")),
			),
			grid.ColWidthPerc(50,
				grid.Widget(widgets.BlockInfo(), container.Border(linestyle.Round), container.BorderTitle(" Current Block ")),
			),
		),
		grid.RowHeightPerc(30,
			grid.ColWidthPerc(50,
				grid.Widget(widgets.InstanceInfo(), container.Border(linestyle.Round), container.BorderTitle(" Harmony Node ")),
			),
			grid.ColWidthPerc(50,
				widgets.CpuLoadGrid(ctx)...,
			),
		),
		grid.RowHeightPerc(40,
			grid.ColWidthPerc(50,
				grid.Widget(widgets.GetLineChart(), container.Border(linestyle.Round), container.BorderTitle(fmt.Sprintf(" Earning Rate every %.0f sec ", viper.GetDuration("EarningRateInterval").Seconds()))),
			),
			grid.ColWidthPerc(50,
				grid.Widget(widgets.LogInfo(ctx), container.Border(linestyle.Round), container.BorderTitle(" Validator Logs ")),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		panic(err)
	}

	c, err := container.New(
		t,
		gridOpts...,
	)
	if err != nil {
		panic(err)
	}

	// logic to quite from TUI
	quit := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' || k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			data.Quitter("")
		}
	}

	// function to handle graceful exit along with exit message
	data.Quitter = func(exitMsg string) {
		cancel()
		t.Close()
		exitMsg = exitMsg + "\n"
		fmt.Fprintf(os.Stderr, exitMsg)
		time.Sleep(3 * time.Second)
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quit)); err != nil {
		panic(err)
	}
}
