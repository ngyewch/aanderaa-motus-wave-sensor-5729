package aanderaa_motus_wave_sensor_5729

import (
	"fmt"
	"io"
)

var (
	errTimeout = fmt.Errorf("timeout")
)

type timeoutReader struct {
	r io.Reader
}

func newTimeoutReader(r io.Reader) timeoutReader {
	return timeoutReader{
		r: r,
	}
}

func (t timeoutReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if (n == 0) && (err == nil) {
		err = errTimeout
	}
	return
}
