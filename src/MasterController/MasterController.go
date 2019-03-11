package MasterController

import "../elevio"

// TODO:
// - Initialisere master
// - Få log, Fylle log, sende log
// - Kalkulere optimal path/finne kommandoer til slavene
// - Sende kommandoer til slavene


func manageCmd(){
  //Inn: heisenes pos(posUpdate), buttonPushed
  //Ut: kommando til heis

  //Declaring variables
  position:= int
  id:= int
  hallMat := [4][2]int {}
  cabMat := [4][3]int {}
  buttonEvent := elevio.ButtonEvent {}

  case position, id = <- posUpdateFromComm:
      elevators[id].floor = position
  case buttonEvent, id = <- buttonPushedFromComm
      if buttonEvent.Button == 2 {
        cabMat[buttonEvent.Floor][id] = 1 //Syntax should be checked
      } else {
        hallMat[buttonEvent.Floor][buttonEvent.Button] = 1 //index[0][1] and [3][0] should not be accessed
      }
      //optimize order init

      //should write to logfile and send out log
  case finishedCmd, id = <- cmdFinishedFromComm 
      cabMat[finishedCmd][id] = 0
      hallMat[finishedCmd][] 
      //Not sure about this one. Should matrices be updated when orders are delegated or finished?

}

func optimize(){
  //In : hallmat, cabmat, elevators

//Optimering må skje på buttonPushed og cmdFinished

//Optimeringsprosessen er tosteget på cmdFinished og buttonpushed
// - Prioriterer cab order og setter cab order som en temp
  // - Sjekker først cab order i current DIR, så motsatt
// - Prøver så å optimalisere veien til temp med hallMat 
// - Sender så temp som cmd

  //Out: floor, id

}

func masterInit(){

}

func detectSlaveError(){

}

numbFloors = 4

type Dir int
struct elevator {
  dir Dir = 0
  floor = 0
  dist_order outsideOrder
}
struct outsideOrder {
  floor int = 0
  dir Dir = 0
}

var elevators [3]elevator

func shortestPath(elevators, order) {
  for elev in elevators {
    if elev.outsideOrder.floor
  }
}


func dist(fromFloor int, fromDir Dir, toFloor int, toDir Dir) {
  toFloorIndex := 0
  fromFloorINdex := 0
  if(toDir == UP) {
    toFloorIndex = toFloor + numbFloors
  } else {
    toFloorIndex = numbFloors - toFloor
  }
  return (toFloor- fromFloor)%(numbFloors-1)*2
}
