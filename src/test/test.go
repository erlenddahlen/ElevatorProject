package test
import "fmt"
import "../elevio"


type ElevPos struct {
	Floor  int
	Id     int
}


func Testfunc(update chan ElevPos, buttonFromSlave chan elevio.ButtonEvent) {
	fmt.Println("init test")
	for {
		v:= <- 	update
		w:= <- buttonFromSlave
		fmt.Println("In test func, printing ElevPos update", v, "buttonFromSlave: ", w)


	}
}




func Routine(toSlave chan<- [4][3] int){
	testMatrix:=[4][3]int{}
	testMatrix[1][1] = 1
	testMatrix[2][2] = 1
	testMatrix[0][0] = 1
	toSlave <- testMatrix
}
