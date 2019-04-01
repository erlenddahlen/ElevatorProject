package FSM

import (
	"time"
	"../DataStructures"
	"../elevio"
)

func FSM(chMan DataStructures.ManagerChannels, chFSM DataStructures.FSMChannels, id string, GState DataStructures.GlobalState) {
	doorTimerDone := time.NewTimer(0)
	doorTimerDone.Stop()

	elev := GState.Map[id]
	elevio.SetMotorDirection(elev.Dir)

	go handleIo(chMan, chFSM)

	for {
		chFSM.LocalStateUpdate <- elev

		select {
		case update := <-chFSM.PingFromGov:
			elev.Queue = update.Map[id].Queue

			switch elev.State {
				case DataStructures.Idle:
					elev.Dir = GetNextDir(elev)
					elevio.SetMotorDirection(elev.Dir)

					if elev.Dir != DataStructures.MD_Stop {
						elev.State = DataStructures.Moving
					} else {
						if elev.Queue[elev.Floor][0] || elev.Queue[elev.Floor][1] || elev.Queue[elev.Floor][2] {
							elevio.SetDoorOpenLamp(true)
							doorTimerDone = time.NewTimer(3 * time.Second)
							elev.Queue[elev.Floor][DataStructures.BT_HallUp] = false
							elev.Queue[elev.Floor][DataStructures.BT_HallDown] = false
							elev.Queue[elev.Floor][DataStructures.BT_Cab] = false
							elev.State = DataStructures.DoorOpen
							chFSM.LocalStateUpdate <- elev
						} else {
							elev.State = DataStructures.Idle
						}
					}

				case DataStructures.DoorOpen:
					if elev.Queue[elev.Floor][0] || elev.Queue[elev.Floor][1] || elev.Queue[elev.Floor][2] {
						elevio.SetDoorOpenLamp(true)
						doorTimerDone = time.NewTimer(3 * time.Second)
						elev.Queue[elev.Floor][DataStructures.BT_HallUp] = false
						elev.Queue[elev.Floor][DataStructures.BT_HallDown] = false
						elev.Queue[elev.Floor][DataStructures.BT_Cab] = false
						elev.State = DataStructures.DoorOpen
					}
				}

		case currentFloor := <-chFSM.AtFloor:
			elev.Floor = currentFloor
			elevio.SetFloorIndicator(currentFloor)

			switch elev.State {
			case DataStructures.Unknown:
				elevio.SetMotorDirection(elevio.MD_Stop)
				elev.State = DataStructures.Idle

			case DataStructures.Moving:
				if ShouldStop(elev) {
					elevio.SetMotorDirection(DataStructures.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					doorTimerDone = time.NewTimer(3 * time.Second)
					elev.Queue[elev.Floor][DataStructures.BT_HallUp] = false
					elev.Queue[elev.Floor][DataStructures.BT_HallDown] = false
					elev.Queue[elev.Floor][DataStructures.BT_Cab] = false
					elev.State = DataStructures.DoorOpen
				}
			}

		case <-doorTimerDone.C:
			elevio.SetDoorOpenLamp(false)
			elev.Dir = GetNextDir(elev)
			elevio.SetMotorDirection(elev.Dir)
			if elev.Dir != DataStructures.MD_Stop {
				elev.State = DataStructures.Moving
			} else {
				elev.State = DataStructures.Idle
			}
		}
	}
}

func handleIo(chMan DataStructures.ManagerChannels, chFSM DataStructures.FSMChannels) {
	go elevio.PollButtons(chFSM.ButtonPushed)
	go elevio.PollFloorSensor(chFSM.AtFloor)
	for {
		select {
		case button := <-chFSM.ButtonPushed:
			if button.Button < 2 {
				chMan.AddHallOrder <- button
			} else {
				chFSM.AddCabOrderGov <- button.Floor
			}
		}
	}
}
