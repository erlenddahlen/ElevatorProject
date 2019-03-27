package governor

import (
	"fmt"
	"../Config"
//	"../elevio"
//	"../PeerFSM"
	"time"
	"../network/bcast"
)


//func SetLights(){} // Do this in slave FSM

//Huske på: Vi MÅ sørge for at alle har den nyeste informasjonen når det gjøres en distribuering!
var latestState Config.GlobalState
var GState Config.GlobalState

func SpamGlobalState(gchan Config.GovernorChannels){ //Update and broadcast latest globalState from one peer
	//This is the global state
	ticker := time.Tick(2500 * time.Millisecond)

	transmitNet := make(chan Config.GlobalState)
	go bcast.Transmitter(16569, transmitNet)

	for{
		select{
			case <- ticker: // ticker
				//fmt.Println(t)
				transmitNet <- latestState //sending latest state to network
				//fmt.Println("Transmitting GlobalState from ID: ", latestState.Id)
			case newUpdate:= <- gchan.InternalState:
				latestState = newUpdate
		}
	}
}

func UpdateGlobalState(PeerUpdate chan Config.Elev, gchan Config.GovernorChannels, FSMchan Config.FSMChannels){
	for{
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		gchan.ExternalState <- GState
		select{
			case Update:=  <- gchan.ExternalState: //StateUpdate from other Peer
			OutsideElev:= Update.Map[Update.Id]

			// 1. Set our info about the OutsideElev to the update from OutsideElev
			GState.Map[Update.Id] = OutsideElev
			// 2. If OutsideElev in state = DOOR_OPEN, remove orders in that floor
			if OutsideElev.State == Config.DOOR_OPEN{
					for button:= 0;  button < Config.NumButtons - 1 ; button++ {
						GState.HallRequests[OutsideElev.Floor][button] = false
					}
			}
			// 3. Fylle lokal HallRequests med de evt nye knappetrykkene fra oppdateringen
			// Ta vare på disse knappetrykkene så de kan bli distribuert
			newOrder := Config.ButtonEvent{}
			for floor:= 0;  floor < Config.NumFloors; floor++ {
				for button:= 0;  button < Config.NumButtons; button++ {
					if Update.HallRequests[floor][button] == true {
						GState.HallRequests[floor][button] = true
						newOrder = Config.ButtonEvent{floor, Config.ButtonType(button)}
						fmt.Println(newOrder)
						// Distribuere denne nye knappen, evt legge den til i en liste over knapper som må distrubueres
					}
				}
			}

		case Id:= <- gchan.LostElev:
			lostElev:= GState.Map[Id]
			delete(GState.Map, Id)
			newOrder := Config.ButtonEvent{}

			for floor:= 0;  floor < Config.NumFloors; floor++ {
				for button:=0;  button < Config.NumButtons-1; button++ {
					if lostElev.Queue[floor][button]==true{
						newOrder = Config.ButtonEvent{floor, Config.ButtonType(button)}
						fmt.Println(newOrder)
						//keyBestElevator:= ChooseElevator(GState.Map, newOrder)
						//if keyBestElevator == GState.Id {
							//GState.Map[GState.Id].Queue[newOrder.Floor][newOrder.Button] What should this do?
							FSMchan.PingFromGov <- 1
						//}
					}
				}
			}


			// case newOrder:= <- elevio.ButtonPushed: //sent when ordefinished or new button
			// // 1. Sjekke om cab eller hall-Button. Oppdatere cab eller oppdatere global HallRequests OG sende ut
			// // denne oppdaterte GStaten FØR optimalisering? Slik at alle optimaliserer likt? Eller er ikke dette nødvendig
			// if newOrder.Button == 2{
			// 	GState.Map["minId"].Queue[newOrder.Floor][newOrder.Button] = true
			// 	FSMchan.pingFromGov <- 1
			// } else {
			// 	GState.HallRequests[newOrder.Floor][newOrder.Button] = true
			// 	// send ut ny update til alle?
			//
			// 	// 2. Kjøre en optimize/cost og distribuere ordren til oss selv om vi får den
			// 	keyBestElevator:= ChooseElevator(GState.Map, newOrder)
			// 	if keyBestElevator == "minId" {
			// 		GState.Map["minId"].Queue[newOrder.Floor][newOrder.Button]
			// 		FSMchan.pingFromGov <- 1
			// 	}
			// }
			//
			// // 3. Gi en ping til lokal om at det har skjedd en update (gjør dette underveis)
			//
			// case peerUpdate:= <- FSMchan.LocalStateUpdate: 		//comes from local, Husk: dette er en Elev struct
			// // 1. Bytte ut hele min egen elevState/peer med denne updaten
			// GState.Map["minId"] = peerUpdate
			//
			// // 2. Staten som oppdateringen har(alstå FSM-state) er DOOR_OPEN, fjern hall i denne etasjen
			// if peerUpdate.State == Config.DOOR_OPEN {
			// 	for button:= 0;  button < Config.NumButtons - 1 ; button++ {
			// 		GState.HallRequests[peerUpdate.Floor][button] = false
			// 	}
			// }
		}
	}
}

func NetworkState(gchan Config.GovernorChannels){
	seen:= make(map[int]time.Time)
	timeOut := Config.TIMEOUT*time.Second

	receiveNet := make(chan Config.GlobalState)
	go bcast.Receiver(16569, receiveNet)

	for{
		select{
		case StateUpdate:= <- receiveNet:
			if latestState.Id != StateUpdate.Id {

			seen[StateUpdate.Id] = time.Now()
			//gchan.ExternalState <- StateUpdate
			//fmt.Println("Received external state from ID:", StateUpdate.Id)
			//fmt.Println(StateUpdate)
			}
		default:
			for k, v := range seen{
				t:= time.Now().Sub(v)
				if (t > timeOut){
					gchan.LostElev <- k
					delete(seen, k)
				}
			}
		}
	}
}
