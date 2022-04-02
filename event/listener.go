package event

type Listener interface {
	Listen() []interface{}
	Handle(interface{}) error
	Priority() int
}

type BaseListener struct {
}

func (l BaseListener) Priority() int {
	return NormalPriority
}

type ListenerItem struct {
	listener Listener
	priority int
}

type ListenerQueue []*ListenerItem

func (lq ListenerQueue) Len() int           { return len(lq) }
func (lq ListenerQueue) Less(i, j int) bool { return lq[i].priority > lq[j].priority }
func (lq ListenerQueue) Swap(i, j int)      { lq[i], lq[j] = lq[j], lq[i] }

func (lq *ListenerQueue) Push(item interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*lq = append(*lq, item.(*ListenerItem))
}

func (lq *ListenerQueue) Pop() interface{} {
	old := *lq
	n := len(old)
	x := old[n-1]
	*lq = old[0 : n-1]
	return x
}
