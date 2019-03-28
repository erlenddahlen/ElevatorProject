package main

import (
    "fmt"
    "./Governor"
    "./Config"
    "./PeerFSM"
    "./elevio"
    //"time"
    "flag"
    //"os"
    "strconv"
    //"./Peertest"
)

func main(){
    // Setting Id
    var id string
    flag.StringVar(&id, "id", "", "id of this peer")

    var port string
    flag.StringVar(&port, "port", "", "port of this peer")
    flag.Parse()

    elevio.Init("localhost:" + port, 4)
    fmt.Println("ID: ", id)

    // if id == "" {
    //     localIP, err := localip.LocalIP()
    //     if err != nil {
    //         fmt.Println(err)
    //         localIP = "DISCONNECTED"
    //     }
    //     id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
    // }

    idInt, error := strconv.Atoi(id)
    if error != nil {
        fmt.Println(error.Error())
    }


    PeerFSMChannels := Config.FSMChannels{
        CurrentFloor:       make(chan int),
        LocalStateUpdate:   make(chan Config.Elev, 10),
        PingFromGov:        make(chan Config.GlobalState),
        ButtonPushed:       make(chan Config.ButtonEvent),
        //AddCabOrder:        make(chan int),
        AddCabOrderGov:     make(chan int),
    }

    GovernorChannels:= Config.GovernorChannels{
        InternalState:      make(chan Config.GlobalState),
        ExternalState:      make(chan Config.GlobalState),
        LostElev:           make(chan int),
        AddHallOrder:       make(chan Config.ButtonEvent),
    }
    var GState Config.GlobalState
    GState = Governor.GovernorInit(GState, idInt)
    //fmt.Println("Gstate init: ", GState)

    go Governor.UpdateGlobalState(GovernorChannels, PeerFSMChannels, idInt, GState)
    go Governor.SpamGlobalState(GovernorChannels)
    go Governor.NetworkState(GovernorChannels)
    go PeerFSM.FSM(GovernorChannels, PeerFSMChannels, idInt, GState)


    for{
    }
}
