package hub

// NewCallbackSubscriber creates a new subscriber which accepts a
// callback function which is executed every time a message is received.
func NewCallbackSubscriber(cb func(msg interface{})) Subscriber {
	return &callbackHandler{
		callback: cb,
	}
}

type callbackHandler struct {
	callback func(msg interface{})
}

func (h *callbackHandler) Receive(msg interface{}) {
	go h.callback(msg)
}
