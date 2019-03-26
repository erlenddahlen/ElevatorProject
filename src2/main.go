package main

import (
    "fmt"
    "./Governor"
    "./Config"
)

func main(){

    Hall:= [[0,0],[1,0],[0,1],[0,0]]

    Cab1:= [0,0,1,1]
    Cab2:= [1,0,0,0]

    temp1 := Elev{IDLE, DirUp, 1, Cab1}
    temp2 := Elev{MOVING, DirDown, 2, Cab2}

    GState := globalState{}

    GState.Map[1] = temp1


    for{}
}
