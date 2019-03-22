package main

import (
	"flag"
	//"fmt"
  "./CommTest"
  "./SlaveFSM"
    "./elevio"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.


func main() {
  var port string
  flag.StringVar(&port, "port", "", "port of this peer")
  flag.Parse()

  elevio.Init("localhost:" + port, 4)

  // Setting Id
  var id string
  flag.StringVar(&id, "id", "", "id of this peer")
  flag.Parse()

  // Creating channels
  SlaveFSMChans := CommTest.SlaveFSMChannels{
  PosUpdate:  			make(chan int, 1),
  CmdFinished: 			make(chan int, 1),
  ButtonPushed:			make(chan elevio.ButtonEvent, 10),
  ButtonFromIo:			make(chan elevio.ButtonEvent),
  CurrentFloor:			make(chan int, 1),
  }

  // MasterControllerChans := MasterController.MasterControllerChannels{
  // cmdElevToFloor:		make(chan Communication.Cmd, 2),
  // updateLog:				make(chan ), // lage datatype til log
  // requestLogInt:		make(chan int, 1),
  // requestLogExt:		make(chan int, 1),
  // lightMat: 				make(chan [4][3] int, 12),
  // }
  //
  //
  // LogChans := Log.LogChannels{	// Trenger kanskje ikke denne structen når vi bare har en kanal
  //   fullLog:				make(chan) // lage datatype til log
  // }


  CommunicationChans := CommTest.CommunicationChannels{
    //masterStateToMaster :       make(chan masterStateE),  //Make this state
    PosUpdateToMaster:	    		make(chan int, 1),
    //posUpdateToMaster:	    		make(chan CommTest.ElevPos, 2),
    CmdFinishedToMaster: 	    	make(chan int, 1),
    ButtonPushedToMaster:	    	make(chan CommTest.ButtonPushed, 10),
    //fullLogToMaster:           	make(chan), //Lage datatype til log

    CmdElevToFloorToSlave:      make(chan CommTest.Cmd, 2),
    LightMatToSlave:		    		make(chan [4][3] int, 12),

    //updateLogToLog:              make(chan), //Lage datatype til log
    ReqLogIntToLog:              make(chan int, 1),
    ReqLogExtToLog:              make(chan int, 1),
    MergeLogToLog:               make(chan int, 1),

    //Must fix OUT channels
  }




  lightsFromMaster := make(chan [4][3]int, 12)
  testMatrix:=[4][3]int{}
  testMatrix[1][1] = 1
  testMatrix[2][2] = 1
  testMatrix[0][0] = 1

  button := make(chan elevio.ButtonEvent)
  command := make(chan int)
  go elevio.PollButtons(button)

  go SlaveFSM.FSM(command, SlaveFSMChans, lightsFromMaster)
  go CommTest.CommunicationHandler(id, port, SlaveFSMChans, CommunicationChans)

  lightsFromMaster <- testMatrix
  for {
		select {
		case btn := <-button:
			command <- btn.Floor
		}
	}

}
