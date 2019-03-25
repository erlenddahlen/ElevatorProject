package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"./CommTest"
	"./MasterController"
	"./SlaveFSM"
	"./elevio"
	"./network/localip"
)

func main() {
	// Setting Id
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")

	var port string
	flag.StringVar(&port, "port", "", "port of this peer")
	flag.Parse()

	elevio.Init("localhost:"+port, 4)
	fmt.Println("ID: ", id)

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	idInt, error := strconv.Atoi(id)
	if error != nil {
		fmt.Println(error.Error())
	}

	// Creating channels
	SlaveFSMChans := CommTest.SlaveFSMChannels{
		PosUpdate:    make(chan CommTest.ElevPos),
		CmdFinished:  make(chan CommTest.Cmd),
		ButtonPushed: make(chan CommTest.ButtonPushed),
		ButtonFromIo: make(chan elevio.ButtonEvent),
		CurrentFloor: make(chan int),
	}

	MasterControllerChans := CommTest.MasterControllerChannels{
		CmdElevToFloor: make(chan CommTest.Cmd),
		//updateLog:				make(chan ), // lage datatype til log
		// requestLogInt:		make(chan int),
		// requestLogExt:		make(chan int),
		LightMat: make(chan [4][3]int),
	}

	// LogChans := Log.LogChannels{	// Trenger kanskje ikke denne structen når vi bare har en kanal
	//   fullLog:				make(chan) // lage datatype til log
	// }

	CommunicationChans := CommTest.CommunicationChannels{
		//masterStateToMaster :       make(chan masterStateE),  //Make this state
		PosUpdateToMaster:    make(chan CommTest.ElevPos),
		CmdFinishedToMaster:  make(chan CommTest.Cmd),
		ButtonPushedToMaster: make(chan CommTest.ButtonPushed),
		//fullLogToMaster:           	make(chan), //Lage datatype til log

		CmdElevToFloorToSlave: make(chan int),
		LightMatToSlave:       make(chan [4][3]int),

		//updateLogToLog:              make(chan), //Lage datatype til log
		ReqLogIntToLog: make(chan int),
		ReqLogExtToLog: make(chan int),
		MergeLogToLog:  make(chan int),

		//numElev
		NumElev: make(chan int),
	}

	// Dersom lavest id i peerList er du master
	peerList := make(chan []string)
	go CommTest.PeerUpdating(id, CommunicationChans.NumElev, peerList)
	//fmt.Println("in main: ", peerList)
	var MasterState bool
	if idInt == 0 {
		MasterState = true
		go MasterController.MasterInit(CommunicationChans, MasterControllerChans)
	} else {
		MasterState = false
	}

	lightsFromMaster := make(chan [4][3]int)
	testMatrix := [4][3]int{}
	testMatrix[1][1] = 1
	testMatrix[2][2] = 1
	testMatrix[0][0] = 1

	go SlaveFSM.FSM(id, SlaveFSMChans, CommunicationChans, lightsFromMaster)

	go CommTest.CommunicationHandler(MasterState, id, SlaveFSMChans, CommunicationChans, MasterControllerChans)
	lightsFromMaster <- testMatrix

	for {
		select {
		case p := <-peerList:
			fmt.Println(p)
			// var min bool
			// min = true
			// for peer := 0; peer < len(peerList); peer++ {
			// 	if id > peerList[peer] {
			// 		MasterState = true
			// 		go MasterController.MasterInit(CommunicationChans, MasterControllerChans)
			// 	} else {
			// 		MasterState = false
			// 	}
		}
	}
}
