package main
import (
	"fmt"
	"time"
)

func main(){
	list := []int{123,453,6,67,34,6576,0,45,245,3,387523,67,2,8756,1233,33,4,53,75,8,5}
	fmt.Println("原始切片：",list)
	//快排
	quickSort(list,0,len(list)-1)

}

func quickSort(nums []int, left, right int){
	if left < right {
		pos := portion(nums, left, right)
		quickSort(nums,left,pos-1)
		quickSort(nums,pos+1, right)
	}
}

func portion(nums []int,left,right int) int{
	sentry := nums[left]	//最左边的为key
	for left < right{
		for left < right && nums[right] >= sentry {
			right--
		}

		nums[left] = nums[right]

		for left < right && nums[left] <= sentry{
			left++
		}
		nums[right] = nums[left]

	} 
	nums[left] = sentry

	return left		//已排好序的元素下标
}
