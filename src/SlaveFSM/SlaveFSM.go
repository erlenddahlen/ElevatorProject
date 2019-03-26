package SlaveFSM

import (
	"fmt"
	"strconv"
	"time"

	"../CommTest"
	"../elevio"
)

type state int

const (
	UNKNOWN state = 1 + iota
	IDLE
	DOOR_OPEN
	MOVING
)

func FSM(idString string, chSlave CommTest.SlaveFSMChannels, chComm CommTest.CommunicationChannels, lightsFromMaster chan [4][3]int) {
	id, error := strconv.Atoi(idString)
	if error != nil {
		fmt.Println(error.Error())
	}

	go elevio.PollFloorSensor(chSlave.CurrentFloor)
	go handleio(id, lightsFromMaster, chSlave.ButtonFromIo, chSlave.ButtonPushed)
	go elevio.PollButtons(chSlave.ButtonFromIo)

	state := UNKNOWN
	elevio.SetMotorDirection(elevio.MD_Up)
	doorTimer := time.NewTimer(0)
	floor := -1
	targetFloor := -1
	elevPosition := CommTest.ElevPos{floor, id}
	cmdFinished := CommTest.Cmd{floor, id}

	for {
		//fmt.Println("state: ", state)
		select {
		case targetFloor = <-chComm.CmdElevToFloorToSlave:
			//	fmt.Println("targetfloor", targetFloor)
			fmt.Println("targetfloor:", targetFloor)
			switch state {
			case IDLE:
				if floor == targetFloor {
					elevio.SetDoorOpenLamp(true)
					doorTimer = time.NewTimer(3 * time.Second)
					state = DOOR_OPEN
				} else if floor < targetFloor {
					state = MOVING
					elevio.SetMotorDirection(elevio.MD_Up)
				} else {
					state = MOVING
					elevio.SetMotorDirection(elevio.MD_Down)
				}
			case DOOR_OPEN:
				if floor == targetFloor {
					doorTimer = time.NewTimer(3 * time.Second)

				}
			case MOVING:
				if floor < targetFloor {
					elevio.SetMotorDirection(elevio.MD_Up)
				} else {
					elevio.SetMotorDirection(elevio.MD_Down)
				}
			}
		case floor = <-chSlave.CurrentFloor:
			//fmt.Println("floor in slave", floor)
			//			fmt.Println("Currentfloor Slave: ", floor)
			elevPosition.Floor = floor
			//fmt.Println("Internal PosUpdate", elevPosition)
			chSlave.PosUpdate <- elevPosition
			elevio.SetFloorIndicator(floor)
			switch state {
			case UNKNOWN:
				elevio.SetMotorDirection(elevio.MD_Stop)
				state = IDLE
			case IDLE:
			case DOOR_OPEN:
			case MOVING:
				if floor == targetFloor {
					//cmdFinished.Floor = floor
					//fmt.Println("Internal cmdFinished", cmdFinished)
					//chSlave.CmdFinished <- cmdFinished
					//fmt.Println("Sending CMD case floor")
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					doorTimer = time.NewTimer(3 * time.Second)
					state = DOOR_OPEN
				} else if floor == 0 {
					elevio.SetMotorDirection(elevio.MD_Up)
				} else if floor == 4 {
					elevio.SetMotorDirection(elevio.MD_Down)
				}
			}

		case <-doorTimer.C:
			//			fmt.Println("doortimer")
			switch state {
			case IDLE:
			case DOOR_OPEN:
				elevio.SetDoorOpenLamp(false)
				if floor == targetFloor {
					state = IDLE
					cmdFinished.Floor = floor
					chSlave.CmdFinished <- cmdFinished
					fmt.Println("Sending CMD case doortimer")
					//					fmt.Println("Sending cmd finished from slave")
				} else if floor < targetFloor {
					state = MOVING
					elevio.SetMotorDirection(elevio.MD_Up)
				} else {
					state = MOVING
					elevio.SetMotorDirection(elevio.MD_Down)
				}
			case MOVING:
			}
		}
	}
}

func handleio(id int, lightsFromMaster chan [4][3]int, buttonFromIo chan elevio.ButtonEvent,
	buttonPushed chan CommTest.ButtonPushed) {
	buttonMat := [4][3]int{}
	boolValue := true
	button := CommTest.ButtonPushed{}
	button.Id = id
	for {
		select {
		case buttonMat = <-lightsFromMaster:
			for btnType := 0; btnType < 3; btnType++ {
				for floorIndex := 0; floorIndex < 4; floorIndex++ {
					if !((btnType == 1 && floorIndex == 0) || (btnType == 0 && floorIndex == 3)) {
						if buttonMat[floorIndex][btnType] == 0 {
							boolValue = false
						} else {
							boolValue = true
						}
						elevio.SetButtonLamp(elevio.ButtonType(btnType), floorIndex, boolValue)
					}
				}
			}
		case buttonEvent := <-buttonFromIo:
			button.Button = buttonEvent
			fmt.Println("ButtonPushed from slave", button.Button)
			buttonPushed <- button
		}
	}
}
