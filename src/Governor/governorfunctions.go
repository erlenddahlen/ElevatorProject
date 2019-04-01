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
	for _,v:= range keys{
		keyString:= strconv.Itoa(v)
		times[i] = TimeToServeOrder(elevators[keyString], NewOrder)
		if times[i] <= times[indexBestElevator]{
			indexBestElevator = i
			keyBestElevator = keyString
			}
		i++
	}
    return keyBestElevator
}

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

func Lights(GState DataStructures.GlobalState, peer DataStructures.Elev, id string) {
	for floor := 0; floor < DataStructures.NumFloors; floor++ {
		if peer.Queue[floor][2] {
			elevio.SetButtonLamp(DataStructures.BT_Cab, floor, true)
		} else {
			elevio.SetButtonLamp(DataStructures.BT_Cab, floor, false)
		}
		if GState.HallRequests[floor][0] {
			elevio.SetButtonLamp(DataStructures.BT_HallUp, floor, true)
		} else {
			elevio.SetButtonLamp(DataStructures.BT_HallUp, floor, false)
		}
		if GState.HallRequests[floor][1] {
			elevio.SetButtonLamp(DataStructures.BT_HallDown, floor, true)
		} else {
			elevio.SetButtonLamp(DataStructures.BT_HallDown, floor, false)
		}
	}

}
