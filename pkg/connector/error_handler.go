package connector

type ErrorHandler interface {
	Handle(err error)
}
