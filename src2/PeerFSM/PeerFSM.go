package PeerFSM

import(
  "fmt"
  "time"
  "../elevio"
  "../Config"
)

func Lights(GState Config.GlobalState){
    for floor:= 0;  floor < Config.NumFloors; floor++ {
        if GState.Map["minId"].Queue[floor][2] {
            elevio.SetButtonLamp(Config.BT_Cab, floor, true)
        } else {
            elevio.SetButtonLamp(Config.BT_Cab, floor, false)
        }
        if GState.HallRequests[floor][0]{
            elevio.SetButtonLamp(Config.BT_HallUp, floor, true)
        } else {
            elevio.SetButtonLamp(Config.BT_HallUp, floor, false)
        }
        if GState.HallRequests[floor][1]{
            elevio.SetButtonLamp(Config.BT_HallDown, floor, true)
        } else {
            elevio.SetButtonLamp(Config.BT_HallDown, floor, false)
        }
    }

}
func FSM(c Config.FSMChannels, peer Config.Elev) {
  // Setting timer
  doorTimerDone := time.NewTimer(0)
  doorTimerDone.Stop()
  elevio.SetMotorDirection(Config.MD_Up)
  fmt.Println(peer.Queue)
  for{
  fmt.Println("STATE: ", peer.State, "DIR: ", peer.Dir, "FLOOR: ", peer.Floor)
    select{
    case <-c.PingFromGov:
        fmt.Println("CASE A")
      switch peer.State {
      case Config.IDLE:
        peer.Dir = GetNextDir(peer)
          fmt.Println(peer.Queue)
        fmt.Println("Next dir: ", peer.Dir)
        elevio.SetMotorDirection(peer.Dir)
        if peer.Dir != Config.MD_Stop {
          peer.State = Config.MOVING
                c.LocalStateUpdate <- peer
        } else {
          if peer.Queue[peer.Floor][0] || peer.Queue[peer.Floor][1] || peer.Queue[peer.Floor][2]{
            doorTimerDone= time.NewTimer(3 * time.Second)
            // Slette ordre fra queue
            peer.Queue[peer.Floor][Config.BT_HallUp] = false
            peer.Queue[peer.Floor][Config.BT_HallDown] = false
            peer.Queue[peer.Floor][Config.BT_Cab] = false
            peer.State = Config.DOOR_OPEN
                  c.LocalStateUpdate <- peer
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



  case currentFloor:= <-c.CurrentFloor:
      fmt.Println("CASE B")
      peer.Floor = currentFloor
      elevio.SetFloorIndicator(currentFloor)
      switch peer.State {
      case Config.UNKNOWN:
              elevio.SetMotorDirection(elevio.MD_Stop)
              peer.State = Config.IDLE
      case Config.MOVING:
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
    }


      //else if floor == 0 {
        //elevio.SetMotorDirection(elevio.MD_Up)
      //} else if floor == 4 {
        //elevio.SetMotorDirection(elevio.MD_Down)
      //}

    case <-doorTimerDone.C:
    fmt.Println("CASE C")
      elevio.SetDoorOpenLamp(false)

      //Oppdatere lys?

      peer.Dir = GetNextDir(peer)
      elevio.SetMotorDirection(peer.Dir)
      fmt.Println("Doortimer done found next dir: ", peer.Dir)
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
    for button:= 0; button < Config.NumButtons; button++{
        if peer.Queue[floor][button]{
            fmt.Println("order above")
            return true
        }
    }
}
return false
}

func OrderBelow(peer Config.Elev)bool {
  for floor:= 0; floor < peer.Floor; floor++{
    for button:= 0; button < Config.NumButtons; button++{
        if peer.Queue[floor][button]{
                        fmt.Println("order below")
            return true
        }
    }
}
return false
}

func ShouldStop(peer Config.Elev)bool {
    //fmt.Println("peer.Dir: ", peer.Dir)
    //fmt.Println("peer.Floor: ", peer.Floor)
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
