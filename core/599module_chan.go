package core

// 阻塞队列
func (fns *InternalFunctionSet) BlockList(args []interface{}) Value {
	queueLen := 128
	if len(args) > 0 {
		if arg, ok := args[0].(int); ok && arg > 0 {
			queueLen = arg
		}
	}
	queue := make(chan interface{}, queueLen)
	obj := &BlockQueue{queue}
	return newClass("BlockQueue", &obj)
}

type BlockQueue struct {
	list chan interface{}
}

func (q *BlockQueue) Get() interface{} {
	return <-q.list
}

func (q *BlockQueue) Set(item interface{}) {
	q.list <- item
}
