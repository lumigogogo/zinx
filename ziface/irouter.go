package ziface

// IRouter ..
type IRouter interface {
	PreHandle(IRequest)
	Handle(IRequest)
	AfterHandle(IRequest)
}
