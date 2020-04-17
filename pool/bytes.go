package pool

type LimitedBytes struct {
	pool    chan []byte
	maxSize int
}

// maxSize: 如果maxSize>0, 并且item大于maxSize，则丢弃，不缓存
func NewLimitedBytes(max, maxSize int) *LimitedBytes {
	return &LimitedBytes{
		pool:    make(chan []byte, max),
		maxSize: maxSize,
	}
}

func (l *LimitedBytes) Get(sz int) []byte {
	var c []byte
	select {
	case c = <-l.pool:
	default:
		// todo 如果不存在，创建
		return make([]byte, sz)
	}
	// todo 如果存在，但是缓存的cap小于需要的cap，则创建一个
	if cap(c) < sz {
		return make([]byte, sz)
	}
	// todo 只返回需要的大小
	return c[:sz]
}

func (l *LimitedBytes) Put(c []byte) {
	if cap(c) >= l.maxSize {
		return
	}
	select {
	case l.pool <- c:
	default:
	}
}
