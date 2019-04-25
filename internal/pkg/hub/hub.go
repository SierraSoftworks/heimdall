package hub

type Hub interface {
	Notify(interface{})
	Subscribe(Subscriber)
	Unsubscribe(Subscriber)
}

type Subscriber interface {
	Receive(interface{})
}
