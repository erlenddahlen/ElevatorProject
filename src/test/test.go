package test


func testfunc() {
		toSlave := make(chan [4][3] int, 12)
		go Routine(toSlave)

		for {

		}
	}



func Routine(toSlave chan<- [4][3] int){
	testMatrix:=[4][3]int{}
	testMatrix[1][1] = 1
	testMatrix[2][2] = 1
	testMatrix[0][0] = 1
	toSlave <- testMatrix
}
