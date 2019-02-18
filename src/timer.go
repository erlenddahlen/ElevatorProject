package main

import (
	"fmt"
	"time"
)

func main() {
	doorTimer := time.NewTimer(0)
	//doorTimer = time.NewTimer(3 * time.Second)
	hei := <-doorTimer.C

	fmt.Println(hei)
}
