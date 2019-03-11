package main

import "./elevio"
//import "./test"
import "./SlaveFSM"


// ny terminal der simulator ligger: ./SimElevatorServer

func main() {
	elevio.Init("localhost:15657", 4)

	button := make(chan elevio.ButtonEvent)
	command := make(chan int)
	go elevio.PollButtons(button)
	go SlaveFSM.FSM(command)


	for {
		select {
		case btn := <-button:
			command <- btn.Floor
		}
	}
}

