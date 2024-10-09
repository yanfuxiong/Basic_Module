package main
import (
	"fmt"
	""
)

func main{
	list := []int{11,13,15,20,24,25,28,32,38,39,41,46,49,66,73}
	fmt.Println("二分查找28："， binSearch(list,28))

}

func binSearch(nums []int, k int) int{
	left, right := 0, len(nums) - 1
	mark := -1 

	for left <= right {
		mid := (left+right) /2
		if nums[mid] == k{
			return mid
		}else if nums[mid]> k{
			right = mid-1
		}else if nums[mid] <k{
			left = mid+1
		}
	}

	return -1
}