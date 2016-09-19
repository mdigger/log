package log

// Handler is used to handle log events, outputting them to stdio or sending
// them to remote services.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	Handle(*Entry) error
}
