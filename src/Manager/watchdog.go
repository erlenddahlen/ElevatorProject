package Manager
import (
	"../DataStructures"
	 //"fmt"
	 "time"
)

//Watchdog should be own go routine
//In SpamGlobalState, send latestState to Watchdog on own channel
//Watchdog should be a case in for/select in UpdateGlobalState
//When activated, put all Hallreq in own queue

func Watchdog(gchan DataStructures.ManagerChannels, GState DataStructures.GlobalState){
	watchdogTimer := time.NewTimer(6 * time.Second)
	var changedState bool
	var hallreqExist bool
	var prevState DataStructures.GlobalState
	prevState = GState


	for{
		select{
		case newState := <- gchan.UpdatefromSpam:
			changedState = false
			hallreqExist = false

			// Iterate through Hallreqs, if exists orders hallreq = true
			for floor:= 0; floor < DataStructures.NumFloors; floor++ {
				for button := 0; button < DataStructures.NumButtons-1; button++ {
					if newState.HallRequests[floor][button] {
						hallreqExist = true
						break
					}
				}
			}

			// Iterate throug and check if change in elev floor and state between latestState
			// and newState, if exists change = true

			for key, _:= range newState.Map{
				if newState.Map[key].State != prevState.Map[key].State ||  newState.Map[key].Floor != prevState.Map[key].Floor{
					changedState = true
					break
				}
			}

			if changedState || !(hallreqExist){
					watchdogTimer = time.NewTimer(6 * time.Second)
			}
			prevState = newState

		case <- watchdogTimer.C:
			gchan.Watchdog <- 1
		}
	}
}
