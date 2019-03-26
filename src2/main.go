package main

import (
    "fmt"
    //"./Governor"
    "./Config"
    "./Peertest"
)

func main(){

    Queue1:= [4][3]bool{{false,false,false},{false,false,false},{false,false,false} ,{true,false,false}}
    Queue2:= [4][3]bool{{false,false,false},{false,false,false},{false,false,false} ,{false,false,false}}

    temp1 := Config.Elev{Config.IDLE, Config.DirStop, 1, Queue1}
    temp2 := Config.Elev{Config.IDLE, Config.DirStop, 2, Queue2}

    GState := Config.GlobalState{}
    GState.Map = make(map[string]Config.Elev)

    GState.Map["1"] = temp1
    GState.Map["2"] = temp2

    NewButton := Config.ButtonEvent{Floor:0, Button: Config.BT_HallUp}

    fmt.Println(ChooseElevator(GState.Map, NewButton))
    //for{}
}

func ChooseElevator(elevators map[string]Config.Elev, NewOrder Config.ButtonEvent) string{ //Should reurn Config.Elev
	times := make([]int, len(elevators))
	indexBestElevator := 0
	keyBestElevator := ""

	//trying to convert to map
	i := 0
	for key, _ := range elevators {
		times[i] = TimeToServeOrder(elevators[key], NewOrder)
		if times[i] <= times[indexBestElevator] {
			indexBestElevator = i
			keyBestElevator = key
        }
		i++
	}

	//return elevators[keyBestElevator]
    return keyBestElevator
}

/*
func AssignOrder(elevator Elevator, Order elevio.ButtonEvent) {
	//Assign the order to the chosen elevator
}
*/

func TimeToServeOrder(e Config.Elev, button Config.ButtonEvent) int {
	tempElevator := e
	tempElevator.Queue[button.Floor][button.Button] = true

	timeUsed := 0

	switch tempElevator.State {
	case Config.IDLE:
		tempElevator.Dir = Peertest.GetNextDir(tempElevator)
		if tempElevator.Dir == Config.DirStop {
			return timeUsed
		}
	case Config.MOVING:
		timeUsed += Config.TRAVEL_TIME / 2
		tempElevator.Floor += int(tempElevator.Dir)
	case Config.DOOR_OPEN:
		timeUsed += Config.DOOR_OPEN_TIME / 2
	}
    count := 0
	for {
		if Peertest.ShouldStop(tempElevator) {
			if tempElevator.Floor == button.Floor {
				//We have "arrived" at destination
				return timeUsed
			}
			tempElevator.Queue[tempElevator.Floor][Config.BT_Cab] = false
			if tempElevator.Dir == Config.DirUp || !Peertest.OrderAbove(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][Config.BT_HallUp] = false
			}
			if tempElevator.Dir == Config.DirDown || !Peertest.OrderBelow(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][Config.BT_HallDown] = false
			}
			timeUsed += Config.DOOR_OPEN_TIME
			tempElevator.Dir = Peertest.GetNextDir(tempElevator)
		}
		tempElevator.Floor += int(tempElevator.Dir)
		timeUsed += Config.TRAVEL_TIME
        count = count + 1
	}
}
