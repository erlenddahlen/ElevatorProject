package CommTest

import (
	"fmt"
	"../network/bcast"
	"../network/peers"
	"../elevio"
	//"../Log"
	"strconv"
	"time"
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
				PosUpdateToMaster	    	chan ElevPos
				CmdFinishedToMaster 	    chan Cmd
				ButtonPushedToMaster	    chan ButtonPushed
        //fullLogToMaster             chan //Lage datatype til log

        CmdElevToFloorToSlave       chan int
        LightMatToSlave		    chan [4][3] int

        //updateLogToLog              chan //Lage datatype til log
        ReqLogIntToLog              chan int
        ReqLogExtToLog              chan int
        MergeLogToLog               chan int

		//NumElev
		NumElev						chan int

}

type SlaveFSMChannels struct {
	PosUpdate  			chan ElevPos
	CmdFinished 		chan Cmd
	ButtonPushed		chan ButtonPushed
	ButtonFromIo		chan elevio.ButtonEvent
	CurrentFloor		chan int
}

type MasterControllerChannels struct {
		CmdElevToFloor						chan Cmd
//		ExtOPtimize							chan Cmd
		//updateLog							chan // lage datatype til log
		//requestLogInt						chan int
		//requestLogExt						chan int
		LightMat							chan [4][3] int
}

//type LogChannels struct {
//	fullLog	chan 		//lage datatype for Log
//}


func CommunicationHandler(MasterState bool, idString string, chSlave SlaveFSMChannels,
	chComm CommunicationChannels, chMaster MasterControllerChannels) {
	id, error:=strconv.Atoi(idString)
	if error != nil{
		fmt.Println(error.Error())
	}


	go PeerUpdating(idString, chComm.NumElev)
	go TransmitMessage(id, MasterState, chSlave, chComm , chMaster)
	go ReceiveMessage(id, MasterState, chSlave, chComm , chMaster)

}

func TransmitMessage(id int,MasterState bool, chSlave SlaveFSMChannels,
	chComm CommunicationChannels, chMaster MasterControllerChannels){

	transPosUpdate := make(chan ElevPos)
	transCmdFinished := make(chan Cmd)
	transButtonPushed := make(chan ButtonPushed)
	transCmd := make(chan Cmd)

	go bcast.Transmitter(16569, transPosUpdate)
	go bcast.Transmitter(16569, transCmdFinished)
	go bcast.Transmitter(16569, transButtonPushed)
	go bcast.Transmitter(16569, transCmd)

	timer:=time.NewTicker(1*time.Second)
	for{
		//Transmitter
		select{
		case <- timer.C:
			select{
			case a:= <- chSlave.PosUpdate:
				if MasterState{
					chComm.PosUpdateToMaster <- a
				} else {
					transPosUpdate <- a
					//fmt.Println("Sent PosExt")
				}
			case b:= <- chSlave.CmdFinished:
				if MasterState{
					chComm.CmdFinishedToMaster <- b
				} else {
					transCmdFinished<- b
					//fmt.Println("Sent cmdExt")
				}
			case msg:= <-chSlave.ButtonPushed:
				if MasterState{
					chComm.ButtonPushedToMaster <- msg
	//				fmt.Println("Button sent to Master", msg.Button)
				} else {
	//				fmt.Println("Button sent to other comm", msg.Button)
					transButtonPushed <- msg
					//fmt.Println("Sent ButtonPushedExt")
				}
			case msg:= <- chMaster.CmdElevToFloor:
				if msg.Id == id{
					chComm.CmdElevToFloorToSlave <- msg.Floor
				} else {
//					fmt.Println("Sending cmf from master")
					transCmd <- msg
					//fmt.Println("Sent ButtonPushedExt")
				}
			}
		}
	}
}


func ReceiveMessage(id int, MasterState bool, chSlave SlaveFSMChannels, chComm CommunicationChannels, chMaster MasterControllerChannels){
	recPosUpdate := make(chan ElevPos)
	recCmdFinished := make(chan Cmd)
	recButtonPushed := make(chan ButtonPushed)
	recCmd := make(chan Cmd)

	go bcast.Receiver(16569, recPosUpdate)
	go bcast.Receiver(16569, recCmdFinished)
	go bcast.Receiver(16569, recButtonPushed)
	go bcast.Receiver(16569, recCmd)


	//fmt.Println("RecMessage Init")
	for{
	//	fmt.Println("new round")
		//Receiver
		select{
		case a:= <- recPosUpdate:
			//fmt.Println("rec Posext")
			if id != a.Id && MasterState{
				chComm.PosUpdateToMaster <- a

			}
		case b:= <- recCmdFinished:
			//fmt.Println("rec cmdext")
			if id != b.Id && MasterState{
				chComm.CmdFinishedToMaster <- b
			}
		case msg:= <- recButtonPushed:
			//fmt.Println("rec buttonPushedExt")
			if id != msg.Id && MasterState{
//				fmt.Println("Button rec from other comm, sent to master: ", msg.Button)
				chComm.ButtonPushedToMaster <- msg
			}
		case msg:= <- recCmd:
			if id == msg.Id {
//				fmt.Println("sending from comm")
				chComm.CmdElevToFloorToSlave <- msg.Floor
			}
		}
	}
}




func PeerUpdating(id string, NumElev chan int) {
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15648, id, peerTxEnable)
	go peers.Receiver(15648, peerUpdateCh)

	for {
		select {
		case p := <-peerUpdateCh:
			// fmt.Printf("Peer update:\n")
			// fmt.Printf("  Peers:    %q\n", p.Peers)
			// fmt.Printf("  New:      %q\n", p.New)
			// fmt.Printf("  Lost:     %q\n", p.Lost)
			NumElev <- len(p.Peers)
		}
	}
}
