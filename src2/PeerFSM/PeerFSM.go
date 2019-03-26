package PeerFSM

import{
  "fmt"
  "time"
  "../elevio"
}

FLOORS := 4
BUTTONS := 3

type state int
const (
	UNKNOWN  state = 1 + iota
	IDLE
	DOOR_OPEN
	MOVING
)

//type dir int
//const(
//  UP        dir = 1
//  DOWN
//  STOP
//)

type localState struct{
  State     state
  Floor     int
  Dir       elevio.MotorDirection
  Queue     [FLOORS][BUTTONS]bool
}

type buttonType int
const(
  HallUp    buttonType = 0
  HallDown
  Cab
)

type order struct{
  Floor         int
  Button        buttonType
}


type FSMChannels struct{
  CurrentFloor            chan int                  //from elevio
  LocalStateUpdate        chan localState           //to governor
  pingFromGov             chan order                //from governor
  //OrderFinished           chan order
}


func FSM(c FSMChannels) {
  go elevio.PollFloorSensor(c.CurrentFloor)
  floorInit := <-c.CurrentFloor

  // Initializing the local/internal state of the peer
  // MEN DETTE LIGGER EGT LOKALT I peer = GLOBALSTATE[minId].Elev
  peer := localState{
      State:  IDLE,
      Floor:  floorInit,
      Dir:    elevio.MD_Stop,
      Queue:  [FLOORS][BUTTONS]bool{},
  }
  // Setting timer
  doorTimerDone := time.NewTimer(0)

  // Sending first update (the initialisation)
  c.LocalStateUpdate<-peer

  for{
    select{
    case <-c.pingFromGov:
      // MÅ FYLLE DEN LOKALE KØEN MED DENNE NYE ORDREN!
      switch peer.state {
      case IDLE:
        peer.Dir = getNextDir(peer)
        elevio.SetMotorDirection(peer.Dir)
        if peer.Dir != elevio.MD_Stop {
          peer.state = MOVING
        } else {
          doorTimerDone= time.NewTimer(3 * time.Second)
          peer.State = DOOR_OPEN
        }
          //BØR VI HA EN ELSE HER? HVA SKJER NÅR getNextDir SIER DIR = STOP?

      case MOVING:
        // IKKE GJØRE NOE FORDI VI SKAL FULLFØRE PÅBEGYNT ORDRE

      case DOOR_OPEN:
        // Hva gjør vi om vi står med døren åpen og får ping fra
        // Gov om at det er oppdatering i kø, og denne oppdateringen er
        // at det er ny ordre i samme etasje som den står i med døren åpen?
        // For da risikerer vi å ikke stå der i 3 ekstra sek
        // Slik gjorde vi det før, men dette er ikke lenger mulig da order.Floor ikke finnes
        //if peer.Floor == order.Floor{
          //doorTimerDone = time.NewTimer(3 * time.Second)

        }

      }
      c.LocalStateUpdate <- peer


    case currentFloor:= <-atFloor:
      peer.Floor = currentFloor
      elevio.SetFloorIndicator(currentFloor)
      if shouldStop(peer){
        elevio.SetMotorDirection(elevio.MD_Stop)
        elevio.SetDoorOpenLamp(true)
        doorTimer = time.NewTimer(3 * time.Second)
        peer.State = DOOR_OPEN
        c.LocalStateUpdate <- peer
      }


      //else if floor == 0 {
        //elevio.SetMotorDirection(elevio.MD_Up)
      //} else if floor == 4 {
        //elevio.SetMotorDirection(elevio.MD_Down)
      //}

    case <-doorTimerDone.C:
      elevio.SetDoorOpenLamp(false)
      //SENDE CMD_FINISHED ?
      //c.OrderFinished <- order

      // Slette ordre fra queue
      peer.Queue[peer.Floor][HallUp] = false
      peer.Queue[peer.Floor][HallDown] = false
      peer.Queue[peer.Floor][Cab] = false

      //Oppdatere lys?

      peer.Dir = getNextDir(peer)
      if peer.Dir != elevio.MD_Stop {
        peer.State = MOVING
      } else {
        peer.State = IDLE
      }
      c.LocalStateUpdate <- peer
    }
  }
}


func getNextDir(peer localState) dir elevio.MotorDirection  {

  if peer.Dir == elevio.MD_Up {
      if orderAbove(peer.Floor){
          return elevio.MD_Up
      } else if orderBelow(peer.Floor + 1){
          return elevio.MD_Down
      } else {
          return elevio.MD_Stop
      }
  } else if (peer.Dir == MD_Down) {
      if orderBelow(peer.Floor){
          return MD_Down
      } else if orderAbove(peer.Floor - 1) {
          return MD_Up
      } else {
          return MD_Stop
      }
  } else {
      if orderAbove(peer.Floor){
          return MD_Up
      } else if orderBelow(peer.Floor){
          return MD_Down
      } else {
          return MD_Stop
      }
  }
}

func orderAbove(peer.Floor int) ans bool {
  for floor:= FLOORS; floor > peer.Floor floor--{
    for button:= 0; button < BUTTONS; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
    return false
}

func orderBelow(peer.Floor int) ans bool {
  for floor:= 0; floor < peer.Floor; floor++{
    for button: = 0; button < BUTTONS; button++{
        if peer.Queue[floor][button]{
            return true
        }
    }
}
return false
}

func shouldStop(peer) ans bool {
  if peer.Dir == MD_Up{
		if(peer.Queue[peer.Floor][HallUp] || peer.Queue[peer.Floor][Cab] || !orderAbove(peer.Floor)){
			return true
		}
	} else if peer.Dir == MD_Down {
		if(peer.Queue[peer.Floor][HallDown] || peer.Queue[peer.Floor][Cab] || !orderBelow(peer.Floor)){
			return true
		}
	} else if peer.Dir == MD_Stop  {
		if(peer.Queue[peer.Floor][HallUp] || peer.Queue[peer.Floor][HallDown|| peer.Queue[peer.Floor][Cab]){
			return true
		}
	}

//  if (peer.Dir == MD_Down && peer.Queue[peer.Floor][HallUp] && !queue_order_below(peer.Floor))
//	{
	//		return true
//	} else if (peer.Dir == MD_Up && peer.Queue[peer.Floor][HallDown] && !queue_order_above(peer.Floor)){
	//		return true
//	}
	return false
}
