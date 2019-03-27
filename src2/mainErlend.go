package main

import (
    "fmt"
    "./Governor"
    "./Config"
//    "./PeerFSM"
    //"./elevio"
  //  "time"
    "flag"
    "strconv"

)

func main(){
    // Setting Id
    var id string
    flag.StringVar(&id, "id", "", "id of this peer")

    var port string
    flag.StringVar(&port, "port", "", "port of this peer")
    flag.Parse()

    idInt, error := strconv.Atoi(id)
  	if error != nil {
  		fmt.Println(error.Error())
  	}

    // elevio.Init("localhost:"+port, 4)
    // fmt.Println("ID: ", id)

    gchan:= Config.GovernorChannels{ //decalred in config file
    	InternalState: make(chan Config.GlobalState),
    	ExternalState: make(chan Config.GlobalState),
    	LostElev: make(chan int),
    }

    Queue1:= [4][3]bool{{false,false,false},{false,false,false},{false,false,true} ,{false,false,false}}
    //Queue2:= [4][3]bool{{false,false,false},{false,false,false},{false,false,false} ,{false,false,false}}

    temp1 := Config.Elev{Config.IDLE, Config.MD_Stop, 3, Queue1}
    //temp2 := Config.Elev{Config.IDLE, Config.DirStop, 2, Queue2}

    GState := Config.GlobalState{}
    GState.Map = make(map[int]Config.Elev)
    GState.Id = idInt

    GState.Map[idInt] = temp1
    //GState.Map["2"] = temp2

    go governor.SpamGlobalState(gchan)
    go governor.NetworkState(gchan)

    gchan.InternalState <- GState

    for{}
}
