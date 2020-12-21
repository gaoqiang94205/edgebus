package controller

import (
	"fmt"
	"testing"
)

func TestMol(t *testing.T) {
	fmt.Println(majorityElement([]int{3, 2, 3, 2, 1, 3}))
}

func majorityElement(nums []int) []int {
	len := len(nums)
	limit := len / 3
	var result = make([]int, 2)
	var count1, count2 int
	//计数统计阶段
	for _, n := range nums {
		if result[0] == n {
			count1++
			continue
		}
		if result[1] == n {
			count2++
			continue
		}
		if count1 == 0 {
			result[0] = n
			count1++
			continue
		}
		if count2 == 0 {
			result[1] = n
			count2++
			continue
		}

		count1--
		count2--
	}
	//for 循环一次进行校验
	var c1, c2 int
	for _, n := range nums {
		if result[0] == n {
			c1++
			continue
		}
		if result[1] == n {
			c2++
			continue
		}
	}
	r := make([]int, 0)
	if c1 > limit {
		r = append(r, result[0])
	}
	if c2 > limit {
		r = append(r, result[1])
	}
	return r
}
