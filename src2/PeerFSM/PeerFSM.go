package PeerFSM

import(
  "fmt"
  "time"
  "../elevio"
  "../Config"
)


func FSM(c Config.FSMChannels, peer Config.Elev) {
  // Setting timer
  doorTimerDone := time.NewTimer(0)

  for{
    select{
    case <-c.PingFromGov:
      switch peer.State {
      case Config.IDLE:
        peer.Dir = GetNextDir(peer)
        elevio.SetMotorDirection(peer.Dir)
        if peer.Dir != Config.MD_Stop {
          peer.State = Config.MOVING
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


  case currentFloor:= <-c.CurrentFloor:
      peer.Floor = currentFloor
      elevio.SetFloorIndicator(currentFloor)
      if ShouldStop(peer){
        elevio.SetMotorDirection(Config.MD_Stop)
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
      if peer.Dir != Config.MD_Stop {
        peer.State = Config.MOVING
      } else {
        peer.State = Config.IDLE
      }
      c.LocalStateUpdate <- peer
    }
  }
}

// FSM HELP FUNCTIONS

func GetNextDir(peer Config.Elev)Config.MotorDirection  {

  if peer.Dir == Config.MD_Stop{
      if OrderAbove(peer){
          return Config.MD_Up
      } else if OrderBelow(peer){
          return Config.MD_Down
      } else {
          return Config.MD_Stop
      }
  } else if peer.Dir == Config.MD_Up {
      if OrderAbove(peer){
          return Config.MD_Up
      } else if OrderBelow(peer) {
          return Config.MD_Down
      } else {
          return Config.MD_Stop
      }
  } else {
      if OrderBelow(peer){
          return Config.MD_Down
      } else if OrderAbove(peer){
          return Config.MD_Up
      } else {
          return Config.MD_Stop
      }
  }
  return Config.MD_Stop
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
    case Config.MD_Up:
        //fmt.Println("MD_Up: ", peer.Queue[peer.Floor][Config.BT_HallUp], peer.Queue[peer.Floor][Config.BT_Cab], !OrderAbove(peer))
        return(peer.Queue[peer.Floor][Config.BT_HallUp] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderAbove(peer))
    case Config.MD_Down:
        //fmt.Println("MD_Down: ", peer.Queue[peer.Floor][Config.BT_HallDown],  peer.Queue[peer.Floor][Config.BT_Cab], !OrderBelow(peer))
        return (peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab] || !OrderBelow(peer))
    case Config.MD_Stop:
        // fmt.Println("MD_Stop: ", peer.Queue[peer.Floor][Config.BT_HallUp ], peer.Queue[peer.Floor][Config.BT_HallDown], peer.Queue[peer.Floor][Config.BT_Cab])
        // if(peer.Queue[peer.Floor][Config.BT_HallUp ] || peer.Queue[peer.Floor][Config.BT_HallDown] || peer.Queue[peer.Floor][Config.BT_Cab]){
        //     return true
        // }
    default:
    }
    return false
}
