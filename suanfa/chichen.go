



import (
 	"fmt"	
)


func Isvalid()bool
{
	var n int64

	fmt.Scanln(&n)

	ValueLis := []int{7,13,29,20,36,42,49}

	//三重循环法
	for i:=0;i<=n; i+=7{
		for j:=0; j<=n; j+=13{
			for k:=0; k<=n; k+=29{
				if i+j+k == n {
					return true
				}
			}
		}
	}


	return false 

}