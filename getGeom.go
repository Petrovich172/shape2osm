// package main

// import (
// 	"fmt"
// )

// func Sqrt(x float64) (int, int) {
// 	z := 1.0
// 	y := 0.0
// 	a := 0
// 	// for (x / (z * z)) > 1.0 {
// 	for z*z < 25 {
// 		z -= (z*z - x) / (2 * z)
// 		y = (x / (z * z))
// 		fmt.Println(y)
// 		a ++
// 	}
// 	return int(z), int(a)
// }

// func main() {
// 	fmt.Println(Sqrt(81))
// }

// package main

// import (
// 	"fmt"
// //	"runtime"
// )

// func switcher_checker(os string) string {
// 	switch os {
// 	case "darwin":
// 		os = "OS X."
// 	case "linux":
// 		os = "Linux."
// 	}
// 	return os
// }

// func main() {
// 	fmt.Print("Go runs on ")
// 	fmt.Printf(switcher_checker("darwin"))
// 	// fmt.Printf("%s.", os)
// }

package main

import (
"fmt"
"encoding/hex"
)


func main() {
	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	i = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i
	fmt.Println(*p)

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j

	 stringBytes, _ := hex.DecodeString("0100000B002900030001FA1A00389405A0979800F6234E0F0202101700F6234E0FE1644B9EB5C0983503008037E39E0004FF00004149")
  fmt.Println("Test HEX string", stringBytes)
}

