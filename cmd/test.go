package main

import (
	"context"
	"log/slog"
	"strings"
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

	iterationCount := 0
	for {
		err = func() error {
			iterationCount++
			slog.Info("iteration",
				slog.Int("count", iterationCount),
			)

			err = tester.Setup(ctx)
			if err != nil {
				return err
			}

			t0 := time.Now()
			for {
				m, err := tester.DoSample(ctx)
				if err != nil {
					return err
				}

				if len(m.Quantities) == 0 {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(interval):
						continue
					}
				}

				duration := time.Since(t0)
				var quantityStrings []string
				for _, quantity := range m.Quantities {
					quantityStrings = append(quantityStrings, quantity.String())
				}
				slog.Info("measurement",
					slog.Duration("duration", duration),
					slog.String("quantities", strings.Join(quantityStrings, " / ")),
				)
				break
			}

			return nil
		}()
		if err != nil {
			slog.Error("error",
				slog.Any("err", err),
			)
		}
	}
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

func (tester *Tester) wait(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
	}
	return nil
}

func (tester *Tester) Setup(ctx context.Context) error {
	err := tester.wait(ctx, 1*time.Second)
	if err != nil {
		return err
	}

	err = tester.driver.Wake()
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.Stop()
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.SetPasskey(1000)
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.SetMode("Smart Sensor Terminal")
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.SetEnableText(true)
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.SetEnablePolledMode(true)
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.SetEnableDecimalformat(false)
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.Save()
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	err = tester.driver.Reset()
	if err != nil {
		return err
	}

	err = tester.wait(ctx, 1*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (tester *Tester) DoSample(ctx context.Context) (*aanderaa_motus_wave_sensor_5729.Measurement, error) {
	err := tester.driver.Wake()
	if err != nil {
		return nil, err
	}

	m, err := tester.driver.DoSample()
	if err != nil {
		return nil, err
	}

	err = tester.wait(ctx, 100*time.Millisecond)
	if err != nil {
		return nil, err
	}

	return m, nil
}
