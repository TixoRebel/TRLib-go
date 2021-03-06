package tmem

import (
	"sync"
)

type node struct {
	data []byte
	next *node
}

type ExpandingBuffer struct {
	head *node
	tail *node
	nodeSize int
	writeLock* sync.Mutex
	readLock* sync.Mutex
	wrote* sync.Cond
	length int
}

func NewExpandingBuffer() (e *ExpandingBuffer) {
	e = new(ExpandingBuffer)
	e.writeLock = new(sync.Mutex)
	e.readLock = new(sync.Mutex)
	e.wrote = sync.NewCond(e.writeLock)
	e.length = 0
	return
}

func (e *ExpandingBuffer) Write(data []byte) {
	e.writeLock.Lock()
	defer e.writeLock.Unlock()
	if e.head == nil {
		e.head = e.makeNode()
		e.tail = e.head
	}

	for len(data) > 0 {
		if cap(e.tail.data) == len(e.tail.data) {
			e.tail.next = e.makeNode()
			e.tail = e.tail.next
		}

		i := copy(e.tail.data[len(e.tail.data):cap(e.tail.data)], data)
		data = data[i:]
		e.length += i
		e.tail.data = e.tail.data[:i + len(e.tail.data)]
	}

	e.wrote.Signal()
}

func (e *ExpandingBuffer) WriteDirect(length int, read func([]byte) (int, error)) (int, error) {
	e.writeLock.Lock()
	defer e.writeLock.Unlock()
	if e.head == nil {
		e.head = e.makeNode()
		e.tail = e.head
	}

	wrote := 0

	for wrote < length {
		if cap(e.tail.data) == len(e.tail.data) {
			e.tail.next = e.makeNode()
			e.tail = e.tail.next
		}

		sz := len(e.tail.data) + length - wrote
		if sz > cap(e.tail.data) {
			sz = cap(e.tail.data)
		}
		i, err := read(e.tail.data[len(e.tail.data):sz])

		wrote += i
		e.length += i

		e.tail.data = e.tail.data[:i + len(e.tail.data)]

		if err != nil {
			return wrote, nil
		}
	}

	e.wrote.Signal()
	return wrote, nil
}

func (e *ExpandingBuffer) Read(data []byte) (read int) {
	e.readLock.Lock()
	defer e.readLock.Unlock()
	e.writeLock.Lock()
	defer e.writeLock.Unlock()

	if e.head == nil {
		e.head = e.makeNode()
		e.tail = e.head
	}

	read = 0
	for read < cap(data) {
		if len(e.head.data) == 0 {
			if e.head.next == nil {
				e.wrote.Wait()
			} else {
				e.head = e.head.next
			}
		}
		i := copy(data[read:], e.head.data)
		e.head.data = e.head.data[i:]
		read += i
		e.length -= i
	}

	return
}

func (e *ExpandingBuffer) makeNode() (n *node) {
	if e.nodeSize == 0 { e.nodeSize = 100 }
	n = new(node)
	n.data = make([]byte, e.nodeSize)
	n.data = n.data[0:0]
	
	return
}

func (e *ExpandingBuffer) SetNodeSize(size int) {
	if size >= 0 {
		e.nodeSize = size
	}
}

func (e *ExpandingBuffer) GetNodeSize() int {
	return e.nodeSize
}

func (e *ExpandingBuffer) GetLength() int {
	return e.length
}