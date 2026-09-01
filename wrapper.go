package aanderaa_motus_wave_sensor_5729

import (
	"errors"
	"syscall"
	"time"

	"go.bug.st/serial"
)

type SerialPortWrapper struct {
	serialPort serial.Port
}

func NewSerialPortWrapper(serialPort serial.Port) *SerialPortWrapper {
	return &SerialPortWrapper{
		serialPort: serialPort,
	}
}

func (port *SerialPortWrapper) SetMode(mode *serial.Mode) error {
	for {
		err := port.serialPort.SetMode(mode)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) Read(b []byte) (int, error) {
	for {
		n, err := port.serialPort.Read(b)
		if !errors.Is(err, syscall.EINTR) {
			return n, err
		}
	}
}

func (port *SerialPortWrapper) Write(b []byte) (int, error) {
	for {
		n, err := port.serialPort.Write(b)
		if !errors.Is(err, syscall.EINTR) {
			return n, err
		}
	}
}

func (port *SerialPortWrapper) Drain() error {
	for {
		err := port.serialPort.Drain()
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) ResetInputBuffer() error {
	for {
		err := port.serialPort.ResetInputBuffer()
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) ResetOutputBuffer() error {
	for {
		err := port.serialPort.ResetOutputBuffer()
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) SetDTR(dtr bool) error {
	for {
		err := port.serialPort.SetDTR(dtr)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) SetRTS(rts bool) error {
	for {
		err := port.serialPort.SetRTS(rts)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	for {
		bits, err := port.serialPort.GetModemStatusBits()
		if !errors.Is(err, syscall.EINTR) {
			return bits, err
		}
	}
}

func (port *SerialPortWrapper) SetReadTimeout(t time.Duration) error {
	for {
		err := port.serialPort.SetReadTimeout(t)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (port *SerialPortWrapper) Close() error {
	return port.serialPort.Close()
}

func (port *SerialPortWrapper) Break(duration time.Duration) error {
	for {
		err := port.serialPort.Break(duration)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}
