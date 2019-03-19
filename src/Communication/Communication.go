package Communication


// Sette min id
myId :=

type ElevPos struct {
	Floor  int
	Id     int
}


type CmdFinished struct {
	Floor  int
	Id     int
}


type ButtonPushed struct {
	Button ButtonEvent
    Id     int
}

func localHandler(updateFromSlave chan ElevPos){
for{
    select{
    // From Slave
		// Adding id to information
		// Sending information either to internal master or to NH(external master)
    case update <- updateFromSlave:
        updatedSlavePos.id = myId
        if masterStateE == true {
            updatetSlavePosToMaster <- update
        } else {
            updatetSlavePosToNH <- update
        }

    case cmdFinished <- cmdFinishedFromSlave
        cmdFinished.id = myId
        if masterStateE == true {
            cmdFinishedToMaster <- cmdFinished
        } else {
            cmdFinishedToNH <- cmdFinished
        }


    case buttonPushed <- buttonPushedFromSlave
			buttonPushed.id = myId
			if masterStateE == true {
					buttonPushedToMaster <- buttonPushed
			} else {
					buttonPushedToNH <- buttonPushed
			}


    // From Master
    case inn <- ElevtoflorrfromfromMaster
    if id == own
    slave <- cmdelevtofloor
    else
    NH <- cmdleevtofloor

    case inn <- updatelogfromamster
        log <- logupdate
        NH <- logupdate

    case inn <- requestlogintfromamster
        log <- reqlog
    case inn <- requestlogextfrommaster
        NH <- reqlog

    //From log
    case inn <- fullogfromlog
        if masterState == own
            master <- log
        else
            NH <- log

    //From NH
    case inn <- posupdate
    master <-

    case inn <- cmdfinished
    master <-

    case inn <- buttonpushed
    master <-

    case inn <- elevtofloor
    master <-

    case inn <- updatelog
    log <-

    case inn <- reqlogext
    log <-
    }
}
}

func networkHandler(){
    for
        select
        //fromLH
        case inn <- posUpdate
        append id
        other NH <- //Need to fix UDP/channel with network module


        //same holds for rest of the cases
        case inn <- orderinished
        case inn <- buttonpushed
        case inn <- elevtofloor
        case inn <- updatelog
        case inn <- reqlogext
        case inn <- fullog

        //fromONH
        case inn <- posUpdate
        // send ACK to posUpdate-channel
        LH <-


        case inn <- cmdinished
        LH <-

        case inn <- buttonpushed
        LH <-

        case inn <- elevtofloor
        LH <-

        case inn <- updatelog
        LH <-

        case inn <- reqlogext
        LH <-

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
