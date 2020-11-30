package znet

import "github.com/lumigogogo/zinx/ziface"

// Router 封装一次请求执行方法(业务方实现)
type Router struct {
}

// PreHandle 请求前执行
func (r *Router) PreHandle(request ziface.IRequest) {

}

// Handle 具体业务
func (r *Router) Handle(request ziface.IRequest) {

}

// AfterHandle 请求后执行
func (r *Router) AfterHandle(request ziface.IRequest) {

}
