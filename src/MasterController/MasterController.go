package MasterController

import (
	"../elevio"
	"fmt"
	"../CommTest"
)

type Dir int

const (
	UP              Dir = 0
	DOWN                = 1
	STOP                = 2
)

type optPassing struct{
	hallMat [4][2]int
	cabMat	[4][3]int
	num int
}

type elevator struct {
  Dir Dir
  Floor int
  OrderButton elevio.ButtonEvent
}

var elevators [3]elevator

func MasterInit(chComm CommTest.CommunicationChannels, chMaster CommTest.MasterControllerChannels){
	for id:=0; id<3; id++{
		elevators[id].Dir=STOP
	}
	go ManageCmd(chComm, chMaster)

	for {
    }
}


func ManageCmd(chComm CommTest.CommunicationChannels, chMaster CommTest.MasterControllerChannels){

	num:=1
	hallMat := [4][2]int {}
	cabMat := [4][3]int {}
	orderMatrices := optPassing{hallMat, cabMat, num}

	for {
		select{
		case orderMatrices.num = <- chComm.NumElev:

		case msg:= <- chComm.PosUpdateToMaster:
			elevators[msg.Id].Floor = msg.Floor

		case msg:= <- chComm.ButtonPushedToMaster:
			if msg.Button.Button == 2 {
				orderMatrices.cabMat[msg.Button.Floor][msg.Id] = 1
				//orderMatrices = intoptimize(msg.Id, chMaster, orderMatrices)
			} else {
				orderMatrices.hallMat[msg.Button.Floor][msg.Button.Button] = 1 //index[0][1] and [3][0] should not be accessed
			}
			//orderMatrices = extoptimize(chMaster, orderMatrices)
			fmt.Println("HallMat: ", orderMatrices.hallMat)
			fmt.Println("CabMat: ", orderMatrices.cabMat)
			//should write to logfile and send out log
		case msg:= <- chComm.CmdFinishedToMaster:
			orderMatrices.cabMat[msg.Floor][msg.Id] = 0
			if elevators[msg.Id].OrderButton.Button < 2{
			orderMatrices.hallMat[elevators[msg.Id].OrderButton.Floor][elevators[msg.Id].OrderButton.Button]=0
			}
			elevators[msg.Id].OrderButton.Floor = 0
			//orderMatrices = intoptimize(msg.Id, chMaster, orderMatrices)
			//orderMatrices = extoptimize(chMaster, orderMatrices)

		//should write to logfile and send out log
		}

	}
 }

func extoptimize(chMaster CommTest.MasterControllerChannels, orderMatrices optPassing)optPassing{

	command:= CommTest.Cmd{}
	hallMat:=orderMatrices.hallMat
	cabMat:=orderMatrices.cabMat
	num:= orderMatrices.num
    // check elevators that are free
    for id := 0; id < num; id++ {
		command.Id=id
        if elevators[id].OrderButton.Floor == 0 {
            for floor := 0; floor < 3; floor++ {

                if hallMat[floor][0] == 1 {
										command.Floor = floor
										elevators[id].OrderButton.Floor = floor
										elevators[id].OrderButton.Button = 	elevio.BT_HallUp
										hallMat[floor][0] = 2
										if elevators[id].Floor <= floor{
											elevators[id].Dir = UP
										} else {
											elevators[id].Dir = DOWN
										}
										chMaster.CmdElevToFloor <- command
										break
                }
                if hallMat[floor+1][1] == 1 {
										command.Floor = floor+1
										elevators[id].OrderButton.Floor = floor + 1
										elevators[id].OrderButton.Button= 	elevio.BT_HallDown
										hallMat[floor+1][1] = 2
										if elevators[id].Floor < floor+1{
											elevators[id].Dir = UP
										} else {
											elevators[id].Dir = DOWN
										}
										chMaster.CmdElevToFloor <- command
										break
                }
            }
        }
	    if elevators[id].Dir == UP{
	    	for floor := elevators[id].Floor; floor < elevators[id].OrderButton.Floor; floor++ {
				command.Floor = floor
	        	if hallMat[floor][0] == 1{
							hallMat[elevators[id].OrderButton.Floor][0] = 1
							elevators[id].OrderButton.Floor = floor
							hallMat[floor][0] = 2
							chMaster.CmdElevToFloor <- command
							break
	        	}
	    	}
		}
		if elevators[id].Dir == DOWN {
	        for floor := elevators[id].Floor; floor > elevators[id].OrderButton.Floor; floor-- {
							command.Floor = floor
	            if hallMat[floor][1] == 1{
									hallMat[elevators[id].OrderButton.Floor][1] = 1
									elevators[id].OrderButton.Floor = floor
									hallMat[floor][1] = 2
									chMaster.CmdElevToFloor <- command
									break
	          	}
		  	}
  		}
	}
	return optPassing{hallMat,cabMat,num}
}



