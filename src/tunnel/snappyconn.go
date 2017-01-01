package tunnel

import (
	"github.com/golang/snappy"
	"net"
)

// Wrap wraps a connection and adds snappy compression on reading and writing.
func Wrap(wrapped net.Conn) net.Conn {
	return &Snappyconn{wrapped, snappy.NewReader(wrapped), snappy.NewBufferedWriter(wrapped)}
}

type Snappyconn struct {
	net.Conn
	r *snappy.Reader
	w *snappy.Writer
}

func (c *Snappyconn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *Snappyconn) Write(b []byte) (n int, err error) {
	n, err = c.w.Write(b)
	if err == nil {
		err = c.w.Flush()
	}
	return n, err
}