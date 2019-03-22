package CommTest

import (
	"fmt"
	"os"
  "../network/bcast"
  "../network/localip"
  "../network/peers"
	"../elevio"
	"../Log"
)

type HelloMsg struct {
	Message string
	Iter    int
}

type ElevPos struct {
	Floor  int
	Id     int
}


type Cmd struct {
	Floor  int
	Id     int
}


type ButtonPushed struct {
	Button elevio.ButtonEvent
        Id     int
}

type CommunicationChannels struct {
        //masterStateToMaster         chan masterStateE   //Make this state
				PosUpdateToMaster	    	chan int
				CmdFinishedToMaster 	    chan int
				ButtonPushedToMaster	    chan ButtonPushed
        //fullLogToMaster             chan //Lage datatype til log

        CmdElevToFloorToSlave       chan Cmd
        LightMatToSlave		    chan [4][3] int

        //updateLogToLog              chan //Lage datatype til log
        ReqLogIntToLog              chan int
        ReqLogExtToLog              chan int
        MergeLogToLog               chan int

        //Must fix OUT channels
}
// Declaring channels
type SlaveFSMChannels struct {
	PosUpdate  			chan int
	CmdFinished 		chan int
	ButtonPushed		chan elevio.ButtonEvent
	ButtonFromIo		chan elevio.ButtonEvent
	CurrentFloor		chan int
}

//type LogChannels struct {
//	fullLog	chan 		//lage datatype for Log
//}


func CommunicationHandler(id string, chSlave SlaveFSMChannels, chComm CommunicationChannels) {
	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}



	// PEER LIST
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15648, id, peerTxEnable)
	go peers.Receiver(15648, peerUpdateCh)



	// SEND - RECEIVE
	// We make channels for sending and receiving our custom data types
	//helloTx := make(chan HelloMsg)
	//helloRx := make(chan HelloMsg)



	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, chSlave.PosUpdate)
	go bcast.Receiver(16569, chComm.PosUpdateToMaster)

	// The example message. We just send one of these every second.
	// go func() {
	// 	//helloMsg := HelloMsg{"Hello from " + id, 0}
	// 	for {
	// 		helloMsg.Iter++
	// 		helloTx <- helloMsg
	//
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-chComm.PosUpdateToMaster:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
