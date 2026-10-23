package aanderaa_motus_wave_sensor_5729

import (
	"errors"
	"syscall"
	"time"

	"go.bug.st/serial"
)

type serialPortWrapper struct {
	serialPort serial.Port
}

func newSerialPortWrapper(serialPort serial.Port) *serialPortWrapper {
	return &serialPortWrapper{
		serialPort: serialPort,
	}
}

func (port *serialPortWrapper) SetMode(mode *serial.Mode) error {
	for {
		err := port.serialPort.SetMode(mode)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) Read(b []byte) (int, error) {
	for {
		n, err := port.serialPort.Read(b)
		if !errors.Is(err, syscall.EINTR) {
			return n, err
		}
	}
}

func (port *serialPortWrapper) Write(b []byte) (int, error) {
	for {
		n, err := port.serialPort.Write(b)
		if !errors.Is(err, syscall.EINTR) {
			return n, err
		}
	}
}

func (port *serialPortWrapper) Drain() error {
	for {
		err := port.serialPort.Drain()
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) ResetInputBuffer() error {
	for {
		err := port.serialPort.ResetInputBuffer()
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) ResetOutputBuffer() error {
	for {
		err := port.serialPort.ResetOutputBuffer()
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) SetDTR(dtr bool) error {
	for {
		err := port.serialPort.SetDTR(dtr)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) SetRTS(rts bool) error {
	for {
		err := port.serialPort.SetRTS(rts)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	for {
		bits, err := port.serialPort.GetModemStatusBits()
		if !errors.Is(err, syscall.EINTR) {
			return bits, err
		}
	}
}

func (port *serialPortWrapper) SetReadTimeout(t time.Duration) error {
	for {
		err := port.serialPort.SetReadTimeout(t)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *serialPortWrapper) Close() error {
	return port.serialPort.Close()
}

func (port *serialPortWrapper) Break(duration time.Duration) error {
	for {
		err := port.serialPort.Break(duration)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}
