package main

import (
	"fmt"
	"flag"

	"./DataStructures"
	"./Governor"
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
	GState = Governor.GovernorInit(GState, id)

	PeerFSMChannels := DataStructures.FSMChannels{
		CurrentFloor:     make(chan int),
		LocalStateUpdate: make(chan DataStructures.Elev, 10),
		PingFromGov:      make(chan DataStructures.GlobalState),
		ButtonPushed:     make(chan DataStructures.ButtonEvent),
		AddCabOrderGov: 	make(chan int),
	}

	GovernorChannels := DataStructures.GovernorChannels{
		InternalState: 		make(chan DataStructures.GlobalState),
		ExternalState: 		make(chan DataStructures.GlobalState),
		LostElev:      		make(chan string),
		AddHallOrder:  		make(chan DataStructures.ButtonEvent),
		UpdatefromSpam: 	make(chan DataStructures.GlobalState),
		Watchdog:					make(chan int),
	}

	go Governor.UpdateGlobalState(GovernorChannels, PeerFSMChannels, id, GState)
	go Governor.SpamGlobalState(GovernorChannels)
	go Governor.NetworkState(GovernorChannels)
	go Governor.Watchdog(GovernorChannels, GState)
	go PeerFSM.FSM(GovernorChannels, PeerFSMChannels, id, GState)

	for {
		//Run elevator
	}
}
