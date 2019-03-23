package MasterController

import (//"../elevio"
	"fmt"
	"../CommTest"
)

type Dir int

const (
	UP              Dir = 0
	DOWN                = 1
	STOP                = 2
)

type outsideOrder struct{
  Floor int
  Dir Dir
}

type elevator struct {
  Dir Dir
  Floor int
  OutOrder outsideOrder
}

var elevators [3]elevator

func MasterInit(chComm CommTest.CommunicationChannels, chMaster CommTest.MasterControllerChannels){
	go ManageCmd(chComm, chMaster)

	for {
        // select {
		// case msg:= <- chComm.PosUpdateToMaster:
		// 	fmt.Println("I am Master, PosUpdate: ", msg)
		// case msg:= <- chComm.CmdFinishedToMaster:
		// 	fmt.Println("I am Master, cmdFinished: ", msg)
		// case msg:= <- chComm.ButtonPushedToMaster:
		// 	fmt.Println("I am Master, buttonpushed: ", msg)
		// 	command:= CommTest.Cmd{}
		// 	command.Floor = msg.Button.Floor
		// 	command.Id = msg.Id
		// 	chMaster.CmdElevToFloor <- command
		// }

    }
}
func ManageCmd(chComm CommTest.CommunicationChannels, chMaster CommTest.MasterControllerChannels){
	//Declaring variables
	hallMat := [4][2]int {}
	cabMat := [4][3]int {}
	var num int
	for {
		select{
		case num = <- chComm.NumElev:


		case msg:= <- chComm.PosUpdateToMaster:

			elevators[msg.Id].Floor = msg.Floor

		case msg:= <- chComm.ButtonPushedToMaster:
			if msg.Button.Button == 2 {
				cabMat[msg.Button.Floor][msg.Id] = 1
				//intoptimize()
			} else {
				hallMat[msg.Button.Floor][msg.Button.Button] = 1 //index[0][1] and [3][0] should not be accessed
				go extoptimize(num, chComm, chMaster, hallMat, cabMat) //should this be a go routine?
			}
			//should write to logfile and send out log
		case msg:= <- chComm.CmdFinishedToMaster:
			fmt.Println("msg from cmd finished:", msg)
			cabMat[msg.Floor][msg.Id] = 0
			if elevators[msg.Id].Dir < 2{
			hallMat[elevators[msg.Id].OutOrder.Floor][elevators[msg.Id].OutOrder.Dir]=0
			}
			// else if elevators[msg.Id].Dir==2{
			// 	hallMat[elevators[msg.Id].OutOrder.Floor][elevators[msg.Id].OutOrder.Dir] = 0
			// }
			elevators[msg.Id].OutOrder.Floor = 0
			//intoptimize()
			go extoptimize(num, chComm, chMaster, hallMat, cabMat)
		//should write to logfile and send out log
		}
	}
 }


func extoptimize(num int, chComm CommTest.CommunicationChannels, chMaster CommTest.MasterControllerChannels,
	hallMat [4][2]int, cabMat[4][3]int){

	command:= CommTest.Cmd{}
	//chMaster.CmdElevToFloor <- command
	fmt.Println("Hallmat: ", hallMat)
    // check elevators that are free
    for id := 0; id < num; id++ {
		command.Id=id
        if elevators[id].OutOrder.Floor == 0 {
            for floor := 0; floor < 3; floor++ {

                if hallMat[floor][0] == 1 {
					command.Floor = floor
					elevators[id].OutOrder.Floor = floor
					elevators[id].OutOrder.Dir = UP
					if elevators[id].Floor < floor{
						elevators[id].Dir = UP
					} else if elevators[id].Floor < floor{
						elevators[id].Dir = STOP
					} else {
						elevators[id].Dir = DOWN
					}
					fmt.Println("Command.Floor",command.Floor)
					chMaster.CmdElevToFloor <- command
					break
                }
                if hallMat[floor+1][1] == 1 {
					command.Floor = floor+1
					elevators[id].OutOrder.Floor = floor + 1
					elevators[id].OutOrder.Dir = DOWN
					if elevators[id].Floor < floor+1{
						elevators[id].Dir = UP
					} else if elevators[id].Floor < floor{
						elevators[id].Dir = STOP
					}  else {
						elevators[id].Dir = DOWN
					}
					fmt.Println("Command.Floor",command.Floor)
					chMaster.CmdElevToFloor <- command
					break
                }
            }
        }


	    if elevators[id].Dir == UP{
	    	for floor := elevators[id].Floor; floor < elevators[id].OutOrder.Floor; floor++ {
				command.Floor = floor
	        	if hallMat[floor][0] == 1{
					elevators[id].OutOrder.Floor = floor
					chMaster.CmdElevToFloor <- command
					fmt.Println("Extra opt")
					break
	        	}
	    	}
		}
		if elevators[id].Dir == DOWN {
	        for floor := elevators[id].Floor; floor > elevators[id].OutOrder.Floor; floor-- {
				command.Floor = floor
	            if hallMat[floor][1] == 1{
					elevators[id].OutOrder.Floor = floor
					chMaster.CmdElevToFloor <- command
					fmt.Println("Extra opt")
					break
	          	}
		  	}
  		}
	}
	fmt.Println("Opt finished")
}



// func intoptimize(int id){
//     temp:= elevators[id].floor
//
//     if elevators[id].dir == UP {
//     	for count := 0; count < 4; count++ {
//         	temp = (temp + count)%4
//         	if cabMat[temp][id] == 1{
//           	//send order and quit for loop
//           	//update dir
// 			//update OutOrder
//
//     		} else {
// 	    		for count := 4; count > 0; count-- {
// 	        		temp = (temp + count)%4
// 	        		if cabMat[temp][id] == 1{
//           		//send order and quit for loop
//           		//update dir
// 				//update OutOrder
//
//       				}
//     			}
// 			}
// 		}
// 	}
// }
//

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
