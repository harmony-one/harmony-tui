package main

import (
	"github.com/mum4k/termdash/keyboard"
	"context"
	"flag"

	"text-based-ui/widgets"
	"text-based-ui/src"
	"text-based-ui/config"
	
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/container/grid"
)

func main() {
	// setting up config varibales
	env := flag.String("env", "local", "environment")
	flag.Parse()
	config.SetConfig(*env)
	
	// start go routine to tail the log file
	go src.TailZeroLogFile()

	t, err := termbox.New()
	if err!=nil {
		panic(err)
	}
	defer t.Close()

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
	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' || k.Key == keyboard.KeyEsc {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}
