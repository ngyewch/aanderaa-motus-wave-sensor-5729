package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"go.bug.st/serial"
)

func newSerialPort(ctx context.Context, cmd *cli.Command) (serial.Port, error) {
	serialPort := cmd.String(serialPortFlag.Name)
	baudRate := cmd.Int(baudRateFlag.Name)
	dataBits := cmd.Int(dataBitsFlag.Name)
	parityString := cmd.String(parityFlag.Name)
	stopBitsString := cmd.String(stopBitsFlag.Name)

	parity, err := func(s string) (serial.Parity, error) {
		switch s {
		case "N":
			return serial.NoParity, nil
		case "E":
			return serial.EvenParity, nil
		case "O":
			return serial.OddParity, nil
		case "M":
			return serial.MarkParity, nil
		case "S":
			return serial.SpaceParity, nil
		default:
			return serial.NoParity, fmt.Errorf("invalid parity: %s", s)
		}
	}(parityString)
	if err != nil {
		return nil, err
	}

	stopBits, err := func(s string) (serial.StopBits, error) {
		switch s {
		case "1":
			return serial.OneStopBit, nil
		case "1.5":
			return serial.OnePointFiveStopBits, nil
		case "2":
			return serial.TwoStopBits, nil
		default:
			return serial.OneStopBit, fmt.Errorf("invalid stop bits: %s", s)
		}
	}(stopBitsString)
	if err != nil {
		return nil, err
	}

	mode := serial.Mode{
		BaudRate: baudRate,
		DataBits: dataBits,
		Parity:   parity,
		StopBits: stopBits,
	}

	port, err := serial.Open(serialPort, &mode)
	if err != nil {
		return nil, err
	}

	return port, nil
}
