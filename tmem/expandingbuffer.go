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
	Size int
	writeLock* sync.Mutex
	readLock* sync.Mutex
	wrote* sync.Cond
}

func NewExpandingBuffer() (e *ExpandingBuffer) {
	e = new(ExpandingBuffer)
	e.writeLock = new(sync.Mutex)
	e.readLock = new(sync.Mutex)
	e.wrote = sync.NewCond(e.writeLock)
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
		e.tail.data = e.tail.data[:i + len(e.tail.data)]
	}
	
	e.wrote.Signal()
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
		i := copy(data[read:], e.head.data)
		e.head.data = e.head.data[i:]
		read += i
		if len(e.head.data) == 0 {
			if e.head.next == nil { return }
			e.head = e.head.next
		}
	}
	
	return
}

func (e *ExpandingBuffer) ReadAll(data []byte) (read int) {
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
	}
	
	return
}

func (e *ExpandingBuffer) makeNode() (n *node) {
	if e.Size == 0 { e.Size = 1024 }
	n = new(node)
	n.data = make([]byte, e.Size)
	n.data = n.data[0:0]
	
	return
}