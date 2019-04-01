package Governor
//
import (
	"../DataStructures"
	"../PeerFSM"
	 "strconv"
	 "fmt"
	 "sort"
)

func ChooseElevator(elevators map[string]DataStructures.Elev, NewOrder DataStructures.ButtonEvent, GState DataStructures.GlobalState) string{ //Should reurn DataStructures.Elev
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

func TimeToServeOrder(e DataStructures.Elev, button DataStructures.ButtonEvent) int {
	tempElevator := e
	tempElevator.Queue[button.Floor][button.Button] = true

	timeUsed := 0

	switch tempElevator.State {
	case DataStructures.Idle:
		tempElevator.Dir = PeerFSM.GetNextDir(tempElevator)
		if tempElevator.Dir == DataStructures.MD_Stop {
			return timeUsed
		}
	case DataStructures.Moving:
		timeUsed += DataStructures.TravelTime / 2
		tempElevator.Floor += int(tempElevator.Dir)
	case DataStructures.DoorOpen:
		timeUsed += DataStructures.DoorOpenTime / 2
	}
    count := 0
	for {
		if PeerFSM.ShouldStop(tempElevator) {
			if tempElevator.Floor == button.Floor {
				//We have "arrived" at destination
				return timeUsed
			}
			tempElevator.Queue[tempElevator.Floor][DataStructures.BT_Cab] = false
			if tempElevator.Dir == DataStructures.MD_Up || !PeerFSM.OrderAbove(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][DataStructures.BT_HallUp] = false
			}
			if tempElevator.Dir == DataStructures.MD_Down || !PeerFSM.OrderBelow(tempElevator) {
				tempElevator.Queue[tempElevator.Floor][DataStructures.BT_HallDown] = false
			}
			timeUsed += DataStructures.DoorOpenTime
			tempElevator.Dir = PeerFSM.GetNextDir(tempElevator)
		}
		tempElevator.Floor += int(tempElevator.Dir)
		timeUsed += DataStructures.TravelTime
        count = count + 1
	}
}
