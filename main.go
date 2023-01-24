package main

import (
	"fmt"
	"github.com/matiux/memo/domain/aggregate"
)

func main() {

	id1 := aggregate.NewUUIDv4()
	id2 := aggregate.NewUUIDv4()

	checkIDs(id1, id2)
}

func checkIDs(id1, id2 aggregate.EntityId) {
	fmt.Printf("%v\n", id1.Equals(id2))
	fmt.Printf("%T\n", id1)
	fmt.Printf("Value id1: %v\n", id1)
	fmt.Printf("Value id2: %v\n", id2)
}
