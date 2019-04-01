package Manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"../DataStructures"

	//	"../elevio"
	"time"

	"../PeerFSM"
	"../network/bcast"
)

var latestState DataStructures.GlobalState
var StateUpdate DataStructures.GlobalState

func SpamGlobalState(gchan DataStructures.ManagerChannels) { //Update and broadcast latest globalState from one peer
	//This is the global state
	ticker := time.Tick(1000 * time.Millisecond)
	transmitNet := make(chan DataStructures.GlobalState)
	//i:= 0
	go bcast.Transmitter(16700, transmitNet)

	// transmitTest := make(chan map[string]DataStructures.Elev)
	// go bcast.Transmitter(16700, transmitTest)

	for {
		select {
		case <-ticker: // ticker
			//fmt.Println(t)
			//fmt.Println("Transmit latestState", latestState)
			transmitNet <- latestState //sending latest state to network
			//fmt.Println("Transmitting GlobalState from ID: ", latestState.Id)

			// fmt.Println("Transmit Map", latestState.Map)
			// transmitTest <- latestState.Map

		case newUpdate := <-gchan.InternalState:
			latestState = newUpdate
			var file, err1 = os.Create(DataStructures.Backupfilename)
			isError(err1)
			GStatejason, _ := json.MarshalIndent(latestState, "", "")
			_ = ioutil.WriteFile(DataStructures.Backupfilename, GStatejason, 0644)
			file.Close()
			//fmt.Println("latestState:", latestState)
		}
	}
}

func UpdateGlobalState(gchan DataStructures.ManagerChannels, FSMchan DataStructures.FSMChannels, id string, GState DataStructures.GlobalState) {
	//ticker := time.NewTicker(500 * time.Millisecond)
	distribute := false
	var hallOrder DataStructures.ButtonEvent

	for {
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		gchan.InternalState <- GState
		//fmt.Println("distribute: ", distribute)
		if distribute {
			//fmt.Println("inside distribute")
			//distribuere
			keyBestElevator := ChooseElevator(GState.Map, hallOrder, GState)
			fmt.Println("Distributing to: ", keyBestElevator)
			if keyBestElevator == GState.Id {
				var x = GState.Map[GState.Id]
				x.Queue[hallOrder.Floor][hallOrder.Button] = true
				GState.Map[GState.Id] = x
				FSMchan.PingFromGov <- GState
				//fmt.Println("Q.i:", GState.Map[GState.Id].Queue)
			}
			distribute = false
		}
		if DataStructures.HasBackup{
			fmt.Println("has backup")
			DataStructures.HasBackup = false
			go func(FSMchan DataStructures.FSMChannels) {
			time.Sleep(3 * time.Second)
			FSMchan.PingFromGov <- GState
			fmt.Println("ping sent to FSM from gov")

			} (FSMchan)
		}

		select {
		// case <- ticker.C:
		// 		FSMchan.PingFromGov <- GState
		case Update := <-gchan.ExternalState: //StateUpdate from other Peer
			//fmt.Println("Gov_1 in")

			//fmt.Println("external state: ", Update.Id)
			OutsideElev := Update.Map[Update.Id]

			// 1. Set our info about the OutsideElev to the update from OutsideElev
			GState.Map[Update.Id] = OutsideElev
			// 2. If OutsideElev in state = DoorOpen, remove orders in that floor
			if OutsideElev.State == DataStructures.DoorOpen {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					GState.HallRequests[OutsideElev.Floor][button] = false
				}
			}
			// 3. Fylle lokal HallRequests med de evt nye knappetrykkene fra oppdateringen
			// Ta vare på disse knappetrykkene så de kan bli distribuert
			newOrder := DataStructures.ButtonEvent{}
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if Update.HallRequests[floor][button] && !GState.HallRequests[floor][button] {
						GState.HallRequests[floor][button] = true
						newOrder = DataStructures.ButtonEvent{floor, DataStructures.ButtonType(button)}
						fmt.Println("Order from other: ", newOrder)
						keyBestElevator := ChooseElevator(GState.Map, newOrder, GState)
						fmt.Println("Distributing to: ", keyBestElevator)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							FSMchan.PingFromGov <- GState
							//fmt.Println("Q.e:", GState.Map[GState.Id].Queue)
						}
					}
				}
			}
			PeerFSM.Lights(GState, GState.Map[id], id)
			//fmt.Println("Gov_1 out")

		case Id := <-gchan.LostElev:
			//fmt.Println("Gov_2 in")
			lostElev := GState.Map[Id]
			delete(GState.Map, Id)
			newOrder := DataStructures.ButtonEvent{}

			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if lostElev.Queue[floor][button] {
						newOrder = DataStructures.ButtonEvent{floor, DataStructures.ButtonType(button)}
						keyBestElevator := ChooseElevator(GState.Map, newOrder, GState)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							FSMchan.PingFromGov <- GState
						}
					}
				}
			}
			//fmt.Println("Gov_2 out")
		case hallOrder = <-gchan.AddHallOrder:
			//fmt.Println("Gov_3 in")
			//fmt.Println("Hall", GState.HallRequests)
			//fmt.Println("Got Hall from Peer")
			GState.HallRequests[hallOrder.Floor][hallOrder.Button] = true
			//			fmt.Println("set distrute to true")
			distribute = true
			//fmt.Println("Gov_3 out")
		case cabOrderFloor := <-FSMchan.AddCabOrderGov:
			//fmt.Println("Gov_4 in")
			//fmt.Println("Got Hall from Peer")
			var x = GState.Map[GState.Id]
			x.Queue[cabOrderFloor][2] = true
			GState.Map[GState.Id] = x
			//fmt.Println("Got Hall from Peer")
			FSMchan.PingFromGov <- GState
			//fmt.Println("Sent update to Peer")
			//fmt.Println("Gov_4 out")
		case update := <-FSMchan.LocalStateUpdate:
			//fmt.Println("Gov_5 in")
			tempQueue := GState.Map[GState.Id].Queue
			GState.Map[GState.Id] = update
			var x = GState.Map[GState.Id]
			x.Queue = tempQueue
			GState.Map[GState.Id] = x
			//Dersom staten som oppdateringen har(alstå FSM-state) er DoorOpen, fjern hall i denne etasjen
			if update.State == DataStructures.DoorOpen {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					GState.HallRequests[update.Floor][button] = false
				}
				var x = GState.Map[GState.Id]
				x.Queue[update.Floor][0] = false
				x.Queue[update.Floor][1] = false
				x.Queue[update.Floor][2] = false
				GState.Map[GState.Id] = x
				//fmt.Println("Q.D:", GState.Map[GState.Id].Queue)
			}
			PeerFSM.Lights(GState, GState.Map[id], id)
			//FSMchan.PingFromGov <- GState
			//fmt.Println("Gov_5 out")

		case <-gchan.Watchdog:
			fmt.Println("Watchdog timed out")
			// En-to heiser henger, vi tildeler alle hallreq til oss selv
			for floor := 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if GState.HallRequests[floor][button] {
						var x = GState.Map[GState.Id]
						x.Queue[floor][button] = true
						GState.Map[GState.Id] = x
					}
				}
			}
			FSMchan.PingFromGov <- GState
		}
	}
}

