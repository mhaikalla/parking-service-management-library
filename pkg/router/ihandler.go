package router

// IHandleRegister interface to handle route registering
type IHandleRegister interface {
	Handle(method, path string, handler func(interface{}) error)
}

// IHandleRegisterV2 consisting IHandleRegister and add method HandleAuth to handle authenticated route
type IHandleRegisterV2 interface {
	IHandleRegister
	HandleAuth(method, path string, handler func(interface{}) error)
}
