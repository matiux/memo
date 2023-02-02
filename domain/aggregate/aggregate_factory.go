package aggregate

import (
	"reflect"
)

type AggregateFactory interface {
	create(aggregateClass reflect.Type, domainEventStream DomainEventStream) (Root, error)
}

type PublicConstructorAggregateFactory struct {
}

func (pc *PublicConstructorAggregateFactory) create(aggregateClass reflect.Type, domainEventStream DomainEventStream) (Root, error) {
	v := reflect.New(aggregateClass)
	aggregateInstance := v.Elem().Interface().(Root)

	es := EventSourcedAggregateRoot{}
	if err := es.InitializeState(domainEventStream, aggregateInstance); err != nil {
		return nil, err
	}

	return aggregateInstance, nil
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
