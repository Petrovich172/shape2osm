// package main

// import (
// 	"fmt"
// )

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

