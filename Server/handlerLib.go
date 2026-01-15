package main


// Define a handler type
type MessageHandler func(message message) error

// Create a Handler that holds handlers
type Handler struct {
	handlers    map[string]MessageHandler
	subHandlers map[string]*Handler
}

func NewHandler() *Handler {
	return &Handler{
		handlers:    make(map[string]MessageHandler),
		subHandlers: make(map[string]*Handler),
	}
}

func (r *Handler) Register(path string, handler MessageHandler) {
	r.handlers[path] = handler
}

func (r *Handler) RegisterSubHandler(path string, subHandler *Handler) {
	r.subHandlers[path] = subHandler
}

func (r *Handler) Route(message message) error {
	if len(message.Headers) == 0 {
		return nil
	}
	key := message.Headers[0]

	// Check for direct handler
	if handler, exists := r.handlers[key]; exists {
		return handler(message)
	}

	// Check for sub-Handler
	if subHandler, exists := r.subHandlers[key]; exists {
		// Remove the first header and route to sub-handler
		message.Headers = message.Headers[1:]
		return subHandler.Route(message)
	}

	return nil
}
