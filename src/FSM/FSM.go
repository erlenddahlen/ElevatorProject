package FSM

import (
	"fmt"
	"time"
	"../DataStructures"
	"../elevio"
)

func handleIo(chGov DataStructures.GovernorChannels, chFSM DataStructures.FSMChannels) {
	go elevio.PollButtons(chFSM.ButtonPushed)
	go elevio.PollFloorSensor(chFSM.CurrentFloor)
	for {
		select {
		case button := <-chFSM.ButtonPushed:
			fmt.Println("Button pushed: ", button)
			if button.Button < 2 {
				chGov.AddHallOrder <- button
				//fmt.Println("sent hall to Gov")
			} else {
				//chFSM.AddCabOrder <- button.Floor
				chFSM.AddCabOrderGov <- button.Floor
				//fmt.Println("sent cab to Gov")
			}
		}
	}

}


func FSM(chGov DataStructures.GovernorChannels, chFSM DataStructures.FSMChannels, id string, GState DataStructures.GlobalState) {

	// Init doorTimer, MotorDirection and peer
	doorTimerDone := time.NewTimer(0)
	doorTimerDone.Stop()
	go handleIo(chGov, chFSM)

	elevio.SetMotorDirection(GState.Map[id].Dir)
	peer := GState.Map[id]

	for {

		chFSM.LocalStateUpdate <- peer
		//fmt.Println("Queue in FSM: ", peer.Queue)
		// fmt.Println("STATE: ", peer.State, "DIR: ", peer.Dir, "FLOOR: ", peer.Floor)
		//fmt.Println("QUEUE in FSM: ", peer.Queue)
		//Lights(GState, peer, id)
		select {
		// case cabOrderFloor:= <- chFSM.AddCabOrder:
		//     peer.Queue[cabOrderFloor][2] = true
		//

		case update := <-chFSM.PingFromGov:
			peer.Queue = update.Map[id].Queue
			//Lights(update, peer, id)
			//fmt.Println("CASE A")
			switch peer.State {
			case DataStructures.Idle:
				peer.Dir = GetNextDir(peer)

				//fmt.Println("Next dir: ", peer.Dir)
				elevio.SetMotorDirection(peer.Dir)
				if peer.Dir != DataStructures.MD_Stop {
					peer.State = DataStructures.Moving
					//c.LocalStateUpdate <- peer
				} else {
					if peer.Queue[peer.Floor][0] || peer.Queue[peer.Floor][1] || peer.Queue[peer.Floor][2] {
						elevio.SetDoorOpenLamp(true)
						doorTimerDone = time.NewTimer(3 * time.Second)
						// Slette ordre fra queue
						peer.Queue[peer.Floor][DataStructures.BT_HallUp] = false
						peer.Queue[peer.Floor][DataStructures.BT_HallDown] = false
						peer.Queue[peer.Floor][DataStructures.BT_Cab] = false
						peer.State = DataStructures.DoorOpen
						//c.LocalStateUpdate <- peer
						chFSM.LocalStateUpdate <- peer
					} else {
						peer.State = DataStructures.Idle
					}
				}
			case DataStructures.Moving:
			case DataStructures.DoorOpen:
				if peer.Queue[peer.Floor][0] || peer.Queue[peer.Floor][1] || peer.Queue[peer.Floor][2] {
					elevio.SetDoorOpenLamp(true)
					doorTimerDone = time.NewTimer(3 * time.Second)
					// Slette ordre fra queue
					peer.Queue[peer.Floor][DataStructures.BT_HallUp] = false
					peer.Queue[peer.Floor][DataStructures.BT_HallDown] = false
					peer.Queue[peer.Floor][DataStructures.BT_Cab] = false

					peer.State = DataStructures.DoorOpen
					//c.LocalStateUpdate <- peer
				}
			}

		case currentFloor := <-chFSM.CurrentFloor:
			//fmt.Println("CASE B")
			peer.Floor = currentFloor
			elevio.SetFloorIndicator(currentFloor)
			switch peer.State {
			case DataStructures.Unknown:
				//fmt.Println("case Unknown")
				elevio.SetMotorDirection(elevio.MD_Stop)
				peer.State = DataStructures.Idle

				//fmt.Println("State i unknown: ", peer.State)
			case DataStructures.Moving:
				if ShouldStop(peer) {
					elevio.SetMotorDirection(DataStructures.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					doorTimerDone = time.NewTimer(3 * time.Second)
					// Slette ordre fra queue
					peer.Queue[peer.Floor][DataStructures.BT_HallUp] = false
					peer.Queue[peer.Floor][DataStructures.BT_HallDown] = false
					peer.Queue[peer.Floor][DataStructures.BT_Cab] = false
					peer.State = DataStructures.DoorOpen
				}
				//c.LocalStateUpdate <- peer
			}

			//else if floor == 0 {
			//elevio.SetMotorDirection(elevio.MD_Up)
			//} else if floor == 4 {
			//elevio.SetMotorDirection(elevio.MD_Down)
			//}
		case <-doorTimerDone.C:
			//fmt.Println("CASE C")
			elevio.SetDoorOpenLamp(false)
			peer.Dir = GetNextDir(peer)
			elevio.SetMotorDirection(peer.Dir)
			//fmt.Println("Doortimer done found next dir: ", peer.Dir)
			if peer.Dir != DataStructures.MD_Stop {
				peer.State = DataStructures.Moving
			} else {
				peer.State = DataStructures.Idle
			}
			//c.LocalStateUpdate <- peer
		}
	}

}
