package event

import (
	"fmt"
	"sync"
	"testing"
)

var eventDispatcher *Dispatcher = New()

type user struct {
	id   int
	name string
	age  int
}

type ExampleEvent struct {
	data string
}

type SecondEvent struct {
	data *user
}

type ExampleEventListener struct {
	BaseListener
}

func (u ExampleEventListener) Listen() []interface{} {
	return []interface{}{ExampleEvent{}}
}

func (u ExampleEventListener) Handle(e interface{}) error {
	exampleEvent := e.(ExampleEvent)
	fmt.Printf(" processing handle,data: %#v\n", exampleEvent.data)
	return nil
}

type SecondEventListener struct {
	BaseListener
}

func (s SecondEventListener) Listen() []interface{} {
	return []interface{}{SecondEvent{}}
}

func (s SecondEventListener) Handle(e interface{}) error {
	switch event := e.(type) {
	case ExampleEvent:
		fmt.Printf("[ExampleEvent] processing handle,data: %#v\n", event.data)
	case SecondEvent:
		fmt.Printf("[SecondEvent] processing handle,data: %#v\n", event.data)
		event.data.age++
	case *SecondEvent:
		fmt.Printf("[*SecondEvent] processing handle,data: %#v\n", event.data)
		event.data.age++
	}

	return nil
}

type SecondEventListener2 struct {
}

func (s SecondEventListener2) Listen() []interface{} {
	return []interface{}{SecondEvent{}}
}

func (s SecondEventListener2) Handle(e interface{}) error {
	switch event := e.(type) {
	case ExampleEvent:
		fmt.Printf("[ExampleEvent2] processing handle,data: %#v\n", event.data)
	case SecondEvent:
		fmt.Printf("[SecondEvent2] processing handle,data: %#v\n", event.data)
		event.data.age++
	case *SecondEvent:
		fmt.Printf("[*SecondEvent2] processing handle,data: %#v\n", event.data)
		event.data.age++
	}

	return nil
}

func (s SecondEventListener2) Priority() int {
	return HighPriority
}

func TestEventDispatcher_Subscribe(t *testing.T) {
	if eventDispatcher.Subscribe(ExampleEventListener{}) != nil {
		t.Fail()
	}
	if eventDispatcher.Subscribe(SecondEventListener{}) != nil {
		t.Fail()
	}

	if eventDispatcher.Subscribe(SecondEventListener2{}) != nil {
		t.Fail()
	}
}

func TestEventDispatcher_Publish(t *testing.T) {
	var w sync.WaitGroup
	w.Add(3)
	go func() {
		defer w.Done()
		eventDispatcher.Publish(ExampleEvent{data: "ExampleEvent.data values"})
	}()
	go func(userData *user) {
		defer w.Done()
		eventDispatcher.Publish(SecondEvent{data: userData})
		eventDispatcher.Publish(user{})
	}(&user{1, "SecondEvent.data", 1})
	go func(userData *user) {
		defer w.Done()
		eventDispatcher.Publish(SecondEvent{data: userData})
	}(&user{2, "SecondEvent.data", 1})
	w.Wait()
}
