package stopper

// The synchronous sender implementation. Useful for if you need your receivers to receive in a very specific order.
type SyncSender struct {
	names    []string
	stoppers []*Receiver

	onStop OnStop
}

func (c *SyncSender) Register(name string) *Receiver {
	r := &Receiver{
		C:    make(chan bool),
		done: make(chan bool),
	}

	c.names = append(c.names, name)
	c.stoppers = append(c.stoppers, r)

	return r
}

func (c *SyncSender) Close() {
	for i := 0; i < len(c.stoppers); i++ {
		r := c.stoppers[i]

		close(r.C)
		<-r.done

		if c.onStop != nil {
			c.onStop(c.names[i])
		}
	}
}

func NewSync(onClose OnStop) *AsyncSender {
	return &AsyncSender{
		stoppers: map[string]*Receiver{},
		onStop:   onClose,
	}
}
