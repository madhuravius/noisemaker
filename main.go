package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2" // imports as package "cli"

	"trashbin/internal"
)

func main() {
	app := &cli.App{
		Name:  "trashbin",
		Usage: "needlessly consume resources and throw it in the bin",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{""},
				Usage:   "start the trashcan",
				Action: func(cCtx *cli.Context) error {
					log.Println("Starting trashbin run")
					internal.RunLoop(cCtx)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  internal.CpuContextKey,
				Value: 50,
				Usage: "as a percentage, specify the percentage of CPU you would like to use. Max: 99",
			}, &cli.Float64Flag{
				Name:  internal.MemoryContextKey,
				Value: 50,
				Usage: "as a percentage, specify the percentage of RAM you would like to use. Max: 99",
			}, &cli.IntFlag{
				Name:  internal.BandwidthContextKey,
				Value: 5,
				Usage: "in MBps, specify how much bandwidth you want to use upload/download",
			}, &cli.IntFlag{
				Name:  internal.PortContextKey,
				Value: 3000,
				Usage: "specify port you wish to run the web server for stressing",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
