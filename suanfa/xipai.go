package (
	"math/rand"
	"fmt"
	"time"
)



func XiPai(){
	Pai := make([]int, 54)
	for i:=0; i<54; i++ {
		fmt.Scanln(@Pai[i])
		// Pai[i] = i
	}

	//随机数
	rand.Seed(time.Now().UnixNano())
	//rand.Intn(100)
	for i:=53; i>0; i--{
		position := rand.Intn(i + 1)
		Pai[position],Pai[i] = Pai[i], Pai[position]
	}

	fmt.Println(Pai)

}




