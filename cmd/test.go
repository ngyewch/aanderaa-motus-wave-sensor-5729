package main

import (
	"context"
	"log/slog"
	"time"

	aanderaa_motus_wave_sensor_5729 "github.com/ngyewch/aanderaa-motus-wave-sensor-5729"
	"github.com/urfave/cli/v3"
	"go.bug.st/serial"
)

func doTest(ctx context.Context, cmd *cli.Command) error {
	serialPort, err := newSerialPort(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port) {
		_ = serialPort.Close()
	}(serialPort)

	interval := cmd.Duration(intervalFlag.Name)

	driver := aanderaa_motus_wave_sensor_5729.NewDriver(serialPort)

	tester := NewTester(driver, interval)
	err = tester.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

type Tester struct {
	driver   *aanderaa_motus_wave_sensor_5729.Driver
	interval time.Duration
}

func NewTester(driver *aanderaa_motus_wave_sensor_5729.Driver, interval time.Duration) *Tester {
	return &Tester{
		driver:   driver,
		interval: interval,
	}
}

func (tester *Tester) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
	}

	err := tester.driver.Wake()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.Stop()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.SetPasskey(1000)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.SetMode("Smart Sensor Terminal")
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.SetEnableText(true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.SetEnablePolledMode(true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.SetEnableDecimalformat(true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.Save()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	err = tester.driver.Reset()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
	}

	for {
		err = tester.driver.Wake()
		if err != nil {
			return err
		}

		lines, err := tester.driver.DoSample()
		if err != nil {
			return err
		}
		slog.Info("Do Sample",
			slog.Any("lines", lines),
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(tester.interval):
		}
	}

	return nil
}
