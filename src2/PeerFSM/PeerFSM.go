package PeerFSM

import(
  "fmt"
  "time"
  "../elevio"
  "../Config"
)

type FSMChannels struct{
  CurrentFloor            chan int                  //from elevio
  LocalStateUpdate        chan localState           //to governor
  pingFromGov             chan int                //from governor
  //OrderFinished           chan order
}


func FSM(c FSMChannels, peer Config.Elev) {
  //go elevio.PollFloorSensor(c.CurrentFloor)
  //floorInit := <-c.CurrentFloor

  // Setting timer
  doorTimerDone := time.NewTimer(0)

  // Sending first update (the initialisation)
  c.LocalStateUpdate<-peer

  for{
    select{
    case <-c.pingFromGov:       // Ny oppdatering har skjedd
      switch peer.state {
      case Config.IDLE:
        peer.Dir = GetNextDir(peer)
        elevio.SetMotorDirection(peer.Dir)
        if peer.Dir != Config.DirStop {
          peer.state = Config.MOVING
        } else {
          if peer.Queue[peer.Floor][0] || peer.Queue[peer.Floor][1] || peer.Queue[peer.Floor][2]{
            doorTimerDone= time.NewTimer(3 * time.Second)
            // Slette ordre fra queue
            peer.Queue[peer.Floor][Config.BT_HallUp] = false
            peer.Queue[peer.Floor][Config.BT_HallDown] = false
            peer.Queue[peer.Floor][Config.BT_Cab] = false
            peer.State = Config.DOOR_OPEN
          } else {
            peer.State = Config.IDLE
            }
        }
      case Config.MOVING:
      case Config.DOOR_OPEN:
        if peer.Queue[peer.Floor][0] || peer.Queue[peer.Floor][1] || peer.Queue[peer.Floor][2]{
          doorTimerDone= time.NewTimer(3 * time.Second)
          // Slette ordre fra queue
          peer.Queue[peer.Floor][Config.BT_HallUp] = false
          peer.Queue[peer.Floor][Config.BT_HallDown] = false
          peer.Queue[peer.Floor][Config.BT_Cab] = false
          peer.State = Config.DOOR_OPEN
        }
      }
      c.LocalStateUpdate <- peer


    case currentFloor:= <-atFloor:
      peer.Floor = currentFloor
      elevio.SetFloorIndicator(currentFloor)
      if ShouldStop(peer){
        elevio.SetMotorDirection(Config.DirStop)
        elevio.SetDoorOpenLamp(true)
        doorTimerDone = time.NewTimer(3 * time.Second)
        // Slette ordre fra queue
        peer.Queue[peer.Floor][Config.BT_HallUp] = false
        peer.Queue[peer.Floor][Config.BT_HallDown] = false
        peer.Queue[peer.Floor][Config.BT_Cab] = false
        peer.State = Config.DOOR_OPEN
      }
      c.LocalStateUpdate <- peer

      //else if floor == 0 {
        //elevio.SetMotorDirection(elevio.MD_Up)
      //} else if floor == 4 {
        //elevio.SetMotorDirection(elevio.MD_Down)
      //}

    case <-doorTimerDone.C:
      elevio.SetDoorOpenLamp(false)

      //Oppdatere lys?

      peer.Dir = GetNextDir(peer)
      if peer.Dir != Config.DirStop {
        peer.State = Config.MOVING
      } else {
        peer.State = Config.IDLE
      }
      c.LocalStateUpdate <- peer
    }
  }
}

// FSM HELP FUNCTIONS

func GetNextDir(peer Config.Elev)Config.Direction  {

  if peer.Dir == Config.DirStop{
      if OrderAbove(peer){
          return Config.DirUp
      } else if OrderBelow(peer){
          return Config.DirDown
      } else {
          return Config.DirStop
      }
  } else if peer.Dir == Config.DirUp {
      if OrderAbove(peer){
          return Config.DirUp
      } else if OrderBelow(peer) {
          return Config.DirDown
      } else {
          return Config.DirStop
      }
  } else {
      if OrderBelow(peer){
          return Config.DirDown
      } else if OrderAbove(peer){
          return Config.DirUp
      } else {
          return Config.DirStop
      }
  }
  return Config.DirStop
}

func OrderAbove(peer Config.Elev)bool {
  for floor:= peer.Floor + 1; floor < Config.NumFloors; floor++{
    for button:= 0; button < Config.NumButtons-1; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func OrderBelow(peer Config.Elev)bool {
  for floor:= 0; floor < peer.Floor; floor++{
    for button:= 0; button < Config.NumButtons-1; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func ShouldStop(peer Config.Elev)bool {
    fmt.Println("peer.Dir: ", peer.Dir)
    fmt.Println("peer.Floor: ", peer.Floor)
    switch peer.Dir {
    case Config.DirUp:
        //fmt.Println("dirUp: ", peer.Queue[peer.Floor][Config.BT_HallUp], peer.Queue[peer.Floor][Config.BT_Cab], !OrderAbove(peer))
        return(peer.Queue[peer.Floor][Config.BT_HallUp] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderAbove(peer))
    case Config.DirDown:
        //fmt.Println("dirDown: ", peer.Queue[peer.Floor][Config.BT_HallDown],  peer.Queue[peer.Floor][Config.BT_Cab], !OrderBelow(peer))
        return (peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderBelow(peer))
    case Config.DirStop:
        // fmt.Println("dirstop: ", peer.Queue[peer.Floor][Config.BT_HallUp ], peer.Queue[peer.Floor][Config.BT_HallDown], peer.Queue[peer.Floor][Config.BT_Cab])
        // if(peer.Queue[peer.Floor][Config.BT_HallUp ] || peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab]){
        //     return true
        // }
    default:
    }
    return false
}
