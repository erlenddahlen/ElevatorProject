package governor

import (
	"fmt"
	"./Config"
	"./elevio"
	"./PeerFSM"
)


func SetLights(){

}

//SPØRSMÅL:
// Nå fjernes ordre fra queue og HallRequests på to måter. Enten at lokalPeer sletter når den har utført ordre i en etasje
// eller at governor sjekker om en heis står i state DOOR_OPEN og sletter ordren derfra. Kan det skape trøbbel?
// Må tenke litt gjennom rekkefølgen vi gjør ting: som når vi distribuerer/kalkulerer optimal heis for ordre og når vi sender ut info


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
			case newOrder:= <- elevio.ButtonPushed //sent when ordefinished or new button
			//Update hallmat
			//distribute order to best elev by also setting that req=1 for that elev in own globalState
			//if internal -- <- ping fsm


			// 1. Sjekke om cab eller hall-Button. Oppdatere cab eller oppdatere global HallRequests OG sende ut
			// denne oppdaterte Gstaten FØR optimalisering? Slik at alle optimaliserer likt? Eller er ikke dette nødvendig
			if newOrder.Button == 2{
				Gstate.Map["minId"].Queue[newOrder.Floor][newOrder.Button] = true
				PeerFSM.pingFromGov <- 1
			} else {
				GState.HallRequests[newOrder.Floor][newOrder.Button] = true
				// send ut ny update til alle?

				// 2. Kjøre en optimize/cost og distribuere ordren, dvs legge den inn i den med lavest kost sin lokale queue
				keyBestElevator:= ChooseElevator(Gstate.Map, newOrder)
				Gstate.Map[keyBestElevator].Queue[newOrder.Floor][newOrder.Button]
				if keyBestElevator == "minId" {
					PeerFSM.pingFromGov <- 1
				}
			}

			// 3. Gi en ping til lokal om at det har skjedd en update (gjør dette underveis)

			case peerUpdate:= <- PeerFSM.LocalStateUpdate 		//comes from local, Husk: dette er en Elev struct

			// 1. Bytte ut hele min egen elevState/peer med denne updaten
			Gstate.Map["minId"] = peerUpdate

			// 2. Oppdatere global HallRequests
			for floor:= 0;  floor < Config.NumFloors; ++ {
				for button:=0;  < Config.NumButtons-1; ++ {
					if peerUpdate.Queue[floor][button]==false &&  Gstate.HallRequests[floor][button]==true{
						Gstate.HallRequests[floor][button]==false
					}
				}
			}
			// 3. Oppdatere alles lokale HallRequests/Queue
			// Oppdatere innebærer å fjerne ordre som er fjernet fra lokalHall i update
			for id := 0;  id < Config.NumOfElev; id++ {
				HallRequestsOfId:= Gstate.Map[id]
				for floor:= 0;  floor < Config.NumFloors; ++ {
					for button:=0;  < Config.NumButtons-1; ++ {
						if peerUpdate.Queue[floor][button]==false &&  HallRequestsOfId.Queue[floor][button]==true{
							HallRequestsOfId.Queue[floor][button]==false
						}
					}
				}
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
