package main

import (
	"fmt"
	"flag"

	"./DataStructures"
	"./Manager"
	"./PeerFSM"
	"./elevio"
)

func main() {

	// Setting id to the individual elevator
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")

	flag.Parse()
	elevio.Init("localhost:15657", 4)
	fmt.Println("ID: ", id)

//	Define a unique backup file name
	DataStructures.Backupfilename = "backup" + id + ".txt"

	var GState DataStructures.GlobalState
	GState = Manager.ManagerInit(GState, id)

	PeerFSMChannels := DataStructures.FSMChannels{
		CurrentFloor:     make(chan int),
		LocalStateUpdate: make(chan DataStructures.Elev, 10),
		PingFromGov:      make(chan DataStructures.GlobalState),
		ButtonPushed:     make(chan DataStructures.ButtonEvent),
		AddCabOrderGov: 	make(chan int),
	}

	ManagerChannels := DataStructures.ManagerChannels{
		InternalState: 		make(chan DataStructures.GlobalState),
		ExternalState: 		make(chan DataStructures.GlobalState),
		LostElev:      		make(chan string),
		AddHallOrder:  		make(chan DataStructures.ButtonEvent),
		UpdatefromSpam: 	make(chan DataStructures.GlobalState),
		Watchdog:					make(chan int),
	}

	go Manager.UpdateGlobalState(ManagerChannels, PeerFSMChannels, id, GState)
	go Manager.SpamGlobalState(ManagerChannels)
	go Manager.NetworkState(ManagerChannels)
	go Manager.Watchdog(ManagerChannels, GState)
	go PeerFSM.FSM(ManagerChannels, PeerFSMChannels, id, GState)

	for {
		//Run elevator
	}
}
