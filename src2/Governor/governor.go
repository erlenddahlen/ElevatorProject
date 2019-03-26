package governor

import (
	"fmt"
	"./Config"
	"./elevio"
	"./PeerFSM"
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

func UpdateGlobalState(PeerUpdate chan Config.Elev){
	for{
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		Update.c <- GlobalState
		select{
			case Update <- StateUpdate.c //StateUpdate from other Peer
			OutsideElev:= Update.Map[Update.Id]
			if OutsideElev.State == DOOR_OPEN{
					//Set HallRequests UP and DOWN in that floor to zero
					GState.HallRequests[OutsideElev.Floor][0] = 0
					GState.HallRequests[OutsideElev.Floor][1] = 0
					//Set HallOrders in other elevs queue to zero
					for _,elev := range GState.Map{
						elev.Queue[OutsideElev.Floor][0] = 0
						elev.Queue[OutsideElev.Floor][1] = 0
						}
			}
			//Copy Caborders for OutsideElev to own map
			GState[Update.Id].Cab = Outside.cab // denne må fikses, cab = slice Queue


			//compare, OR global hallmat
				//if 1 when initially 0, distribute order to best elev by also setting that queue=1 for that elev in own globalState

			case Id <- lostElev.c
			//for req in elev[ID] distribute orders
			//delete ID from map in global state
			case <- Button //sent when ordefinished or new button
			//Update hallmat
			//distribute order to best elev by also setting that req=1 for that elev in own globalState
			//if internal -- <- ping fsm
			case globalState[Id] = <- PeerUpdate //from local
			//clear hallmatorders, and queue for other elevs

			// 1. Bytte ut hele min egen elevState/peer med denne updaten
			// 2. Bytte ut alle sin HallRequests
			for id := 0;  id < Config.NumOfElev; id++ {
				GlobalState[id].

			}

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
