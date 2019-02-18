package main

import "./elevio"
import "time"
import "fmt"

// TO DO: Turn all the lights on, fix problem line 97, targetfloor/queue,
// ny terminal der simulator ligger: ./SimElevatorServer
func main() {
	elevio.Init("localhost:15657", 4)
	command := make(chan int)
	currentFloor := make(chan int, 1)
	button := make(chan elevio.ButtonEvent)

	go elevio.PollFloorSensor(currentFloor)
	go elevio.PollButtons(button)
	go run(currentFloor, command)
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

func run(currentFloor chan int, command chan int) {

	state := UNKNOWN
	elevio.SetMotorDirection(elevio.MD_Up)
	doorTimer := time.NewTimer(0)
	floor := -1
	targetFloor := -1
	for {
		fmt.Println(state)
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
				if targetFloor != floor {
					fmt.Println("targetfloor ikke lik floor, sender kommando på nytt")
					state = IDLE
					//command <- targetFloor // kan kopiere inn if betingelser inn her for å løse problem
				} else {
					//SendFinished()
					state = IDLE
				}
			case MOVING:
			}
		}
	}
}
