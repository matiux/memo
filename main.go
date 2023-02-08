package main

import (
	"fmt"
	"reflect"
)

type DomainEventStream []string

type EntityRoot struct {
}

func (e *EntityRoot) InitializeState(stream DomainEventStream, aggregate Root) error {
	for _, event := range stream {
		if err := aggregate.Apply(event); err != nil {
			return err
		}
	}

	return nil
}

type Root interface {
	Apply(event string) (err error)
	InitializeState(stream DomainEventStream, aggregate Root) error
}

type Memo struct {
	name    string
	surname string
	EntityRoot
}

func (m *Memo) Apply(event string) (err error) {

	switch event {
	case "Matteo":
		m.name = event
	case "Galacci":
		m.surname = event
	}

	return nil
}

func create(aggregateClass reflect.Value, domainEventStream DomainEventStream) Root {

	inputs := make([]reflect.Value, 2)
	inputs[0] = reflect.ValueOf(domainEventStream)
	//inputs[1] = reflect.ValueOf(aggregateClass.Elem())
	inputs[1] = aggregateClass.Elem().Addr().Convert(reflect.TypeOf((*Root)(nil)).Elem())

	aggregateClass.MethodByName("InitializeState").Call(inputs)

	return aggregateClass.Interface().(Root)

}

func main() {
	root := create(reflect.ValueOf(&Memo{}), DomainEventStream{"Matteo", "Galacci"})
	memo := root.(*Memo)
	fmt.Println(memo)
}

//type MyClass struct {
//	Name string
//}
//
//func create(t reflect.Type) {
//	v := reflect.New(t)
//	v.Elem().Field(0).SetString("John")
//	myClassInstance := v.Elem().Interface()
//
//	fmt.Println(myClassInstance)
//	fmt.Printf("%T\n", myClassInstance)
//	fmt.Printf("%v\n", t)
//	fmt.Printf("%v\n", t)
//}
//
//func main() {
//	t := reflect.TypeOf(MyClass{})
//	create(t)
//
//}
