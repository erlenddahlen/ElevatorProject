package test

import (
	"fmt"
)

type Cmd struct {
	Floor  int
	Id     int
}

func main() {
	x := Cmd{1, 1}
	y := Cmd{1, 1}
	c := make(chan Cmd)
	c<-x
	fmt.Println(<-c)
	c<-x
	c<-y
	fmt.Println(<-c)
	fmt.Println(<-c)

}
