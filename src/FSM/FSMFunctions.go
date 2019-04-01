package FSM

import (
	"../DataStructures"
)

func GetNextDir(elev DataStructures.Elev) DataStructures.MotorDirection {

	if elev.Dir == DataStructures.MD_Stop {

		if OrderAbove(elev) {
			return DataStructures.MD_Up
		} else if OrderBelow(elev) {
			return DataStructures.MD_Down
		} else {
			return DataStructures.MD_Stop
		}
	} else if elev.Dir == DataStructures.MD_Up {
		if OrderAbove(elev) {
			return DataStructures.MD_Up
		} else if OrderBelow(elev) {
			return DataStructures.MD_Down
		} else {
			return DataStructures.MD_Stop
		}
	} else {
		if OrderBelow(elev) {
			return DataStructures.MD_Down
		} else if OrderAbove(elev) {
			return DataStructures.MD_Up
		} else {
			return DataStructures.MD_Stop
		}
	}
	return DataStructures.MD_Stop
}

func OrderAbove(elev DataStructures.Elev) bool {
	for floor := elev.Floor + 1; floor < DataStructures.NumFloors; floor++ {
		for button := 0; button < DataStructures.NumButtons; button++ {
			if elev.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func OrderBelow(elev DataStructures.Elev) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for button := 0; button < DataStructures.NumButtons; button++ {
			if elev.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func ShouldStop(elev DataStructures.Elev) bool {
	switch elev.Dir {
	case DataStructures.MD_Up:
		return (elev.Queue[elev.Floor][DataStructures.BT_HallUp] || elev.Queue[elev.Floor][DataStructures.BT_Cab] || !OrderAbove(elev))

	case DataStructures.MD_Down:
		return (elev.Queue[elev.Floor][DataStructures.BT_HallDown] || elev.Queue[elev.Floor][DataStructures.BT_Cab] || !OrderBelow(elev))

	case DataStructures.MD_Stop:

	default:
	}
	return false
}
