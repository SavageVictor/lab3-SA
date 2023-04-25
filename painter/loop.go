package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver   Receiver
	next       screen.Texture
	prev       screen.Texture
	mq         messageQueue
	operations []Operation // Added a slice to store past operations
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.mq = newMessageQueue()

	go l.processMessages()
}

func (l *Loop) Post(op Operation) {
	l.operations = append(l.operations, op) // Add new operation to the operations slice
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	close(l.mq.queue)
}

func (l *Loop) processMessages() {
	for {
		op := l.mq.pull()
		if op == nil {
			continue
		}

		// Apply all previous operations
		for _, pastOp := range l.operations {
			pastOp.Do(l.next)
		}

		update := op.Do(l.next)
		if update {
			l.Receiver.Update(l.next)
			l.next, l.prev = l.prev, l.next
		}
	}
}

type messageQueue struct {
	queue chan Operation
}

func newMessageQueue() messageQueue {
	return messageQueue{
		queue: make(chan Operation, 100),
	}
}

func (mq *messageQueue) push(op Operation) {
	mq.queue <- op
}

func (mq *messageQueue) pull() Operation {
	return <-mq.queue
}
