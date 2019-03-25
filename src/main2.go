//package main

// import "./elevio"
// //import "./test"
// import "./SlaveFSM"
//
//
// // ny terminal der simulator ligger: ./SimElevatorServer
//
// func main() {
//
// 	// Creating channels
// 	SlaveFSMChans := SlaveFSM.SlaveFSMChannels{
// 	posUpdate:  			make(chan Communication.ElevPos, 2),
// 	cmdFinished: 			make(chan int, 1),
// 	buttonPushed:			make(chan Communication.ButtonPushed, 10),
// 	buttonFromIo:			make(chan elevio.ButtonEvent),
// 	currentFloor:			make(chan int, 1),
// 	}
//
// 	MasterControllerChans := MasterController.MasterControllerChannels{
// 	cmdElevToFloor:		make(chan Communication.Cmd, 2),
// 	updateLog:				make(chan ), // lage datatype til log
// 	requestLogInt:		make(chan int, 1),
// 	requestLogExt:		make(chan int, 1),
// 	lightMat: 				make(chan [4][3] int, 12),
// 	}
//
//
// 	LogChans := Log.LogChannels{	// Trenger kanskje ikke denne structen når vi bare har en kanal
// 		fullLog:				make(chan) // lage datatype til log
// 	}
//
//
// 	CommunicationChans := Communication.CommunicationChannels{
// 		masterStateToMaster :       make(chan masterStateE),  //Make this state
// 		posUpdateToMaster:	    		make(chan Communication.ElevPos, 2),
// 		cmdFinishedToMaster: 	    	make(chan int, 1),
// 		buttonPushedToMaster:	    	make(chan Communication.ButtonPushed, 10),
// 		fullLogToMaster:           	make(chan), //Lage datatype til log
//
// 		cmdElevToFloorToSlave:      make(chan Communication.Cmd, 2),
// 		lightMatToSlave:		    		make(chan [4][3] int, 12),
//
// 		updateLogToLog:              make(chan), //Lage datatype til log
// 		reqLogIntToLog:              make(chan int, 1),
// 		reqLogExtToLog:              make(chan int, 1),
// 		mergeLogToLog:               make(chan int, 1),
//
// 		//Must fix OUT channels
// 	}
//
//
// 	elevio.Init("localhost:15657", 4)
// 	posUpdateToMaster := make(chan int, 1) // Posisjon
// 	posUpdate := make(chan int, 2)					// Posisjon + id
//
// 	button := make(chan elevio.ButtonEvent)
// 	command := make(chan int)
// 	go elevio.PollButtons(button)
//
//
// 	go Communication.CommTest(SlaveKanal, MasterKanal)
// 	//go SlaveFSM.FSM(command, SlaveFSMChans, lightsFromMaster)
// 	//go Communication.communicationHandler(SlaveFSMChans, )
// 	//go masterInit()
//
//
// }
