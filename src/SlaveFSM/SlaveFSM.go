package SlaveFSM

import "../elevio"
import "../test"
import "time"
import "fmt"

// TO DO:
// - tell master that order is finished, not necessary, 
// master could check this from pos update and targetFloor(the command it has sent to the slave) (***)


type state int

const (
	UNKNOWN state = 1 + iota
	IDLE
	DOOR_OPEN
	MOVING
)

func FSM(command chan int) {

	// Channels in
	currentFloor := make(chan int, 1)
	lightsFromMaster := make(chan [4][3] int, 12)
	buttonFromIo := make(chan elevio.ButtonEvent)

	// Channels out
	posUpdate := make(chan test.ElevPos, 2)
	cmdFinished := make(chan int, 1)
	buttonToComm := make(chan elevio.ButtonEvent, 10)

	go elevio.PollFloorSensor(currentFloor)
	go test.Routine(lightsFromMaster)
	go handleio(lightsFromMaster, buttonFromIo, buttonToComm)
	go test.Testfunc(posUpdate, buttonToComm)
	go elevio.PollButtons(buttonFromIo)

	state := UNKNOWN
	elevio.SetMotorDirection(elevio.MD_Up)
	doorTimer := time.NewTimer(0)
	floor := -1
	targetFloor := -1
	elevPosition := test.ElevPos{Floor: -1, Id: 0}								// sette egen ip her selv??

	for {
	
		fmt.Println("state: ", state)
		select {
		case targetFloor = <-command:
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
		case floor = <-currentFloor:
			elevPosition.Floor = floor
			posUpdate <- elevPosition
			elevio.SetFloorIndicator(floor)
			switch state {
			case UNKNOWN:
				elevio.SetMotorDirection(elevio.MD_Stop)
				state = IDLE
			case IDLE:
			case DOOR_OPEN:
			case MOVING:
				if floor == targetFloor {
					cmdFinished <- floor
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					doorTimer = time.NewTimer(3 * time.Second)
					state = DOOR_OPEN
				}

			}
		case <-doorTimer.C:
			switch state {
			case IDLE:
			case DOOR_OPEN:
				elevio.SetDoorOpenLamp(false)
				if floor == targetFloor {
					state = IDLE
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


func handleio(lightsFromMaster chan [4][3]int, buttonFromIo chan elevio.ButtonEvent, buttonToComm chan elevio.ButtonEvent) {
	buttonMat := [4][3]int {}
	boolValue := true
	buttonEvent := elevio.ButtonEvent {}
	for {
		select{
		case buttonMat = <- lightsFromMaster:
			for btnType := 0; btnType < 3; btnType++ {
					for floorIndex := 0;  floorIndex < 4; floorIndex++ {
						if !((btnType == 1 && floorIndex == 0) || (btnType == 0 && floorIndex == 3)) {
							if  buttonMat[floorIndex][btnType] == 0{
								boolValue = false
							} else {
								boolValue = true
							}
							elevio.SetButtonLamp(elevio.ButtonType(btnType), floorIndex, boolValue)
						}
					}
			}
		case buttonEvent = <- buttonFromIo:
			fmt.Println("before updating channel")
			buttonToComm <- buttonEvent
		}
	}
}
