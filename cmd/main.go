package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

var (
	version string

	serialPortFlag = &cli.StringFlag{
		Name:     "serial-port",
		Usage:    "serial port",
		Required: true,
		Sources:  cli.EnvVars("SERIAL_PORT"),
	}
	baudRateFlag = &cli.IntFlag{
		Name:    "baud-rate",
		Usage:   "baud rate",
		Value:   115200,
		Sources: cli.EnvVars("BAUD_RATE"),
	}
	dataBitsFlag = &cli.IntFlag{
		Name:    "data-bits",
		Usage:   "data bits",
		Value:   8,
		Sources: cli.EnvVars("DATA_BITS"),
	}
	parityFlag = &cli.StringFlag{
		Name:    "parity",
		Usage:   "parity",
		Value:   "N",
		Sources: cli.EnvVars("PARITY"),
		Action: func(ctx context.Context, cmd *cli.Command, s string) error {
			switch s {
			case "N", "O", "E", "M", "S":
				return nil
			default:
				return fmt.Errorf("unknown parity: %s", s)
			}
		},
	}
	stopBitsFlag = &cli.StringFlag{
		Name:    "stop-bits",
		Usage:   "stop bits",
		Value:   "1",
		Sources: cli.EnvVars("STOP_BITS"),
		Action: func(ctx context.Context, cmd *cli.Command, s string) error {
			switch s {
			case "1", "1.5", "2":
				return nil
			default:
				return fmt.Errorf("unknown stop bits: %s", s)
			}
		},
	}
	intervalFlag = &cli.DurationFlag{
		Name:    "interval",
		Usage:   "interval",
		Value:   1 * time.Minute,
		Sources: cli.EnvVars("INTERVAL"),
	}

	app = &cli.Command{
		Name:    "aanderaa-motus-wave-sensor-5729",
		Usage:   "Aanderaa Motus Wave Sensor 5729",
		Version: version,
		Action:  nil,
		Commands: []*cli.Command{
			{
				Name:   "test",
				Usage:  "test",
				Action: doTest,
				Flags: []cli.Flag{
					serialPortFlag,
					baudRateFlag,
					dataBitsFlag,
					parityFlag,
					stopBitsFlag,
					intervalFlag,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		slog.Error("error",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}
