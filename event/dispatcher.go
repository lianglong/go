package event

import (
	"container/heap"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Dispatcher struct {
	handlers map[string]*ListenerQueue
	lock     *sync.RWMutex // a lock for the map
}

func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]*ListenerQueue),
		lock:     &sync.RWMutex{},
	}
}

//Publish event
func (ed *Dispatcher) Publish(e interface{}) {
	eventName := getStructFullName(e)
	if handlers := ed.getHandlers(eventName); len(handlers) > 0 {
		for _, item := range handlers {
			item.listener.Handle(e)
		}
	}
}

func (ed *Dispatcher) getHandlers(eventName string) ListenerQueue {
	ed.lock.RLock()
	defer ed.lock.RUnlock()
	if lq, ok := ed.handlers[eventName]; ok && lq != nil {
		return *lq
	}
	return ListenerQueue{}
}

//Subscribe event
func (ed *Dispatcher) Subscribe(listener Listener) error {
	ed.lock.Lock()
	defer ed.lock.Unlock()

	if listener.Listen() == nil || len(listener.Listen()) < 1 {
		return errors.New("incorrect listen parameters")
	}
	for _, event := range listener.Listen() {
		eventName := getStructFullName(event)
		if _, ok := ed.handlers[eventName]; !ok {
			ed.handlers[eventName] = &ListenerQueue{}
		}
		heap.Push(ed.handlers[eventName], &ListenerItem{listener: listener, priority: listener.Priority()})
	}

	return nil
}

func getStructFullName(s interface{}) string {
	ref := reflect.TypeOf(s)
	switch ref.Kind() {
	case reflect.Ptr:
		ref = ref.Elem()
		if ref.Kind() != reflect.Struct {
			panic("parameter type error")
		}
	case reflect.Struct:
	default:
		panic("parameter type error")
	}
	return fmt.Sprintf("%s.%s", ref.PkgPath(), ref.Name())
}
