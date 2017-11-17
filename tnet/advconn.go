package tnet

import (
	"net"
	"encoding/binary"
	"math"
)

type AdvConn struct {
	net.Conn
}

func (a *AdvConn) WriteByte(b byte) (e error) {
	_, e = a.Write([]byte {b})
	return
}

func (a *AdvConn) ReadByte() (b byte, e error) {
	var buf [1]byte
	_, e = a.Read(buf[:])
	b = buf[0]

	return
}

func (a *AdvConn) WriteUInt8(i uint8) error {
	return a.WriteByte(i)
}

func (a *AdvConn) ReadUInt8() (uint8, error) {
	return a.ReadByte()
}

func (a *AdvConn) WriteInt8(i int8) error {
	return a.WriteByte(uint8(i))
}

func (a *AdvConn) ReadInt8() (int8, error) {
	b, e := a.ReadByte()
	return int8(b), e
}

func (a *AdvConn) WriteUInt16(i uint16) (int, error) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], i)
	return a.Write(buf[:])
}

func (a *AdvConn) ReadUInt16() (uint16, error) {
	var buf [2]byte
	i, e := a.Read(buf[:])
	if i != 2 || e != nil {
		return 0, e
	}
	return binary.BigEndian.Uint16(buf[:]), nil
}

func (a *AdvConn) WriteInt16(i int16) (int, error) {
	return a.WriteUInt16(uint16(i))
}

func (a *AdvConn) ReadInt16() (int16, error) {
	b, e := a.ReadUInt16()
	return int16(b), e
}

func (a *AdvConn) WriteUInt32(i uint32) (int, error) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], i)
	return a.Write(buf[:])
}

func (a *AdvConn) ReadUInt32() (uint32, error) {
	var buf [4]byte
	i, e := a.Read(buf[:])
	if i != 4 || e != nil {
		return 0, e
	}
	return binary.BigEndian.Uint32(buf[:]), nil
}

func (a *AdvConn) WriteInt32(i int32) (int, error) {
	return a.WriteUInt32(uint32(i))
}

func (a *AdvConn) ReadInt32() (int32, error) {
	b, e := a.ReadUInt32()
	return int32(b), e
}

func (a *AdvConn) WriteUInt64(i uint64) (int, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return a.Write(buf[:])
}

func (a *AdvConn) ReadUInt64() (uint64, error) {
	var buf [8]byte
	i, e := a.Read(buf[:])
	if i != 8 || e != nil {
		return 0, e
	}
	return binary.BigEndian.Uint64(buf[:]), nil
}

func (a *AdvConn) WriteInt64(i int64) (int, error) {
	return a.WriteUInt64(uint64(i))
}

func (a *AdvConn) ReadInt64() (int64, error) {
	b, e := a.ReadUInt64()
	return int64(b), e
}

func (a *AdvConn) WriteFloat32(f float32) (int, error) {
	return a.WriteUInt32(math.Float32bits(f))
}

func (a *AdvConn) ReadFloat32() (float32, error) {
	b, e := a.ReadUInt32()
	return math.Float32frombits(b), e
}

func (a *AdvConn) WriteFloat64(f float64) (int, error) {
	return a.WriteUInt64(math.Float64bits(f))
}

func (a *AdvConn) ReadFloat64() (float64, error) {
	b, e := a.ReadUInt64()
	return math.Float64frombits(b), e
}

func (a *AdvConn) WriteRune(r rune) (int, error) {
	return a.WriteUInt32(uint32(r))
}

func (a *AdvConn) ReadRune() (rune, error) {
	b, e := a.ReadUInt32()
	return rune(b), e
}

func (a *AdvConn) WriteNum(n uint64) (int, error) {
	var buf [10]byte
	i := 0
	for n > 0x7F {
		buf[i] = uint8(n) | 0x80
		n >>= 7
		i++
	}
	buf[i] = uint8(n) & 0x7F
	_, err := a.Write(buf[:i + 1])
	return i + 1, err
}

func (a *AdvConn) ReadNum() (uint64, error) {
	var n uint64 = 0
	for i := byte(0); ; i++ {
		b, e := a.ReadByte()
		if e != nil {
			return 0, e
		}
		n |= uint64(b & 0x7F) << (7 * i)
		if b & 0x80 == 0 {
			return n, nil
		}
	}
}

func (a *AdvConn) WriteString(str string) (int, error) {
	buf := []byte(str)
	i, err := a.WriteNum(uint64(len(buf)))
	if err != nil {
		return i, err
	}
	j, err := a.Write(buf)
	return i + j, err
}

func (a *AdvConn) ReadString() (string, error) {
	n, err := a.ReadNum()
	if err != nil {
		return "", err
	}
	buf := make([]byte, n)
	_, err = a.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}