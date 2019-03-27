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
			lostElev:= Gstate.Map[Id]
			delete(GState, Id)
			newOrder := Config.ButtonEvent{}

			for floor:= 0;  floor < Config.NumFloors; ++ {
				for button:=0;  < Config.NumButtons-1; ++ {
					if lostElev.Queue[floor][button]==true{
						newOrder = Config.ButtonEvent{floor, button}
						keyBestElevator:= ChooseElevator(Gstate.Map, newOrder)
						if keyBestElevator == "minId" {
							Gstate.Map["minId"].Queue[newOrder.Floor][newOrder.Button]
							PeerFSM.pingFromGov <- 1
						}
					}
				}
			}

			
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

			// 2. Staten som oppdateringen har(alstå FSM-state) er DOOR_OPEN, fjern hall i denne etasjen
			if peerUpdate.State == Config.DOOR_OPEN {
				for button:= 0;  button < Config.NumButtons - 1 ; ++ {
					GState.HallRequests[peerUpdate.Floor][button] = false
				}
			}

			// Fra tirsdag:
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