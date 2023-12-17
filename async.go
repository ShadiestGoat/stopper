package stopper

import "sync"

// The async sender implementation. Useful for if you don't care in which order your stuff gets closed.
type AsyncSender struct {
	stoppers map[string]*Receiver
	onStop   OnStop
}

func (c *AsyncSender) Register(name string) *Receiver {
	r := &Receiver{
		C:    make(chan bool),
		done: make(chan bool),
	}

	c.stoppers[name] = r

	return r
}

func (c *AsyncSender) Stop() {
	wg := &sync.WaitGroup{}

	wg.Add(len(c.stoppers))

	for n := range c.stoppers {
		go func(n string) {
			r := c.stoppers[n]

			close(r.C)
			<-r.done

			if c.onStop != nil {
				c.onStop(n)
			}

			wg.Done()
		}(n)
	}

	wg.Wait()
}

func NewAsync(onStop OnStop) *AsyncSender {
	return &AsyncSender{
		stoppers: map[string]*Receiver{},
		onStop:   onStop,
	}
}
