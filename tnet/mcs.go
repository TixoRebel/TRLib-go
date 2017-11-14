package tnet

import (
	"net"
	"sync"
	"encoding/binary"
	"time"
	"trlib/tmem"
)

type MultiChannelStream struct {
	buffer [255]*tmem.ExpandingBuffer
	channels [255]net.Conn
	con net.Conn
	lock sync.Mutex
	readErr error
	running bool
}

type multiChannel struct {
	mcs *MultiChannelStream
	channel byte
}

func (m *MultiChannelStream) getBuffer(c byte) *tmem.ExpandingBuffer {
	if m.buffer[c] == nil {
		m.buffer[c] = tmem.NewExpandingBuffer()
	}
	return m.buffer[c]
}

func (m *MultiChannelStream) GetChannel(c byte) net.Conn {
	if m.channels[c] == nil {
		mc := new(multiChannel)
		mc.channel = c
		mc.mcs = m
		m.channels[c] = mc
	}
	return m.channels[c]
}

func (m *MultiChannelStream) Start() {
	m.running = true
	defer func() { m.con.Close(); m.running = false }()
	
	for {
		channel := make([]byte, 1, 1)
		if i, err := m.con.Read(channel); err != nil && i == 0 {
			return
		}
		size := make([]byte, 2, 2)
		m.con.Read(size)
		nsize := binary.BigEndian.Uint16(size)
		data := make([]byte, nsize, nsize)
		m.con.Read(data)
		m.getBuffer(channel[0]).Write(data)
	}
}

func NewMultiChannelStream(c net.Conn) (m *MultiChannelStream) {
	m = new(MultiChannelStream)
	m.running = false
	m.con = c
	return
}

func (c *multiChannel) Read(b []byte) (int, error) {
	return c.mcs.getBuffer(c.channel).ReadAll(b), c.mcs.readErr
}

func (c *multiChannel) Write(b []byte) (written int, err error) {
	c.mcs.lock.Lock()
	defer c.mcs.lock.Unlock()
	
	toWrite := len(b)
	
	for written < toWrite {
		var i int
		i, err = c.mcs.con.Write([]byte { c.channel })
		if err != nil {
			return
		}
		
		shouldWrite := toWrite - written
		if shouldWrite > 65535 {
			shouldWrite = 65535
		}
		sizeWrite := make([]byte, 2, 2)
		binary.BigEndian.PutUint16(sizeWrite, uint16(shouldWrite))
		i, err = c.mcs.con.Write(sizeWrite)
		if err != nil {
			return
		}
		
		i, err = c.mcs.con.Write(b[written:written + shouldWrite])
		written += i
		if err != nil {
			return
		}
	}
	
	return
}

func (c *multiChannel) Close() error {
	c.mcs.buffer[c.channel] = nil
	return nil
}

func (c *multiChannel) LocalAddr() net.Addr {
	return c.mcs.con.LocalAddr()
}

func (c *multiChannel) RemoteAddr() net.Addr {
	return c.mcs.con.RemoteAddr()
}

func (c *multiChannel) SetDeadline(t time.Time) error {
	return c.mcs.con.SetDeadline(t)
}

func (c *multiChannel) SetReadDeadline(t time.Time) error {
	return c.mcs.con.SetReadDeadline(t)
}

func (c *multiChannel) SetWriteDeadline(t time.Time) error {
	return c.mcs.con.SetWriteDeadline(t)
}
