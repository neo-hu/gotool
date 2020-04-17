package pool

type Pool interface {
	Get(args interface{}) interface{}
	Put(c interface{})
}

type pool struct {
	pool chan interface{}
	fn   func(args interface{}) interface{}
}

func NewPool(max int, fn func(args interface{}) interface{}) Pool {
	return &pool{
		pool: make(chan interface{}, max),
		fn:   fn,
	}
}

// 把item放回pool， 如果pool满了，则丢弃
func (p *pool) Put(c interface{}) {
	select {
	case p.pool <- c:
	default:
	}
}

// 从pool中获取一个，如果没有调用fn方法创建一个
func (p *pool) Get(args interface{}) interface{} {
	var c interface{}
	select {
	case c = <-p.pool:
	default:
		c = p.fn(args)
	}
	return c
}
