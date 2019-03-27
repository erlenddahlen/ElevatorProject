package governor

// func ChooseElevator(elevators map[int]Elevator, NewOrder elevio.ButtonEvent) Elevator {
// 	times := make([]int, len(elevators))
// 	indexBestElevator := 0
// 	keyBestElevator := ""
//
// 	//trying to convert to map
// 	i := 0
// 	for key := range elevators {
// 		times[i] = TimeToServeOrder(elevators[key], NewOrder)
// 		if times[i] < times[indexBestElevator] {
// 			indexBestElevator = i
// 			keyBestElevator = key
// 		}
// 		i++
// 	}
//
// 	return elevators[keyBestElevator]
// }

/*
func AssignOrder(elevator Elevator, Order elevio.ButtonEvent) {
	//Assign the order to the chosen elevator
}
*/

// func TimeToServeOrder(e Elevator, button elevio.ButtonEvent) int {
// 	tempElevator := e
// 	tempElevator.Queue[button.Floor][button.Button] = true
//
// 	timeUsed := 0
//
// 	switch tempElevator.State {
// 	case Idle:
// 		tempElevator.Direction = fsmFunc.ElevChooseDirection(tempElevator)
// 		if tempElevator.Direction == DIR_STOP {
// 			return timeUsed
// 		}
// 	case Moving:
// 		timeUsed += TRAVEL_TIME / 2
// 		tempElevator.Floor += int(tempElevator.Direction)
// 	case DoorOpen:
// 		timeUsed += DOOR_OPEN_TIME / 2
// 	}
//
// 	for {
// 		if fsmFunc.ElevShouldStop(tempElevator) {
// 			if tempElevator.Floor == button.Floor {
// 				//We have "arrived" at destination
// 				return timeUsed
// 			}
// 			tempElevator.Queue[tempElevator.Floor][elevio.BT_Cab] = false
// 			if tempElevator.Direction == DIR_UP || !fsmFunc.RequestAbove(tempElevator) {
// 				tempElevator.Queue[tempElevator.Floor][elevio.BT_HallUp] = false
// 			}
// 			if tempElevator.Direction == DIR_DOWN || !fsmFunc.RequestBelow(tempElevator) {
// 				tempElevator.Queue[tempElevator.Floor][elevio.BT_HallDown] = false
// 			}
// 			timeUsed += DOOR_OPEN_TIME
// 			tempElevator.Direction = fsmFunc.ElevChooseDirection(tempElevator)
// 		}
// 		tempElevator.Floor += int(tempElevator.Direction)
// 		timeUsed += TRAVEL_TIME
// 	}
// }
