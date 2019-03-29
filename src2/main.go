package main

import (
	"fmt"

	"./Config"
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

	Config.Backupfilename = "backup" + id + ".txt"

	// if id == "" {
	//     localIP, err := localip.LocalIP()
	//     if err != nil {
	//         fmt.Println(err)
	//         localIP = "DISCONNECTED"
	//     }
	//     id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	// }

	var GState Config.GlobalState
	GState = Governor.GovernorInit(GState, id)
	//fmt.Println("Gstate init: ", GState)

	PeerFSMChannels := Config.FSMChannels{
		CurrentFloor:     make(chan int),
		LocalStateUpdate: make(chan Config.Elev, 10),
		PingFromGov:      make(chan Config.GlobalState),
		ButtonPushed:     make(chan Config.ButtonEvent),
		//AddCabOrder:        make(chan int),
		AddCabOrderGov: make(chan int),
	}

	GovernorChannels := Config.GovernorChannels{
		InternalState: make(chan Config.GlobalState),
		ExternalState: make(chan Config.GlobalState),
		LostElev:      make(chan string),
		AddHallOrder:  make(chan Config.ButtonEvent),
	}

	go Governor.UpdateGlobalState(GovernorChannels, PeerFSMChannels, id, GState)
	go Governor.SpamGlobalState(GovernorChannels)
	go Governor.NetworkState(GovernorChannels)
	go Governor.Watchdog(GovernorChannels, GState)
	go PeerFSM.FSM(GovernorChannels, PeerFSMChannels, id, GState)

	for {
	}
}
