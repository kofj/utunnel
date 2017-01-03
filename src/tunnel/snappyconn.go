package tunnel

import (
	"github.com/golang/snappy"
	"fmt"
	"io"
	"net"
	"compress/flate"
)

// Wrap wraps a connection and adds snappy compression on reading and writing.
func Wrap(wrapped net.Conn) net.Conn {
	cc := &CompressConn{Conn: wrapped}
	
	r := io.Reader(cc.Conn)
	r = flate.NewReader(r)
	cc.r = r
	
	w := io.Writer(cc.Conn)
	zw, err := flate.NewWriter(w, flate.DefaultCompression)
	if err != nil {
		panic(fmt.Sprintf("BUG: flate.NewWriter(%d) returned non-nil err: %s", flate.DefaultCompression, err))
	}
	w = &writeFlusher{w: zw}

	cc.w = w

	//return cc
	return &Snappyconn{wrapped, snappy.NewReader(wrapped), snappy.NewBufferedWriter(wrapped)}
}

type writeFlusher struct {
	w *flate.Writer
}

func (wf *writeFlusher) Write(p []byte) (int, error) {
	n, err := wf.w.Write(p)
	if err != nil {
		return n, err
	}
	if err := wf.w.Flush(); err != nil {
		return 0, err
	}
	return n, nil
}

type CompressConn struct {
	net.Conn
	r io.Reader
	w io.Writer
}

func (c *CompressConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *CompressConn) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
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
	//n, err = c.w.Write(b)
	//if err == nil {
	//	err = c.w.Flush()
	//}
	//return n, err
	return c.w.Write(b)
}