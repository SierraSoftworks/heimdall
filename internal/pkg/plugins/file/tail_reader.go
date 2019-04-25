package file

import (
	"io"
	"time"
)

// TailReader provides a blocking file reader which
// will continue to attempt reads once it reaches the
// end of a file.
type TailReader struct {
	f     io.Reader
	delay time.Duration
}

// NewTailReader creates a new instance of the TailReader capable
// of blocking reads until more data is available at the end of
// a stream.
func NewTailReader(r io.Reader, interval time.Duration) io.Reader {
	return &TailReader{
		f:     r,
		delay: interval,
	}
}

func (r *TailReader) Read(b []byte) (n int, err error) {
	n, err = r.f.Read(b)

	for n == 0 && err == io.EOF {
		time.Sleep(r.delay)
		n, err = r.f.Read(b)
	}

	return n, err
}
