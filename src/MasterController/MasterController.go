package MasterController

import "../elevio"

// TODO:
// - Initialisere master
// - Få log, Fylle log, sende log
// - Kalkulere optimal path/finne kommandoer til slavene
// - Sende kommandoer til slavene
type Dir int

const (
	UP              Dir = 0
	DOWN                = 1
	STOP                = 2
)

struct elevator {
  dir Dir = 2
  floor = 0
  outOrder = 0
}

var elevators [3]elevator


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
        intoptimize()
      } else {
        hallMat[buttonEvent.Floor][buttonEvent.Button] = 1 //index[0][1] and [3][0] should not be accessed
        extoptimize()
      }
      //should write to logfile and send out log
  case finishedCmd, id = <- cmdFinishedFromComm
      cabMat[finishedCmd][id] = 0
      elevators[id].outOrder = 0
      if elevators[id].dir < 2 {
        hallMat[finishedCmd][elevators[id].dir]
      }
      intoptimize()
      //should write to logfile and send out log
 }

func intoptimize(int id){
    temp:= elevators[id].floor
    runOptimize:= true
    if elevators[id].dir == UP {
    for count := 0; count < 4; count++ {
        temp = (temp + count)%4
        if cabMat[temp][id] == 1
          //send order and quit for loop
          //update dir
          //rO = false
    }else {
    for count := 4; count > 0; count-- {
        temp = (temp + count)%4
        if cabMat[temp][id] == 1
          //send order and quit for loop
          //update dir
          //rO = false
      }
    }
    if bool == true {
        extoptimize()
    }
    }

func extoptimize(){
    temp:= int
    // check elevators that are free
    for id := 0; id < 3; id++ {
        if elevators[id].outOrder == 0 {
            for i := 0; i < 3; i++ {
                if hallMat[i][0] == 1{
                    //quit for loop, do not send order
                    //update dir, outOrder
                }
                if hallMat[i+1][1] == 1 {
                    //quit for loop, do not send order
                    //update dir, outOrder
                }
            }
        }
        temp = elevators[id].floor
        if elevators[id].dir == UP
        for count := 0; count < 4; count++ {
            temp = (temp + count)%4
            if hallMat[temp][0] == 1{
              //send order and quit for loop
            }
        }else if elevators[id].dir == DOWN {
            for count := 4; count > 0; count-- {
                temp = (temp + count)%4
                if hallMat[temp][1] == 1
                  //send order and quit for loop
              }
        }
    }

    }


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
    for {
        switch {}


    }
}

func detectSlaveError(){
    //declare/initialized

    //for -- select -- switch
    for {
        select{
            switch{

            }
        }
    }
}

numbFloors = 4


struct outsideOrder {
  floor int = 0
  dir Dir = 0
}


func shortestPath(elevators, order) {
  for elev in elevators {
    if elev.outsideOrder.floor
  }
}


func dist(fromFloor int, fromDir Dir, toFloor int, toDir Dir) {
  toFloorIndex := 0
  fromFloorIndex := 0
  if(toDir == UP) {
    toFloorIndex = toFloor + numbFloors
  } else {
    toFloorIndex = numbFloors - toFloor
  }
  return (toFloor- fromFloor)%(numbFloors-1)*2
}
