//package Communication

import "../elevio"

// Sette min id, gjøres i main? så kan man ta den inn i funskjonen her?
myId :=

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
        masterStateToMaster         chan masterStateE   //Make this state
				posUpdateToMaster	    chan ElevPos
				cmdFinishedToMaster 	    chan int
				buttonPushedToMaster	    chan ButtonPushed
        fullLogToMaster             chan //Lage datatype til log

        cmdElevToFloorToSlave       chan Cmd
        lightMatToSlave		    chan [4][3] int

        updateLogToLog              chan //Lage datatype til log
        reqLogIntToLog              chan int
        reqLogExtToLog              chan int
        mergeLogToLog               chan int

        //Must fix OUT channels
}

func communicationHandler(   ){
    // Channels inn
        updateFromSlave := make(chan ElevPos, 2)
        cmdFinishedFromSlave := make(chan Cmd, 2)
        buttonPushedFromSlave := make(chan ButtonPushed, 10)

        cmdElevToFloorFromMaster := make(chan Cmd, 2)
        updateLogFromMaster := make(chan , )            // definere og lage type for log
        reqLogIntFromMaster := make(chan int, 1)
        reqLogExtFromMaster := make(chan int, 1)

        fullLogFromLog := make(chan, )                 // definere og lage type for log

    // Channels out
        // Must make the OUT channel, the connection to the GO-network module
        posUpdateToMaster := make(chan ElevPos, 2)
        cmdFinishedToMaster := make(chan Cmd, 2)
        buttonPushedToMaster := make(chan ButtonPushed, 10)

        cmdElevToFloorToSlave := make(chan Cmd, 2)

        updateLogToLog := make(chan , )            // definere og lage type for log
        reqLogIntToLog := make(chan int, 1)
        reqLogExtToLog := make(chan int, 1)

        fullLogToMaster := make(chan , )            // definere og lage type for log



for{
    select{
    // From Slave

	// Adding id to information
	// Sending information either to internal master or to NH(external master)
    case update <- updateFromSlave:
        updatedSlavePos.id = myId
        if masterStateE == true {
            updatedSlavePosToMaster <- update
        } else {
            OUT<- update
        }

    case cmdFinished <- cmdFinishedFromSlave
        cmdFinished.id = myId
        if masterStateE == true {
            cmdFinishedToMaster <- cmdFinished
        } else {
            OUT <- cmdFinished
        }

    case buttonPushed <- buttonPushedFromSlave
			buttonPushed.id = myId
			if masterStateE == true {
					buttonPushedToMaster <- buttonPushed
			} else {
					OUT <- buttonPushed
			}


    // From Master

    // Checking if id is the internal id
    // Sending information either to internal slave or NH(external slave)
    case cmdElevToFloor <- cmdElevToFloorFromMaster
        if cmdElevToFloor.id == myId  {
            cmdElevToFloorToSlave <- cmdElevToFloor
        } else {
            OUT <- cmdElevToFloor
        }

    // Checking if id is the internal id
    // Sending information either to internal slave or NH(external slave)
    case updateLog <- updateLogFromMaster
        if updateLog.id == myId  {
            updateLogToLog <- updateLog
        } else {
            OUT <- updateLog
        }

    // Requesting internal or external log
    case reqLogInt <- reqLogIntFromMaster
        reqLogIntToLog <- reqLogInt

    case reqLogExt <- reqLogExtFromMaster
        OUT <- reqLogExt

    // From log
    // Sending fullLog to either master or NH(external master)
    case fullLog <- fullLogFromLog
        if masterStateE == true {
                fullLogToMaster <- fullLog
        } else {
                OUT <- fullLog
        }

    // From OUT
    // All external information coming into the node from other nodes
    case posUpdate <- posUpdateFromOUT
        updatedSlavePosToMaster <- posUpdate

    case cmdFinished <- cmdFinishedFromOUT
        cmdFinishedToMaster <- cmdFinished

    case buttonPushed <- buttonPushedFromOUT
        buttonPushedToMaster <- buttonPushed

    case cmdElevToFloor <- cmdElevToFloorFromOUT
        cmdElevToFloorToSlave <- cmdElevToFloor

    case updateLog <- updateLogFromOUT
        updateLogToLog <- updateLog

    case reqLogExt <- reqLogExtFromOUT
        reqLogExtToLog <- reqLogExt
        }
    }
}

func checkNetworkStatus(){


}

//ACKING:
//Use prefix to determine if hash or msg
//Use global datastruct
//Use channelID/counter
//Universal timer used to resend messages from datastruct

//All incoming chans in NH --> send ack
//All outgoing chans in NH --> list of ACK, resend
