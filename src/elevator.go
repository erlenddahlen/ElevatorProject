package main

import "./elevio"
import "./test"
import "time"
import "fmt"

// TO DO:
// - tell master that order is finished
// - tell master which buttons are pushed
// - tell master its position

// ny terminal der simulator ligger: ./SimElevatorServer



func main() {
	elevio.Init("localhost:15657", 4)
	command := make(chan int)
	currentFloor := make(chan int, 1)
	button := make(chan elevio.ButtonEvent)
	//buttonPushed := make(chan elevio.ButtonType)
	lightsFromMaster := make(chan [4][3] int, 12)


	go elevio.PollFloorSensor(currentFloor)
	go elevio.PollButtons(button)
	go test.Routine(lightsFromMaster)

	go elevFSM(currentFloor, command)
	go handleio(lightsFromMaster)
	for {
		select {
		case btn := <-button:
			command <- btn.Floor
		}
	}
}

type state int

const (
	UNKNOWN state = 1 + iota
	IDLE
	DOOR_OPEN
	MOVING
)

func elevFSM(currentFloor chan int, command chan int) {

	state := UNKNOWN
	elevio.SetMotorDirection(elevio.MD_Up)
	doorTimer := time.NewTimer(0)
	floor := -1
	targetFloor := -1
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
			elevio.SetFloorIndicator(floor)
			switch state {
			case UNKNOWN:
				elevio.SetMotorDirection(elevio.MD_Stop)
				state = IDLE
			case IDLE:
			case DOOR_OPEN:
			case MOVING:
				if floor == targetFloor {
					fmt.Println("Stopping at target floor")
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					doorTimer = time.NewTimer(3 * time.Second)
					state = DOOR_OPEN
				}

			}
		case <-doorTimer.C:
			fmt.Println("Timer finished")
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


func handleio(lightsFromMaster chan [4][3]int) {
	buttonMat := [4][3]int {}
	boolValue := true
	select{
	case buttonMat = <- lightsFromMaster:
		for btnType := 0; btnType < 3; btnType++ {
				for floorIndex := 0;  floorIndex < 4; floorIndex++ {
					fmt.Println(btnType, floorIndex)
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
	}
}
