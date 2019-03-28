package Governor

import (
	"fmt"
	"../Config"
//	"../elevio"
//	"../PeerFSM"
	"time"
	"../network/bcast"
)

//Huske på: Vi MÅ sørge for at alle har den nyeste informasjonen når det gjøres en distribuering!
var latestState Config.GlobalState


func SpamGlobalState(gchan Config.GovernorChannels){ //Update and broadcast latest globalState from one peer
	//This is the global state
	ticker := time.Tick(2500 * time.Millisecond)

	transmitNet := make(chan Config.GlobalState)
	go bcast.Transmitter(16700, transmitNet)

	for{
		select{
			case <- ticker: // ticker
				//fmt.Println(t)
				fmt.Println("Transmit latestState", latestState)
				transmitNet <- latestState //sending latest state to network
				//fmt.Println("Transmitting GlobalState from ID: ", latestState.Id)
			case newUpdate:= <- gchan.InternalState:
				latestState = newUpdate
		}
	}
}

func UpdateGlobalState(gchan Config.GovernorChannels, FSMchan Config.FSMChannels, id int, GState Config.GlobalState){
	    //ticker := time.NewTicker(500 * time.Millisecond)
			distribute := false
			var hallOrder Config.ButtonEvent



	for{
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		gchan.InternalState <- GState
		if distribute{
			fmt.Println("GState: ", GState)
			//distribuere
			keyBestElevator:= ChooseElevator(GState.Map, hallOrder)
			fmt.Println("Distributing to: ", keyBestElevator)
			if keyBestElevator == GState.Id {
				var x = GState.Map[GState.Id]
				x.Queue[hallOrder.Floor][hallOrder.Button] = true
				GState.Map[GState.Id] = x
				FSMchan.PingFromGov <- GState
			}
			distribute = false
		}

		select{
			// case <- ticker.C:
			// 		FSMchan.PingFromGov <- GState
			// case Update:=  <- gchan.ExternalState: //StateUpdate from other Peer
			// fmt.Println("external state: ", Update.Id)
			// OutsideElev:= Update.Map[Update.Id]
			//
			// // 1. Set our info about the OutsideElev to the update from OutsideElev
			// GState.Map[Update.Id] = OutsideElev
			// // 2. If OutsideElev in state = DOOR_OPEN, remove orders in that floor
			// if OutsideElev.State == Config.DOOR_OPEN{
			// 		for button:= 0;  button < Config.NumButtons - 1 ; button++ {
			// 			GState.HallRequests[OutsideElev.Floor][button] = false
			// 		}
			// }
			// // 3. Fylle lokal HallRequests med de evt nye knappetrykkene fra oppdateringen
			// // Ta vare på disse knappetrykkene så de kan bli distribuert
			// newOrder := Config.ButtonEvent{}
			// for floor:= 0;  floor < Config.NumFloors; floor++ {
			// 	for button:= 0;  button < Config.NumButtons-1; button++ {
			// 		if Update.HallRequests[floor][button] && !GState.HallRequests[floor][button]{
			// 			GState.HallRequests[floor][button] = true
			// 			newOrder = Config.ButtonEvent{floor, Config.ButtonType(button)}
			// 			keyBestElevator:= ChooseElevator(GState.Map, newOrder)
			// 			if keyBestElevator == GState.Id {
			// 				var x = GState.Map[GState.Id]
			// 				x.Queue[hallOrder.Floor][hallOrder.Button] = true
			// 				GState.Map[GState.Id] = x
			// 				//FSMchan.PingFromGov <- GState
			// 			}
			// 		}
			// 	}
			// }

		case Id:= <- gchan.LostElev:
			lostElev:= GState.Map[Id]
			delete(GState.Map, Id)
			newOrder := Config.ButtonEvent{}

			for floor:= 0;  floor < Config.NumFloors; floor++ {
				for button:=0;  button < Config.NumButtons-1; button++ {
					if lostElev.Queue[floor][button]{
						newOrder = Config.ButtonEvent{floor, Config.ButtonType(button)}
						keyBestElevator:= ChooseElevator(GState.Map, newOrder)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[hallOrder.Floor][hallOrder.Button] = true
							GState.Map[GState.Id] = x
							//FSMchan.PingFromGov <- GState
						}
					}
				}
			}

		case hallOrder= <- gchan.AddHallOrder:
				GState.HallRequests[hallOrder.Floor][hallOrder.Button] = true
				distribute = true
		case cabOrderFloor:= <- FSMchan.AddCabOrderGov:
				var x = GState.Map[GState.Id]
				x.Queue[cabOrderFloor][2] = true
				GState.Map[GState.Id] = x
				FSMchan.PingFromGov <- GState

		case update:= <- FSMchan.LocalStateUpdate:
			tempQueue:= GState.Map[GState.Id].Queue
			GState.Map[GState.Id] = update
			var x = GState.Map[GState.Id]
			x.Queue = tempQueue
			GState.Map[GState.Id] = x


			//Dersom staten som oppdateringen har(alstå FSM-state) er DOOR_OPEN, fjern hall i denne etasjen
			if update.State == Config.DOOR_OPEN {
				for button:= 0;  button < Config.NumButtons - 1 ; button++ {
					GState.HallRequests[update.Floor][button] = false
				}
				var x = GState.Map[GState.Id]
				x.Queue[update.Floor][0] = false
				x.Queue[update.Floor][1] = false
				x.Queue[update.Floor][2] = false
				GState.Map[GState.Id] = x

			}
		}
	}
}

func NetworkState(gchan Config.GovernorChannels){
	seen:= make(map[int]time.Time)
	timeOut := Config.TIMEOUT*time.Second

	receiveNet := make(chan Config.GlobalState)
	go bcast.Receiver(16700, receiveNet)

	for{
		select{
		case StateUpdate:= <- receiveNet:
			fmt.Println("Receive state: ", StateUpdate)
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


func GovernorInit(GState Config.GlobalState, id int)Config.GlobalState  {
	Queue1:= [4][3]bool{{false,false,false},{false,false,false},{false,false,false} ,{false,false,false}}
	ElevState := Config.Elev{Config.UNKNOWN, Config.MD_Up, 0, Queue1}
	GState.Map = make(map[int]Config.Elev)
	GState.HallRequests= [4][2]bool{{false,false},{false,false},{false,false} ,{false,false}}
	GState.Id = id
	GState.Map[GState.Id] = ElevState
	fmt.Println("Init: ", GState)
	return GState
}
