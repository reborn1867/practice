package main

import (
	"fmt"
	cases "practice/cases" 
)

func main() {
	fmt.Println("hello world")
	test_nums := []int{1,2,3,5,54}
	target := 8
	res1 := cases.TwoSum(test_nums, target)
	res2 := cases.TwoSumNew(test_nums, target)
	fmt.Printf("result of first solution is %d\n", res1)
	fmt.Printf("result of better solution is %d\n", res2)
}