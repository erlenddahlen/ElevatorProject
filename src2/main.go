package main

import (
    "fmt"
    //"./Governor"
    "./Config"
    "./PeerFSM"
    "./elevio"
    "time"
    "flag"
    //"os"
    //"strconv"
    //"./Peertest"
)

func main(){
    // Setting Id
    var id string
    flag.StringVar(&id, "id", "", "id of this peer")

    var port string
    flag.StringVar(&port, "port", "", "port of this peer")
    flag.Parse()

    elevio.Init("localhost:"+port, 4)
    fmt.Println("ID: ", id)

    // if id == "" {
    //     localIP, err := localip.LocalIP()
    //     if err != nil {
    //         fmt.Println(err)
    //         localIP = "DISCONNECTED"
    //     }
    //     id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
    // }

    // idInt, error := strconv.Atoi(id)
    // if error != nil {
    //     fmt.Println(error.Error())
    // }


    PeerFSMChannels := Config.FSMChannels{
        CurrentFloor:       make(chan int),
        LocalStateUpdate:   make(chan Config.Elev),
        PingFromGov:        make(chan Config.GlobalState),
        ButtonPushed:       make(chan Config.ButtonEvent),
        AddCabOrder:        make(chan int),
    }

    GovernorChannels:= Config.GovernorChannels{
        AddHallOrder:       make(chan Config.ButtonEvent),
    }

    Queue1:= [4][3]bool{{true,false,false},{false,true,false},{false,false,true} ,{false,false,true}}
    //Queue2:= [4][3]bool{{false,false,false},{false,false,false},{false,false,false} ,{false,false,false}}
    temp1 := Config.Elev{Config.UNKNOWN, Config.MD_Up, 0, Queue1}
    //temp2 := Config.Elev{Config.IDLE, Config.DirStop, 2, Queue2}

    GState := Config.GlobalState{}
    GState.Map = make(map[string]Config.Elev)
    GState.HallRequests= [4][2]bool{{true,false},{true,false},{false,false} ,{false,false}}
    GState.Map[id] = temp1
    //GState.Map["2"] = temp2
    ticker := time.NewTicker(1 * time.Second)
    //PeerFSM.Lights(GState, id)
    //go elevio.PollFloorSensor(PeerFSMChannels.CurrentFloor)
    go PeerFSM.FSM(GovernorChannels, PeerFSMChannels, id, GState)


    for{
        select{
        case <- ticker.C:
            PeerFSMChannels.PingFromGov <- GState
            fmt.Println("QUEUE in main: ", GState.Map[id].Queue)
            //PeerFSM.Lights(GState, id)
        case update:= <- PeerFSMChannels.LocalStateUpdate:
            fmt.Println("in main, floor: ", update.Floor)
            GState.Map[id] = update
            fmt.Println("QUEUE in main2: ", GState.Map[id].Queue)
        case hallOrder:= <- GovernorChannels.AddHallOrder:
            GState.HallRequests[hallOrder.Floor][hallOrder.Button] = true
            //GState.Map[id].Queue[hallOrder.Floor][hallOrder.Button] = true
        }
    }

    // NewButton := Config.ButtonEvent{Floor:0, Button: Config.BT_HallUp}
    // TEST OF OPTIMIZING/DISTRIBUTION
    // Husk at OrderAbove og OrderBelow går gjennom alle knappene fra FSM, mens her vil vi bare ha 0 og 1
    // fmt.Println(ChooseElevator(GState.Map, NewButton))
    // //for{}
}

// func ChooseElevator(elevators map[string]Config.Elev, NewOrder Config.ButtonEvent) string{ //Should reurn Config.Elev
// 	times := make([]int, len(elevators))
// 	indexBestElevator := 0
// 	keyBestElevator := ""
//
// 	//trying to convert to map
// 	i := 0
// 	for key, _ := range elevators {
// 		times[i] = TimeToServeOrder(elevators[key], NewOrder)
// 		if times[i] <= times[indexBestElevator] {
// 			indexBestElevator = i
// 			keyBestElevator = key
//         }
// 		i++
// 	}
//
// 	//return elevators[keyBestElevator]
//     return keyBestElevator
// }

/*
func AssignOrder(elevator Elevator, Order elevio.ButtonEvent) {
	//Assign the order to the chosen elevator
}
*/

// func TimeToServeOrder(e Config.Elev, button Config.ButtonEvent) int {
// 	tempElevator := e
// 	tempElevator.Queue[button.Floor][button.Button] = true
//
// 	timeUsed := 0
//
// 	switch tempElevator.State {
// 	case Config.IDLE:
// 		tempElevator.Dir = Peertest.GetNextDir(tempElevator)
// 		if tempElevator.Dir == Config.DirStop {
// 			return timeUsed
// 		}
// 	case Config.MOVING:
// 		timeUsed += Config.TRAVEL_TIME / 2
// 		tempElevator.Floor += int(tempElevator.Dir)
// 	case Config.DOOR_OPEN:
// 		timeUsed += Config.DOOR_OPEN_TIME / 2
// 	}
//     count := 0
// 	for {
// 		if Peertest.ShouldStop(tempElevator) {
// 			if tempElevator.Floor == button.Floor {
// 				//We have "arrived" at destination
// 				return timeUsed
// 			}
// 			tempElevator.Queue[tempElevator.Floor][Config.BT_Cab] = false
// 			if tempElevator.Dir == Config.DirUp || !Peertest.OrderAbove(tempElevator) {
// 				tempElevator.Queue[tempElevator.Floor][Config.BT_HallUp] = false
// 			}
// 			if tempElevator.Dir == Config.DirDown || !Peertest.OrderBelow(tempElevator) {
// 				tempElevator.Queue[tempElevator.Floor][Config.BT_HallDown] = false
// 			}
// 			timeUsed += Config.DOOR_OPEN_TIME
// 			tempElevator.Dir = Peertest.GetNextDir(tempElevator)
// 		}
// 		tempElevator.Floor += int(tempElevator.Dir)
// 		timeUsed += Config.TRAVEL_TIME
//         count = count + 1
// 	}
// }
