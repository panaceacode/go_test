package main

import (
	"fmt"
	"go_test/slice_delete"
)

func main() {
	testSlice := []int{
		2,
	}
	//testSlice1 := []int{
	//	2, 0, 2, 4, 0, 8, 1, 4,
	//}
	fmt.Println(testSlice)
	test1, _ := slice_delete.DeleteAt(testSlice, 0)
	fmt.Println(test1)
}
