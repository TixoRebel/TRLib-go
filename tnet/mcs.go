package tnet

import (
	"net"
	"sync"
	"time"
	"trlib/tmem"
)

type MultiChannelStream struct {
	buffer [255]*tmem.ExpandingBuffer
	channels [255]*AdvConn
	con *AdvConn
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

func (m *MultiChannelStream) GetChannel(c byte) *AdvConn {
	if m.channels[c] == nil {
		mc := new(multiChannel)
		mc.channel = c
		mc.mcs = m
		m.channels[c] = &AdvConn{mc}
	}
	return m.channels[c]
}

func (m *MultiChannelStream) Start() error {
	m.running = true
	defer func() { m.con.Close(); m.running = false }()

	for {
		channel, err := m.con.ReadByte()
		if err != nil {
			return err
		}
		nsize, err := m.con.ReadUInt16()
		if err != nil {
			return err
		}
		data := make([]byte, nsize, nsize)
		_, err = m.con.Read(data)
		if err != nil {
			return err
		}
		m.getBuffer(channel).Write(data)
	}
}

func (m *MultiChannelStream) Close() error {
	return m.con.Close()
}

func NewMultiChannelStream(c net.Conn) (m *MultiChannelStream) {
	m = new(MultiChannelStream)
	m.running = false
	m.con = &AdvConn{c}
	return
}

func (c *multiChannel) Read(b []byte) (int, error) {
	return c.mcs.getBuffer(c.channel).Read(b), c.mcs.readErr
}

func (c *multiChannel) Write(b []byte) (written int, err error) {
	c.mcs.lock.Lock()
	defer c.mcs.lock.Unlock()

	toWrite := len(b)

	for written < toWrite {
		var i int
		err = c.mcs.con.WriteByte(c.channel)
		if err != nil {
			return
		}
		
		shouldWrite := toWrite - written
		if shouldWrite > 0xFFFF {
			shouldWrite = 0xFFFF
		}
		i, err = c.mcs.con.WriteUInt16(uint16(shouldWrite))
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
