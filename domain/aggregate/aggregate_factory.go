package aggregate

type AggregateFactory interface {
	create(aggregateClass Root, domainEventStream DomainEventStream) (Root, error)
}

type PublicConstructorAggregateFactory struct {
}

func (pc *PublicConstructorAggregateFactory) create(aggregateClass Root, domainEventStream DomainEventStream) (Root, error) {

	if err := aggregateClass.InitializeState(domainEventStream, aggregateClass); err != nil {
		return nil, err
	}

	return aggregateClass, nil
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
