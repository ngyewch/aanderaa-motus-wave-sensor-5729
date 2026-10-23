package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	aanderaa_motus_wave_sensor_5729 "github.com/ngyewch/aanderaa-motus-wave-sensor-5729"
	"github.com/urfave/cli/v3"
	"go.bug.st/serial"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	interval := cmd.Duration(intervalFlag.Name)

	tester := Tester{
		Interval: interval,
		SerialPortSupplier: func() (serial.Port, error) {
			return newSerialPort(ctx, cmd)
		},
	}
	err := tester.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

type Tester struct {
	Interval           time.Duration
	SerialPortSupplier func() (serial.Port, error)
}

func (tester Tester) Run(ctx context.Context) error {
	for {
		err := tester.doRun(ctx)
		if err != nil {
			slog.Error("error",
				slog.Any("err", err),
			)
		}
	}
}

func (tester Tester) doRun(ctx context.Context) error {
	serialPort, err := tester.SerialPortSupplier()
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port) {
		_ = serialPort.Close()
	}(serialPort)

	driver := aanderaa_motus_wave_sensor_5729.NewDriver(serialPort)
	client := aanderaa_motus_wave_sensor_5729.NewClient(driver, nil)

	err = doSetup(ctx, client, 1*time.Minute)
	if err != nil {
		return err
	}

	for {
		m, err := doSample(ctx, client, 5*time.Second)
		if err != nil {
			return err
		}

		if len(m.Quantities) > 0 {
			var logAttrs []any
			for _, quantity := range m.Quantities {
				logAttrs = append(logAttrs, slog.Float64(fmt.Sprintf("%s[%s]", quantity.Name, quantity.Units), quantity.Value))
			}
			slog.Info("measurement", logAttrs...)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(tester.Interval):
		}
	}
}

func doSetup(ctx context.Context, client *aanderaa_motus_wave_sensor_5729.Client, timeout time.Duration) error {
	slog.Info("setting up...")

	setupContext, setupCancel := context.WithTimeout(ctx, timeout)
	defer setupCancel()

	err := client.Setup(setupContext)
	if err != nil {
		return err
	}

	slog.Info("setup completed")

	return nil
}

func doSample(ctx context.Context, client *aanderaa_motus_wave_sensor_5729.Client, timeout time.Duration) (*aanderaa_motus_wave_sensor_5729.Measurement, error) {
	sampleContext, sampleCancel := context.WithTimeout(ctx, timeout)
	defer sampleCancel()

	return client.DoSample(sampleContext)
}
