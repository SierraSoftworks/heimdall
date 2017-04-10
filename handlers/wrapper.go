package handlers

type Wrapper struct {
	Handler Handler
	Filter  map[string]interface{}
}
