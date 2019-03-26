package governor

import (
	"fmt"
	"./Config"
)


func SetLights(){

}

func SpamGlobalState(){ //Update and broadcast latest globalState from one peer
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

func UpdateGlobalState(){
	for{
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		Update.c <- GlobalState
		select{
			case Update <- StateUpdate.c//Arriving from NumThree
			//clear hallmatorders and Queue for each elev by checking floor for both elevators, delete UP and DOWN
			//compare, OR global hallmat
				//if 1 when initially 0, distribute order to best elev by also setting that queue=1 for that elev in own globalState

			case Id <- lostElev.c
			//for req in elev[ID] distribute orders
			//delete ID from map in global state
			case <- Button //sent when ordefinished or new button
			//Update hallmat
			//distribute order to best elev by also setting that req=1 for that elev in own globalState
			//if internal -- <- ping fsm
			case globalState[Id] = <- Elev
			//clear hallmatorders, and queue for other elevs 
		}
	}
}

func NetworkState(){
	seen map[int]time
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
