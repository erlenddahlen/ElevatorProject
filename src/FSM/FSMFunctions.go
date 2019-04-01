package FSM

import (
	"../DataStructures"
)

func GetNextDir(elev DataStructures.Elev) DataStructures.MotorDirection {

	if elev.Dir == DataStructures.MotorDirStop {

		if OrderAbove(elev) {
			return DataStructures.MotorDirUp
		} else if OrderBelow(elev) {
			return DataStructures.MotorDirDown
		} else {
			return DataStructures.MotorDirStop
		}
	} else if elev.Dir == DataStructures.MotorDirUp {
		if OrderAbove(elev) {
			return DataStructures.MotorDirUp
		} else if OrderBelow(elev) {
			return DataStructures.MotorDirDown
		} else {
			return DataStructures.MotorDirStop
		}
	} else {
		if OrderBelow(elev) {
			return DataStructures.MotorDirDown
		} else if OrderAbove(elev) {
			return DataStructures.MotorDirUp
		} else {
			return DataStructures.MotorDirStop
		}
	}
	return DataStructures.MotorDirStop
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
	case DataStructures.MotorDirUp:
		return (elev.Queue[elev.Floor][DataStructures.BT_HallUp] || elev.Queue[elev.Floor][DataStructures.BT_Cab] || !OrderAbove(elev))

	case DataStructures.MotorDirDown:
		return (elev.Queue[elev.Floor][DataStructures.BT_HallDown] || elev.Queue[elev.Floor][DataStructures.BT_Cab] || !OrderBelow(elev))

	case DataStructures.MotorDirStop:

	default:
	}
	return false
}
