package FSM

import (
	"../DataStructures"
	"../elevio"
	"time"
)

func FSM(chMan DataStructures.ManagerChannels, chFSM DataStructures.FSMChannels, id string, GState DataStructures.GlobalState) {
	doorTimerDone := time.NewTimer(0)
	doorTimerDone.Stop()

	elev := GState.Map[id]
	elevio.SetMotorDirection(elev.Dir)

	go handleIo(chMan, chFSM)

	for {
		chFSM.UpdateFromFSM <- elev

		select {
		case update := <-chFSM.UpdateFromManager:
			elev.Queue = update.Map[id].Queue

			switch elev.State {
			case DataStructures.Idle:
				elev.Dir = GetNextDir(elev)
				elevio.SetMotorDirection(elev.Dir)

				if elev.Dir != DataStructures.MotorDirStop { // If  new order is not on the current floor
					elev.State = DataStructures.Moving
				} else {
					//Open door if new order is on current floor
					if elev.Queue[elev.Floor][0] || elev.Queue[elev.Floor][1] || elev.Queue[elev.Floor][2] {
						elevio.SetDoorOpenLamp(true)
						doorTimerDone = time.NewTimer(3 * time.Second)
						elev.Queue[elev.Floor][DataStructures.HallUp] = false
						elev.Queue[elev.Floor][DataStructures.HallDown] = false
						elev.Queue[elev.Floor][DataStructures.Cab] = false
						elev.State = DataStructures.DoorOpen
						chFSM.UpdateFromFSM <- elev
					} else { //No orders in queue
						elev.State = DataStructures.Idle
					}
				}

			case DataStructures.DoorOpen:
				if elev.Queue[elev.Floor][0] || elev.Queue[elev.Floor][1] || elev.Queue[elev.Floor][2] {
					elevio.SetDoorOpenLamp(true)
					doorTimerDone = time.NewTimer(3 * time.Second)
					elev.Queue[elev.Floor][DataStructures.HallUp] = false
					elev.Queue[elev.Floor][DataStructures.HallDown] = false
					elev.Queue[elev.Floor][DataStructures.Cab] = false
					elev.State = DataStructures.DoorOpen
				}
			}

		case currentFloor := <-chFSM.AtFloor:
			elev.Floor = currentFloor
			elevio.SetFloorIndicator(currentFloor)

			switch elev.State {
			case DataStructures.Unknown:
				elevio.SetMotorDirection(DataStructures.MotorDirStop)
				elev.State = DataStructures.Idle

			case DataStructures.Moving:
				if ShouldStop(elev) {
					elevio.SetMotorDirection(DataStructures.MotorDirStop)
					elevio.SetDoorOpenLamp(true)
					doorTimerDone = time.NewTimer(3 * time.Second)
					elev.Queue[elev.Floor][DataStructures.HallUp] = false
					elev.Queue[elev.Floor][DataStructures.HallDown] = false
					elev.Queue[elev.Floor][DataStructures.Cab] = false
					elev.State = DataStructures.DoorOpen
				}
			}

		case <-doorTimerDone.C:
			elevio.SetDoorOpenLamp(false)
			elev.Dir = GetNextDir(elev)
			elevio.SetMotorDirection(elev.Dir)
			if elev.Dir != DataStructures.MotorDirStop {
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
				chFSM.AddCabOrderManager <- button.Floor
			}
		}
	}
}
