
func GetNextDir(peer DataStructures.Elev) DataStructures.MotorDirection {

	if peer.Dir == DataStructures.MD_Stop {

		if OrderAbove(peer) {
			return DataStructures.MD_Up
		} else if OrderBelow(peer) {
			return DataStructures.MD_Down
		} else {
			return DataStructures.MD_Stop
		}
	} else if peer.Dir == DataStructures.MD_Up {
		if OrderAbove(peer) {
			return DataStructures.MD_Up
		} else if OrderBelow(peer) {
			return DataStructures.MD_Down
		} else {
			return DataStructures.MD_Stop
		}
	} else {
		if OrderBelow(peer) {
			return DataStructures.MD_Down
		} else if OrderAbove(peer) {
			return DataStructures.MD_Up
		} else {
			return DataStructures.MD_Stop
		}
	}
	return DataStructures.MD_Stop
}

func OrderAbove(peer DataStructures.Elev) bool {
	for floor := peer.Floor + 1; floor < DataStructures.NumFloors; floor++ {
		for button := 0; button < DataStructures.NumButtons; button++ {
			if peer.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func OrderBelow(peer DataStructures.Elev) bool {
	for floor := 0; floor < peer.Floor; floor++ {
		for button := 0; button < DataStructures.NumButtons; button++ {
			if peer.Queue[floor][button] {
				return true
			}
		}
	}
	return false
}

func ShouldStop(peer DataStructures.Elev) bool {
	//fmt.Println("peer.Dir: ", peer.Dir)
	//fmt.Println("peer.Floor: ", peer.Floor)
	switch peer.Dir {
	case DataStructures.MD_Up:
		//fmt.Println("MD_Up: ", peer.Queue[peer.Floor][DataStructures.BT_HallUp], peer.Queue[peer.Floor][DataStructures.BT_Cab], !OrderAbove(peer))
		return (peer.Queue[peer.Floor][DataStructures.BT_HallUp] || peer.Queue[peer.Floor][DataStructures.BT_Cab] || !OrderAbove(peer))
	case DataStructures.MD_Down:
		//fmt.Println("MD_Down: ", peer.Queue[peer.Floor][DataStructures.BT_HallDown],  peer.Queue[peer.Floor][DataStructures.BT_Cab], !OrderBelow(peer))
		return (peer.Queue[peer.Floor][DataStructures.BT_HallDown] || peer.Queue[peer.Floor][DataStructures.BT_Cab] || !OrderBelow(peer))
	case DataStructures.MD_Stop:
		// fmt.Println("MD_Stop: ", peer.Queue[peer.Floor][DataStructures.BT_HallUp ], peer.Queue[peer.Floor][DataStructures.BT_HallDown], peer.Queue[peer.Floor][DataStructures.BT_Cab])
		// if(peer.Queue[peer.Floor][DataStructures.BT_HallUp ] || peer.Queue[peer.Floor][DataStructures.BT_HallDown] || peer.Queue[peer.Floor][DataStructures.BT_Cab]){
		//     return true
		// }
	default:
	}
	return false
}