func NetworkState(gchan DataStructures.ManagerChannels) {
	seen := make(map[string]time.Time)
	timeOut := DataStructures.Timeout * time.Second

	receiveNet := make(chan DataStructures.GlobalState)
	go bcast.Receiver(16700, receiveNet)

	// receiveTest := make(chan map[int]DataStructures.Elev)
	// go bcast.Receiver(16701, receiveTest)

	for {
		select {
		case StateUpdate = <-receiveNet:
			//fmt.Println("GovNet_1 in")

			if latestState.Id != StateUpdate.Id {
				//fmt.Println("Receive state: ", StateUpdate)
				seen[StateUpdate.Id] = time.Now()
				//fmt.Println("Net in")
				gchan.ExternalState <- StateUpdate
				fmt.Println("hei")
				gchan.UpdatefromSpam <- StateUpdate
				fmt.Println("hallaaa")
				//fmt.Println("Net out")
				//fmt.Println("Received external state from ID:", StateUpdate.Id)
				//fmt.Println(StateUpdate)
			}
		//fmt.Println("GovNet_1 out")
		// case Update:= <- receiveTest:
		// 	fmt.Println("recTest: ", Update)
		default:
			for k, v := range seen {
				t := time.Now().Sub(v)
				if t > timeOut {
					gchan.LostElev <- k
					delete(seen, k)
					fmt.Println("deleted: ", k)
				}
			}
		}
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func ManagerInit(GState DataStructures.GlobalState, id string) DataStructures.GlobalState {
	DataStructures.HasBackup = false
	var _, err1 = os.Stat(DataStructures.Backupfilename)
	// If there is no backup, create one
	if os.IsNotExist(err1) {


		Queue1 := [4][3]bool{{false, false, false}, {false, false, false}, {false, false, false}, {false, false, false}}
		ElevState := DataStructures.Elev{DataStructures.Unknown, DataStructures.MD_Up, 0, Queue1}
		GState.Map = make(map[string]DataStructures.Elev)
		GState.HallRequests = [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
		GState.Id = id
		GState.Map[GState.Id] = ElevState

		var file, err1 = os.Create(DataStructures.Backupfilename)
		isError(err1)
		defer file.Close()
		GStatejason, _ := json.MarshalIndent(GState, "", "")
		_ = ioutil.WriteFile(DataStructures.Backupfilename, GStatejason, 0644)
		return GState
	} else {
		DataStructures.HasBackup = true
		//Read backup file
		GStateByte, err := ioutil.ReadFile(DataStructures.Backupfilename) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

		// Undo json
		error := json.Unmarshal(GStateByte, &GState)
		if error != nil {
			fmt.Println("error:", error)
		}
		GState.Map[GState.Id] = DataStructures.Elev{DataStructures.Unknown, DataStructures.MD_Up, 0, GState.Map[GState.Id].Queue}
		fmt.Printf("State: %s", GStateByte)
		return GState
	}
}

// //Watchdog should be own go routine
// //In SpamGlobalState, send latestState to Watchdog on own channel
// //Watchdog should be a case in for/select in UpdateGlobalState
// //When activated, put all Hallreq in own queue
//
//
// func Watchdog(gchan DataStructures.ManagerChannels){
// 	latestState
// 	ticker := time.Tick(6 * time.Second)
// 	var change bool
// 	var hallreqs bool
//
// 	for{
// 		select{
// 			case newUpdate := <- gchan.UpdatefromSpam
// 			change = false
// 			hallreqs = false
//
// 			//Iterate through Hallreqs
// 				//If exists orders hallreq = true
//
// 			//Iterate throug and check if chang in elev flor and state between latestState and newUpdate
// 				//If exists change = true
//
// 			if change OR !(hallreqs){
// 				reset ticker
// 			}
// 			case <- ticker.C
// 			gchan.pingfromwatch <- 1
// 		}
//
// 	}
//
// }
