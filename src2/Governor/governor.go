package Governor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"../Config"

	//	"../elevio"
	"time"

	"../PeerFSM"
	"../network/bcast"
)

//Huske på: Vi MÅ sørge for at alle har den nyeste informasjonen når det gjøres en distribuering!
var latestState Config.GlobalState
var StateUpdate Config.GlobalState

func SpamGlobalState(gchan Config.GovernorChannels) { //Update and broadcast latest globalState from one peer
	//This is the global state
	ticker := time.Tick(500 * time.Millisecond)
	transmitNet := make(chan Config.GlobalState)
	go bcast.Transmitter(16700, transmitNet)

	// transmitTest := make(chan map[string]Config.Elev)
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
			gchan.TakeBackup <- latestState
			//fmt.Println("latestState:", latestState)
		}
	}
}

func UpdateGlobalState(gchan Config.GovernorChannels, FSMchan Config.FSMChannels, id string, GState Config.GlobalState) {
	//ticker := time.NewTicker(500 * time.Millisecond)
	distribute := false
	var hallOrder Config.ButtonEvent

	for {
		//Watchdog <- globalState, after 6 sec without activity in network, peer takes ownership of all global orders
		gchan.InternalState <- GState
		fmt.Println("distribute: ", distribute)
		if distribute {
			fmt.Println("inside distribute")
			//distribuere
			keyBestElevator := ChooseElevator(GState.Map, hallOrder, GState)
			fmt.Println("Distributing to: ", keyBestElevator)
			if keyBestElevator == GState.Id {
				var x = GState.Map[GState.Id]
				x.Queue[hallOrder.Floor][hallOrder.Button] = true
				GState.Map[GState.Id] = x
				FSMchan.PingFromGov <- GState
				fmt.Println("Q.i:", GState.Map[GState.Id].Queue)
			}
			distribute = false
		}

		select {
		// case <- ticker.C:
		// 		FSMchan.PingFromGov <- GState
		case Update := <-gchan.ExternalState: //StateUpdate from other Peer
			//fmt.Println("external state: ", Update.Id)
			OutsideElev := Update.Map[Update.Id]

			// 1. Set our info about the OutsideElev to the update from OutsideElev
			GState.Map[Update.Id] = OutsideElev
			// 2. If OutsideElev in state = DOOR_OPEN, remove orders in that floor
			if OutsideElev.State == Config.DOOR_OPEN {
				for button := 0; button < Config.NumButtons-1; button++ {
					GState.HallRequests[OutsideElev.Floor][button] = false
				}
			}
			// 3. Fylle lokal HallRequests med de evt nye knappetrykkene fra oppdateringen
			// Ta vare på disse knappetrykkene så de kan bli distribuert
			newOrder := Config.ButtonEvent{}
			for floor := 0; floor < Config.NumFloors; floor++ {
				for button := 0; button < Config.NumButtons-1; button++ {
					if Update.HallRequests[floor][button] && !GState.HallRequests[floor][button] {
						GState.HallRequests[floor][button] = true
						newOrder = Config.ButtonEvent{floor, Config.ButtonType(button)}
						keyBestElevator := ChooseElevator(GState.Map, newOrder, GState)
						fmt.Println("Distributing to: ", keyBestElevator)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[newOrder.Floor][newOrder.Button] = true
							GState.Map[GState.Id] = x
							FSMchan.PingFromGov <- GState
							fmt.Println("Q.e:", GState.Map[GState.Id].Queue)
						}
					}
				}
			}
			PeerFSM.Lights(GState, GState.Map[id], id)

		case Id := <-gchan.LostElev:
			lostElev := GState.Map[Id]
			delete(GState.Map, Id)
			newOrder := Config.ButtonEvent{}

			for floor := 0; floor < Config.NumFloors; floor++ {
				for button := 0; button < Config.NumButtons-1; button++ {
					if lostElev.Queue[floor][button] {
						newOrder = Config.ButtonEvent{floor, Config.ButtonType(button)}
						keyBestElevator := ChooseElevator(GState.Map, newOrder, GState)
						if keyBestElevator == GState.Id {
							var x = GState.Map[GState.Id]
							x.Queue[hallOrder.Floor][hallOrder.Button] = true
							GState.Map[GState.Id] = x
							FSMchan.PingFromGov <- GState
						}
					}
				}
			}

		case hallOrder = <-gchan.AddHallOrder:
			//fmt.Println("Hall", GState.HallRequests)
			fmt.Println("Got Hall from Peer")
			GState.HallRequests[hallOrder.Floor][hallOrder.Button] = true
			fmt.Println("set distrute to true")
			distribute = true
		case cabOrderFloor := <-FSMchan.AddCabOrderGov:
			fmt.Println("Got Hall from Peer")
			var x = GState.Map[GState.Id]
			x.Queue[cabOrderFloor][2] = true
			GState.Map[GState.Id] = x
			fmt.Println("Got Hall from Peer")
			FSMchan.PingFromGov <- GState
			fmt.Println("Sent update to Peer")
		case update := <-FSMchan.LocalStateUpdate:
			tempQueue := GState.Map[GState.Id].Queue
			GState.Map[GState.Id] = update
			var x = GState.Map[GState.Id]
			x.Queue = tempQueue
			GState.Map[GState.Id] = x

			//Dersom staten som oppdateringen har(alstå FSM-state) er DOOR_OPEN, fjern hall i denne etasjen
			if update.State == Config.DOOR_OPEN {
				for button := 0; button < Config.NumButtons-1; button++ {
					GState.HallRequests[update.Floor][button] = false
				}
				var x = GState.Map[GState.Id]
				x.Queue[update.Floor][0] = false
				x.Queue[update.Floor][1] = false
				x.Queue[update.Floor][2] = false
				GState.Map[GState.Id] = x
				fmt.Println("Q.D:", GState.Map[GState.Id].Queue)
			}
			PeerFSM.Lights(GState, GState.Map[id], id)
		}
	}
}

