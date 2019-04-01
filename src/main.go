package main

import (
	"flag"
	"fmt"

	"./DataStructures"
	"./FSM"
	"./Manager"
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
	//Init state, use backup if backup file available
	GState = Manager.ManagerInit(GState, id)

	FSMChannels := DataStructures.FSMChannels{
		AtFloor:            make(chan int),
		UpdateFromFSM:      make(chan DataStructures.Elev, 10),
		UpdateFromManager:   make(chan DataStructures.GlobalState),
		ButtonPushed:       make(chan DataStructures.ButtonEvent),
		AddCabOrder:        make(chan int),
		AddCabOrderManager: make(chan int),
	}

	ManagerChannels := DataStructures.ManagerChannels{
		InternalState: make(chan DataStructures.GlobalState),
		ExternalState:     make(chan DataStructures.GlobalState),
		LostElev:          make(chan string),
		AddHallOrder:      make(chan DataStructures.ButtonEvent),
		UpdatefromSpam:    make(chan DataStructures.GlobalState),
		MotorstopWatchdog: make(chan int),
	}

	go Manager.UpdateGlobalState(ManagerChannels, FSMChannels, id, GState)
	go Manager.SpamGlobalState(ManagerChannels)
	go Manager.UpdateNetworkPeers(ManagerChannels)
	go Manager.MotorstopWatchdog(ManagerChannels, GState)
	go FSM.FSM(ManagerChannels, FSMChannels, id, GState)

	for {
		//Run elevator
	}
}