func intoptimize(id int, chMaster CommTest.MasterControllerChannels,  orderMatrices optPassing) optPassing{
		command:= CommTest.Cmd{}
		command.Id = id
    currentFloor:= elevators[id].Floor
		hallMat:=orderMatrices.hallMat
		cabMat:=orderMatrices.cabMat
		num:= orderMatrices.num

    if elevators[id].Dir == UP || elevators[id].Dir == STOP{
    	for floor := currentFloor; floor < 4; floor++ {
        	if cabMat[floor][id] == 1{
						if elevators[id].OrderButton.Button < 2{
								if hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] == 2{
								hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] = 1
							}
						}
						command.Floor = floor
						elevators[id].Dir = UP
						elevators[id].OrderButton.Floor = floor
						elevators[id].OrderButton.Button = elevio.BT_Cab
						chMaster.CmdElevToFloor <- command
						break
					}

		}
		for floor := currentFloor; floor >= 0; floor-- {
				if cabMat[floor][id] == 1{
					if elevators[id].OrderButton.Button < 2{
							if hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] == 2{
							hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] = 1
						}
					}
					command.Floor = floor
					elevators[id].Dir = DOWN
					elevators[id].OrderButton.Floor = floor
					elevators[id].OrderButton.Button = elevio.BT_Cab
					chMaster.CmdElevToFloor <- command
					break
					}
				}
  		} else {
				for floor := currentFloor; floor >= 0; floor-- {
						if cabMat[floor][id] == 1{
							if elevators[id].OrderButton.Button < 2{
									if hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] == 2{
									hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] = 1
								}
							}
							command.Floor = floor
							elevators[id].OrderButton.Floor = floor
							elevators[id].OrderButton.Button = elevio.BT_Cab
							chMaster.CmdElevToFloor <- command
							break
						}
					}
			for floor := currentFloor; floor < 4; floor++ {
					if cabMat[floor][id] == 1{
						if elevators[id].OrderButton.Button < 2{
							if hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] == 2{
								hallMat[elevators[id].OrderButton.Floor][elevators[id].OrderButton.Button] = 1
							}
						}
						command.Floor = floor
						elevators[id].Dir = UP
						elevators[id].OrderButton.Floor = floor
						elevators[id].OrderButton.Button = elevio.BT_Cab
						chMaster.CmdElevToFloor <- command
						break
					}
			}
	}
return optPassing{hallMat,cabMat,num}
}


//
//
//   //In : hallmat, cabmat, elevators
//
// //Optimering må skje på buttonPushed og cmdFinished
//
// //Optimeringsprosessen er tosteget på cmdFinished og buttonpushed
// // - Prioriterer cab order og setter cab order som en temp
//   // - Sjekker først cab order i current DIR, så motsatt
// // - Prøver så å optimalisere veien til temp med hallMat
// // - Sender så temp som cmd
//
//   //Out: floor, id
//
//
//
//
//
// func detectSlaveError(){//own go routine
//   timerSlave0 = time.NewTimer(7 * time.Second)
//   timerSlave1 = time.NewTimer(7 * time.Second)
//   timerSlave2 = time.NewTimer(7 * time.Second)
//   for{
//     select{
//     case elevatorPos, elevatorId = <- posUpdateFromComm:
//       switch elevatorId{
//       case 0:
//         timerSlave0 = time.NewTimer(7 * time.Second)
//       case 1:
//         timerSlave1 = time.NewTimer(7 * time.Second)
//       case 2:
//         timerSlave2 = time.NewTimer(7 * time.Second)
//
//       }
//     case <- timerSlave0.C:
//       if numberOfSlaves > 0{
//           slaveError <- 0
//       }
//     case <- timerSlave1.C:
//       if numberOfSlaves > 1{
//           slaveError <- 1
//       }
//     case <- timerSlave2.C:
//       if numberOfSlaves > 2{
//           slaveError <- 2
//       }
//     }
//   }
// }
//
//
// numbFloors = 4
//
//

//
//
// func shortestPath(elevators, order) {
//   for elev in elevators {
//     if elev.outsideOrder.floor
//   }
// }
//
//
// func dist(fromFloor int, fromDir Dir, toFloor int, toDir Dir) {
//   toFloorIndex := 0
//   fromFloorIndex := 0
//   if(toDir == UP) {
//     toFloorIndex = toFloor + numbFloors
//   } else {
//     toFloorIndex = numbFloors - toFloor
//   }
//   return (toFloor- fromFloor)%(numbFloors-1)*2
// }
