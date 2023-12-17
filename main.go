package stopper

type OnStop func(closerName string)

type Sender interface {
	// Register a new receiver
	Register(name string) *Receiver
	// Send a close signal to all the receivers (ie. close everything)
	Stop()
}

// A receiver of close signals. Receive from <- Receiver.C, and don't forget to run CLoseReceiver.Done()!
type Receiver struct {
	C    chan bool
	done chan bool
}

func (c *Receiver) Done() {
	close(c.done)
}

// Create a new sender. onStop is usually used for logging.
// The difference between an async and a sync sender is that sync sender guarantee order, while async runs them all at once
func NewSender(onStop OnStop, async bool) Sender {
	if async {
		return NewAsync(onStop)
	} else {
		return NewSync(onStop)
	}
}
