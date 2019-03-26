package governor

import (
	"fmt"
)

//Create variables that should be used
var timeOut int

type Elev struct {
	State ElevState
	Dir   Direction
	Floor int
	Queue [NumFloors][NumButtons]bool
}

var globalState map


func Govern(){
	var ()
	for{
		select{
			case newLocalOrder
			case CompleteOrder.Floor
		}
	}
}

func SetLights(){

}

func NumOne(){ //Update and broadcast latest globalState from one peer
	latestState //This is the global state
	for{
		select{
			case <- t.c // ticker 
				network.c <- latestState //sending latest state to network
			case newUpdate <- Update.c
				latestState = newUpdate
		}
	}
}

func NumTwo(){
	for{
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		select{
			case Update <- StateUpdate.c//Arriving from NumThree
			case Id <- lostElev.c
			case <- Button //sent when ordefinished or new button
		}
	}
}

func NumThree(){
	seen map[Id]time
	timeOut = 1 second
	for{
		select{
			StateUpdate <- network.c
			map[StateUpdate.Id] = time.Now()
			StateUpdate.c <- StateUpdate

		default: 
			for v, k := range seen: 
				if (t > timeout){ //t = time.Now.Sub(v)
					lostElev.c <- v.Id
					delete(seen, v.Id)
				}
		}

	}
}