func NetworkState(gchan Config.GovernorChannels) {
	seen := make(map[string]time.Time)
	timeOut := Config.TIMEOUT * time.Second

	receiveNet := make(chan Config.GlobalState)
	go bcast.Receiver(16700, receiveNet)

	// receiveTest := make(chan map[int]Config.Elev)
	// go bcast.Receiver(16701, receiveTest)

	for {
		select {
		case StateUpdate = <-receiveNet:

			if latestState.Id != StateUpdate.Id {
				//fmt.Println("Receive state: ", StateUpdate)
				seen[StateUpdate.Id] = time.Now()
				gchan.ExternalState <- StateUpdate
				//fmt.Println("Received external state from ID:", StateUpdate.Id)
				//fmt.Println(StateUpdate)
			}
		// case Update:= <- receiveTest:
		// 	fmt.Println("recTest: ", Update)
		default:
			for k, v := range seen {
				t := time.Now().Sub(v)
				if t > timeOut {
					gchan.LostElev <- k
					delete(seen, k)
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
func BackupState(gchan Config.GovernorChannels) {
	for {
		GState := <-gchan.TakeBackup
		var file, err1 = os.Create(Config.Backupfilename)
		isError(err1)
		GStatejason, _ := json.MarshalIndent(GState, "", "")
		_ = ioutil.WriteFile(Config.Backupfilename, GStatejason, 0644)
		file.Close()
	}
}

func GovernorInit(GState Config.GlobalState, id string) Config.GlobalState {
	var _, err1 = os.Stat(Config.Backupfilename)
	// If there is no backup, create one
	if os.IsNotExist(err1) {

		Queue1 := [4][3]bool{{false, false, false}, {false, false, false}, {false, false, false}, {false, false, false}}
		ElevState := Config.Elev{Config.UNKNOWN, Config.MD_Up, 0, Queue1}
		GState.Map = make(map[string]Config.Elev)
		GState.HallRequests = [4][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
		GState.Id = id
		GState.Map[GState.Id] = ElevState

		var file, err1 = os.Create(Config.Backupfilename)
		isError(err1)
		defer file.Close()
		GStatejason, _ := json.MarshalIndent(GState, "", "")
		_ = ioutil.WriteFile(Config.Backupfilename, GStatejason, 0644)
		return GState
	} else {
		//Read backup file
		GStateByte, err := ioutil.ReadFile(Config.Backupfilename) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

		// Undo json
		error := json.Unmarshal(GStateByte, &GState)
		if error != nil {
			fmt.Println("error:", error)
		}
		GState.Id = id //Do not overwrite Id
		return GState
	}
}
