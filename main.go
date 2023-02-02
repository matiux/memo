package main

import (
	"fmt"
	"reflect"
)

type MyClass struct {
	Name string
}

func create(t reflect.Type) {
	v := reflect.New(t)
	v.Elem().Field(0).SetString("John")
	myClassInstance := v.Elem().Interface()

	fmt.Println(myClassInstance)
	fmt.Printf("%T\n", myClassInstance)
	fmt.Printf("%v\n", t)
	fmt.Printf("%v\n", t)
}

func main() {
	t := reflect.TypeOf(MyClass{})
	create(t)

}

//func main() {
//
//	id1 := aggregate.NewUUIDv4()
//	id2 := aggregate.NewUUIDv4()
//
//	checkIDs(id1, id2)
//}

//func checkIDs(id1, id2 aggregate.EntityId) {
//	fmt.Printf("%v\n", id1.Equals(id2))
//	fmt.Printf("%T\n", id1)
//	fmt.Printf("Value id1: %v\n", id1)
//	fmt.Printf("Value id2: %v\n", id2)
//}
