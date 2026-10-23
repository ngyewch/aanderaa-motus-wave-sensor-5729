package aanderaa_motus_wave_sensor_5729

import (
	"context"
	"time"
)

const (
	DefaultPostCommandWait = 100 * time.Millisecond
	DefaultPostResetWait   = 1 * time.Second
)

type Client struct {
	driver  *Driver
	options ClientOptions
}

type ClientOptions struct {
	postCommandWait time.Duration
	postResetWait   time.Duration
}

func NewClient(driver *Driver, options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{
			postCommandWait: DefaultPostCommandWait,
			postResetWait:   DefaultPostResetWait,
		}
	}
	return &Client{
		driver:  driver,
		options: *options,
	}
}

func (client *Client) wait(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
	}
	return nil
}

func (client *Client) Setup(ctx context.Context) error {
	err := client.wait(ctx, client.options.postResetWait)
	if err != nil {
		return err
	}

	err = client.driver.Wake()
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.Stop(ctx)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.SetPasskey(ctx, 1000)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.SetMode(ctx, "Smart Sensor Terminal")
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.SetEnableText(ctx, true)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.SetEnablePolledMode(ctx, true)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.SetEnableDecimalformat(ctx, false)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.Save(ctx)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return err
	}

	err = client.driver.Reset(ctx)
	if err != nil {
		return err
	}

	err = client.wait(ctx, client.options.postResetWait)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DoSample(ctx context.Context) (*Measurement, error) {
	err := client.driver.Wake()
	if err != nil {
		return nil, err
	}

	m, err := client.driver.DoSample(ctx)
	if err != nil {
		return nil, err
	}

	err = client.wait(ctx, client.options.postCommandWait)
	if err != nil {
		return nil, err
	}

	return m, nil
}
