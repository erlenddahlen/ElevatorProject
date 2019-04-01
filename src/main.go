package main

import (
	"fmt"

	"./DataStructures"
	"./Governor"
	"./PeerFSM"
	"./elevio"

	//"time"
	"flag"
	//"os"
	//"./Peertest"
)

func main() {
	// Setting Id
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")

	// var port string
	// flag.StringVar(&port, "port", "", "port of this peer")
	flag.Parse()

	elevio.Init("localhost:15657", 4)
	fmt.Println("ID: ", id)

	DataStructures.Backupfilename = "backup" + id + ".txt"

	// if id == "" {
	//     localIP, err := localip.LocalIP()
	//     if err != nil {
	//         fmt.Println(err)
	//         localIP = "DISCONNECTED"
	//     }
	//     id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	// }

	var GState DataStructures.GlobalState
	GState = Governor.GovernorInit(GState, id)
	//fmt.Println("Gstate init: ", GState)

	PeerFSMChannels := DataStructures.FSMChannels{
		CurrentFloor:     make(chan int),
		LocalStateUpdate: make(chan DataStructures.Elev, 10),
		PingFromGov:      make(chan DataStructures.GlobalState),
		ButtonPushed:     make(chan DataStructures.ButtonEvent),
		//AddCabOrder:        make(chan int),
		AddCabOrderGov: make(chan int),
	}

	GovernorChannels := DataStructures.GovernorChannels{
		InternalState: make(chan DataStructures.GlobalState),
		ExternalState: make(chan DataStructures.GlobalState),
		LostElev:      make(chan string),
		AddHallOrder:  make(chan DataStructures.ButtonEvent),
		UpdatefromSpam: make(chan DataStructures.GlobalState),
		Watchdog:		make(chan int),
	}

	go Governor.UpdateGlobalState(GovernorChannels, PeerFSMChannels, id, GState)
	go Governor.SpamGlobalState(GovernorChannels)
	go Governor.NetworkState(GovernorChannels)
	go Governor.Watchdog(GovernorChannels, GState)
	go PeerFSM.FSM(GovernorChannels, PeerFSMChannels, id, GState)

	for {
	}
}
