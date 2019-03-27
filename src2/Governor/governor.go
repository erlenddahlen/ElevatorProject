package governor

import (
	"fmt"
	"./Config"
	"./elevio"
	"./PeerFSM"
)


func SetLights(){

}

//Huske på: Vi MÅ sørge for at alle har den nyeste informasjonen når det gjøres en distribuering!


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
			case Update <- StateUpdate //StateUpdate from other Peer
			OutsideElev:= Update.Map[Update.Id]

			// 1. Set our info about the OutsideElev to the update from OutsideElev
			Gstate.Map[Update.Id] = OutsideElev
			// 2. If OutsideElev in state = DOOR_OPEN, remove orders in that floor
			if OutsideElev.State == DOOR_OPEN{
					for button:= 0;  button < Config.NumButtons - 1 ; ++ {
						GState.HallRequests[peerUpdate.Floor][button] = false
					}
			}
			// 3. Fylle lokal HallRequests med de evt nye knappetrykkene fra oppdateringen
			// Ta vare på disse knappetrykkene så de kan bli distribuert
			var newButton Config.ButtonEvent
			for floor:= 0;  floor < Config.NumFloors; ++ {
				for button:= 0;  button < Config.NumButtons; ++ {
					if Update.HallRequests[floor][button] == true {
						GState.HallRequests[floor][button] = true
						newButton.Floor = floor
						newButton.Button = button
						// Distribuere denne nye knappen, evt legge den til i en liste over knapper som må distrubueres
					}
				}
			}


			case Id <- lostElev.c
			//for req in elev[ID] distribute orders
			//delete ID from map in global state
			case newOrder:= <- elevio.ButtonPushed //sent when ordefinished or new button
			// 1. Sjekke om cab eller hall-Button. Oppdatere cab eller oppdatere global HallRequests OG sende ut
			// denne oppdaterte Gstaten FØR optimalisering? Slik at alle optimaliserer likt? Eller er ikke dette nødvendig
			if newOrder.Button == 2{
				Gstate.Map["minId"].Queue[newOrder.Floor][newOrder.Button] = true
				PeerFSM.pingFromGov <- 1
			} else {
				GState.HallRequests[newOrder.Floor][newOrder.Button] = true
				// send ut ny update til alle?

				// 2. Kjøre en optimize/cost og distribuere ordren til oss selv om vi får den
				keyBestElevator:= ChooseElevator(Gstate.Map, newOrder)
				if keyBestElevator == "minId" {
					Gstate.Map["minId"].Queue[newOrder.Floor][newOrder.Button]
					PeerFSM.pingFromGov <- 1
				}
			}

			// 3. Gi en ping til lokal om at det har skjedd en update (gjør dette underveis)

			case peerUpdate:= <- PeerFSM.LocalStateUpdate 		//comes from local, Husk: dette er en Elev struct
			// 1. Bytte ut hele min egen elevState/peer med denne updaten
			Gstate.Map["minId"] = peerUpdate

			// 2. Staten som oppdateringen har(alstå FSM-state) er DOOR_OPEN, fjern hall i denne etasjen
			if peerUpdate.State == Config.DOOR_OPEN {
				for button:= 0;  button < Config.NumButtons - 1 ; ++ {
					GState.HallRequests[peerUpdate.Floor][button] = false
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
