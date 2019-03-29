package Governor
//
import (
	"../Config"
	"../PeerFSM"
	 "strconv"
	 "fmt"
	 "sort"
)

func ChooseElevator(elevators map[string]Config.Elev, NewOrder Config.ButtonEvent, GState Config.GlobalState) string{ //Should reurn Config.Elev
	times := make([]int, len(elevators))
	indexBestElevator := 0
	keyBestElevator := GState.Id

	//trying to convert to map
	i := 0
	keys := make([]int, 0)
	for key, _ := range elevators {
			keyInt, error := strconv.Atoi(key)
			if error != nil {
					fmt.Println(error.Error())
		}
		keys = append(keys, keyInt)
	}
	sort.Ints(keys)
	//fmt.Println(keys)
	for _,v:= range keys{
		keyString:= strconv.Itoa(v)
		times[i] = TimeToServeOrder(elevators[keyString], NewOrder)
		if times[i] <= times[indexBestElevator]{
			indexBestElevator = i
			keyBestElevator = keyString
			}
		i++
	}


	// for key, _ := range elevators {
	// 	keyInt, error := strconv.Atoi(key)
	// 	if error != nil {
	// 			fmt.Println(error.Error())
	// 	}

	// 	times[i] = TimeToServeOrder(elevators[key], NewOrder)
	// 	if times[i] <= times[indexBestElevator] && keyInt < lowestkey{
	// 		indexBestElevator = i
	// 		keyBestElevator = key
	// 		lowestkey = keyInt
  //       }
	// 	i++
	// }

	//return elevators[keyBestElevator]
	fmt.Println("Choose finished, key: ", keyBestElevator)
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
		tempElevator.Dir = PeerFSM.GetNextDir(tempElevator)
		if tempElevator.Dir == Config.MD_Stop {
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
		if PeerFSM.ShouldStop(tempElevator) {
			if tempElevator.Floor == button.Floor {
				//We have "arrived" at destination
				return timeUsed
			}
			tempElevator.Queue[tempElevator.Floor][Config.BT_Cab] = false
			if tempElevator.Dir == Config.MD_Up || !PeerFSM.OrderAbove(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][Config.BT_HallUp] = false
			}
			if tempElevator.Dir == Config.MD_Down || !PeerFSM.OrderBelow(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][Config.BT_HallDown] = false
			}
			timeUsed += Config.DOOR_OPEN_TIME
			tempElevator.Dir = PeerFSM.GetNextDir(tempElevator)
		}
		tempElevator.Floor += int(tempElevator.Dir)
		timeUsed += Config.TRAVEL_TIME
        count = count + 1
	}
}
