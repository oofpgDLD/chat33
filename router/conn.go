package router

type Connection interface {
	WriteResponse(interface{})
	Args() interface{}
	Start() error
	Close(...interface{}) error
}
