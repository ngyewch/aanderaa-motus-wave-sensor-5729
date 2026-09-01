package aanderaa_motus_wave_sensor_5729

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go.bug.st/serial"
)

const (
	lineTerminator = "\r\n"
)

type Driver struct {
	serialPort serial.Port
}

func NewDriver(serialPort serial.Port) *Driver {
	return &Driver{
		serialPort: serialPort,
	}
}

func (driver *Driver) Wake() error {
	return driver.wake(100 * time.Millisecond)
}

func (driver *Driver) wake(wait time.Duration) error {
	err := driver.serialPort.ResetInputBuffer()
	if err != nil {
		return err
	}
	err = driver.send("")
	if err != nil {
		return err
	}
	err = driver.serialPort.Drain()
	if err != nil {
		return err
	}
	err = driver.serialPort.SetReadTimeout(wait)
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port, t time.Duration) {
		_ = serialPort.SetReadTimeout(t)
	}(driver.serialPort, 5*time.Second)
	var byteBuffer [1]byte
	for {
		n, err := driver.serialPort.Read(byteBuffer[:])
		if err != nil {
			if errors.Is(err, syscall.EINTR) {
				continue
			} else {
				return err
			}
		}
		if (n == 0) || (byteBuffer[0] == '!') {
			break
		}
	}
	return nil
}

func (driver *Driver) SetPasskey(level int) error {
	return driver.setIntProperty("Passkey", level)
}

func (driver *Driver) SetMode(mode string) error {
	return driver.setStringProperty("Mode", mode)
}

func (driver *Driver) SetEnablePolledMode(enable bool) error {
	return driver.setBoolProperty("Enable Polled Mode", enable)
}

func (driver *Driver) SetEnableText(enable bool) error {
	return driver.setBoolProperty("Enable Text", enable)
}

func (driver *Driver) SetEnableDecimalformat(enable bool) error {
	return driver.setBoolProperty("Enable Decimalformat", enable)
}

func (driver *Driver) Start() error {
	_, err := driver.request("Start")
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) Stop() error {
	_, err := driver.request("Stop")
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) Save() error {
	_, err := driver.request("Save")
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) Reset() error {
	_, err := driver.request("Reset")
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) DoSample() (*Measurement, error) {
	lines, err := driver.request("Do Sample")
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("no response")
	}

	line := lines[0]
	parts := strings.Split(line, "\t")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid response")
	}
	if parts[0] != "MEASUREMENT" {
		return nil, fmt.Errorf("invalid response")
	}

	parts = parts[3:]
	var quantities []Quantity
	for i := 0; i < len(parts); i += 2 {
		name := parts[i]
		if strings.HasPrefix(name, "*") {
			continue
		}
		p := strings.LastIndex(name, "[")
		if (p < 0) || !strings.HasSuffix(name, "]") {
			continue
		}
		units := name[p+1 : len(name)-1]
		name = name[:p]
		valueString := parts[i+1]
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			continue
		}
		quantities = append(quantities, Quantity{
			Name:  name,
			Value: value,
			Units: units,
		})
	}

	return &Measurement{
		Quantities: quantities,
	}, nil
}

func (driver *Driver) setStringProperty(name string, value string) error {
	_, err := driver.request(fmt.Sprintf("Set %s(%s)", name, value))
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) setIntProperty(name string, value int) error {
	_, err := driver.request(fmt.Sprintf("Set %s(%d)", name, value))
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) setBoolProperty(name string, value bool) error {
	var valueString string
	if value {
		valueString = "Yes"
	} else {
		valueString = "No"
	}
	_, err := driver.request(fmt.Sprintf("Set %s(%s)", name, valueString))
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) request(req string) ([]string, error) {
	err := driver.serialPort.ResetInputBuffer()
	if err != nil {
		return nil, err
	}
	err = driver.send(req)
	if err != nil {
		return nil, err
	}
	err = driver.serialPort.Drain()
	if err != nil {
		return nil, err
	}
	return driver.readResponse()
}

func (driver *Driver) send(text string) error {
	slog.Debug(">>>",
		slog.String("text", text),
	)
	_, err := driver.serialPort.Write([]byte(text + lineTerminator))
	if err != nil {
		return err
	}
	return nil
}

func (driver *Driver) readResponse() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(driver.serialPort)
	for {
		for scanner.Scan() {
			line := stripXonXoff(scanner.Text())
			slog.Debug("<<<",
				slog.String("text", line),
			)
			if line == "#" {
				return lines, nil
			} else if strings.HasPrefix(line, "*") {
				parts := strings.SplitN(line, "\t", 3)
				if len(parts) == 3 {
					if (parts[0] == "*") && (parts[1] == "ERROR") {
						return lines, fmt.Errorf("%s", parts[2])
					}
				}
			}
			lines = append(lines, line)
		}
		if !errors.Is(scanner.Err(), syscall.EINTR) {
			return lines, scanner.Err()
		}
	}
}

func stripXonXoff(line string) string {
	input := []byte(line)
	var output []byte
	for _, b := range input {
		if (b != 0x13) && (b != 0x11) {
			output = append(output, b)
		}
	}
	return string(output)
}

type Measurement struct {
	Quantities []Quantity
}

type Quantity struct {
	Name  string
	Value float64
	Units string
}

func (q Quantity) String() string {
	return fmt.Sprintf("%s: %ef %s", q.Name, q.Value, q.Units)
}
