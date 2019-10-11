package main

import (
	"context"
	"flag"
	"os"
	"path"
	"fmt"

	"github.com/harmony-one/harmony-tui/src/data"
	"github.com/harmony-one/harmony-tui/widgets"
	"github.com/harmony-one/harmony-tui/src"
	"github.com/harmony-one/harmony-tui/config"
	
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/container/grid"
)

var (
	version string
	commit  string
	builtAt string
	builtBy string
)

func main() {
	// setting up config varibales
	env := flag.String("env", "local", "environment of system binary is running on option 1- \"local\" option 2- \"ec2\"")
	showVersion := flag.Bool("version", false, "version of the binary")
	flag.Parse()
	config.SetConfig(*env)
	
	if *showVersion {
		fmt.Fprintf(os.Stderr,
			"Harmony (C) 2019. %v, version %v-%v (%v %v)\n",
			path.Base(os.Args[0]), version, commit, builtBy, builtAt)
		os.Exit(0)
	}
	// start go routine to tail the log file
	go src.TailZeroLogFile()

	t, err := termbox.New()
	if err!=nil {
		panic(err)
	}
	
	ctx, cancel := context.WithCancel(context.Background())

	builder := grid.New()

	// Placing widgets in grids in TUI
	builder.Add(
		grid.RowHeightPerc(30,
			grid.ColWidthPerc(50,
				grid.Widget(widgets.ChainInfo(), container.Border(linestyle.Round), container.BorderTitle("Harmony Blockchain")),
			),
			grid.ColWidthPerc(50,
				grid.Widget(widgets.BlockInfo(), container.Border(linestyle.Round), container.BorderTitle("Current Block")),
			),
		),
		grid.RowHeightPerc(30,
			grid.ColWidthPerc(50,
				grid.Widget(widgets.InstanceInfo(), container.Border(linestyle.Round), container.BorderTitle("Harmony Node")),
			),
			grid.ColWidthPerc(50,
				widgets.CpuLoadGrid(ctx)..., 
			),
		),
		grid.RowHeightPerc(40,
			grid.Widget(widgets.LogInfo(ctx), container.Border(linestyle.Round), container.BorderTitle("Validator Logs")),
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
		if k.Key == 'q' || k.Key == 'Q' || k.Key == keyboard.KeyEsc {
			data.Quitter("")
		}
	}
	
	// function to handle graceful exit along with exit message
	data.Quitter = func(exitMsg string) {
		cancel()
		t.Close()
		exitMsg = exitMsg + "\n"
		fmt.Fprintf(os.Stderr, exitMsg)
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quit)); err != nil {
		panic(err)
	}
}
