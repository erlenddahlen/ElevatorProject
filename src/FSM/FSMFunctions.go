package FSM

import (
	"../DataStructures"
)

func getNextDir(elev DataStructures.Elev) DataStructures.MotorDirection {

	if elev.Dir == DataStructures.MotorDirStop {

		if orderAbove(elev) {
			return DataStructures.MotorDirUp
		} else if orderBelow(elev) {
			return DataStructures.MotorDirDown
		} else {
			return DataStructures.MotorDirStop
		}
	} else if elev.Dir == DataStructures.MotorDirUp {
		if orderAbove(elev) {
			return DataStructures.MotorDirUp
		} else if orderBelow(elev) {
			return DataStructures.MotorDirDown
		} else {
			return DataStructures.MotorDirStop
		}
	} else {
		if orderBelow(elev) {
			return DataStructures.MotorDirDown
		} else if orderAbove(elev) {
			return DataStructures.MotorDirUp
		} else {
			return DataStructures.MotorDirStop
		}
	}
	return DataStructures.MotorDirStop
}

func orderAbove(elev DataStructures.Elev) bool {
	for floor := elev.Floor + 1; floor < DataStructures.NumFloors; floor++ {
		for button := 0; button < DataStructures.NumButtons; button++ {
			if elev.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func orderBelow(elev DataStructures.Elev) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for button := 0; button < DataStructures.NumButtons; button++ {
			if elev.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func shouldStop(elev DataStructures.Elev) bool {
	switch elev.Dir {
	case DataStructures.MotorDirUp:
		return (elev.Queue[elev.Floor][DataStructures.HallUp] || elev.Queue[elev.Floor][DataStructures.Cab] || !orderAbove(elev))

	case DataStructures.MotorDirDown:
		return (elev.Queue[elev.Floor][DataStructures.HallDown] || elev.Queue[elev.Floor][DataStructures.Cab] || !orderBelow(elev))

	case DataStructures.MotorDirStop:

	default:
	}
	return false
}

func SetHallAndCabLights(GState DataStructures.GlobalState, elev DataStructures.Elev, id string) {
	for floor := 0; floor < DataStructures.NumFloors; floor++ {
		if elev.Queue[floor][2] {
			elevio.SetButtonLamp(DataStructures.Cab, floor, true)
		} else {
			elevio.SetButtonLamp(DataStructures.Cab, floor, false)
		}
		if GState.HallRequests[floor][0] {
			elevio.SetButtonLamp(DataStructures.HallUp, floor, true)
		} else {
			elevio.SetButtonLamp(DataStructures.HallUp, floor, false)
		}
		if GState.HallRequests[floor][1] {
			elevio.SetButtonLamp(DataStructures.HallDown, floor, true)
		} else {
			elevio.SetButtonLamp(DataStructures.HallDown, floor, false)
		}
	}
}
