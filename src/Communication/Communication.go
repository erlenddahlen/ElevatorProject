
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
    //From Slave
    case updatedSlavePos <- updateFromSlave:
        updatedSlavePos.id = myId
        if masterStateE == true {
            updatetSlavePosToMaster <- updatedSlavePos
        } else {
            updatetSlavePosToNH <- updatedSlavePos
        }

        //append id and create struct elevpos
        //if master in same node as slave
        //    master <- updatedSlavePos
        //else
        //    NH <- updatedSlavePos


    case cmdFinished <- cmdFinishedFromSlave
        cmdFinished.id = myId
        if masterStateE == true {
            ToMaster <- cmdFinished
        } else {
            ToNH <- cmdFinishedFromSlave
        }


    append id and create struct cmdfinished
    if master in same node as slave
        master <- struct
    else
        NH <- struct

    case inn <- buttonpushedfromSlave
    append id and create struct buttonpushed
    if master in same node as slave
        master <- struct
    else
        NH <- struct

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
